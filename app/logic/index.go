package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"vote/app/model"
	"vote/app/tools"
)

// Index godoc
// @Summary 显示投票页面
// @Description 显示投票页面
// @Tags vote
// @Produce json
// @Success 200 {object} tools.ECode
// @Router / [get]
func Index(context *gin.Context) {
	ret := model.GetVotes()
	context.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
}

// GetVotes godoc
// @Summary 获取所有投票信息
// @Description 获取所有投票信息
// @Tags vote
// @Produce json
// @Success 200 {object} tools.ECode
// @Router /votes [get]
func GetVotes(context *gin.Context) {
	ret := model.GetVotes()
	context.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
}

// GetVoteInfo godoc
// @Summary 获取特定投票信息
// @Description 获取特定投票信息
// @Tags vote
// @Produce json
// @param   id query string true "id"
// @Success 200 {object} tools.ECode
// @Router /vote [get]
func GetVoteInfo(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret, _ := model.GetVoteV5(id)
	if ret.Vote.Id <= 0 {
		context.JSON(http.StatusNotFound, tools.ECode{})
		return
	}
	context.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
}

// DoVote godoc
// @Summary 进行投票
// @Description 进行投票
// @Tags vote
// @Accept json
// @Produce json
// @Param user_id cookie string true "用户ID"
// @Param vote_id body int64 true "投票ID"
// @Param opt[] body []int64 true "投票选项ID列表"
// @Success 200 {object} tools.ECode
// @Router /vote/do [post]
type VoteData struct {
	VoteID int64   `json:"vote_id"`
	Opt    []int64 `json:"opt"`
}

func DoVote(context *gin.Context) {
	var Votedata VoteData
	if err := context.ShouldBind(&Votedata); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}
	//tokenstr, _ := context.Get("token")
	//tokeStr, _ := tokenstr.(string)
	//user, _ := model.CheckJwt(tokeStr)
	//UserID := user.Id
	token := context.GetHeader("Authorization")
	user, err := model.CheckJwt(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, tools.NotLogin)
		context.Abort()
		return
	}
	UserID := user.Id

	//UserID := model.UserId(context)
	voteID := Votedata.VoteID
	opt := Votedata.Opt
	//fmt.Printf("能打印出什么东西%s,%d,%d", userIDstr, UserID, voteID, opt)
	//voteIdstr, _ := context.GetPostForm("vote_id")
	//optstr, _ := context.GetPostFormArray("opt[]")

	//	前置查询 悲观锁，乐观锁，分布式锁，消息队列
	//fmt.Printf("%d\n,%d\n,%d\n", UserID, voteID, opt)
	old := model.GetVoteHistoryV1(context, UserID, voteID)
	fmt.Println("检查是否投过票失效了??old %", old)
	if len(old) >= 1 {
		context.JSON(http.StatusOK, tools.VoteExits)
		return
	}

	isSuccess := model.DoVoteV2(UserID, voteID, opt, context)
	if isSuccess {
		context.JSON(http.StatusOK, tools.Success)
	}
}

func CheckXYZ(context *gin.Context) bool {
	ip := context.ClientIP()
	ua := context.GetHeader("user-agent")
	fmt.Printf("ip:%s\nua:%s\n", ip, ua)

	//转为MD5
	hash := md5.New()
	hash.Write([]byte(ip + ua))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	flag, _ := model.Rdb.Get(context, "ban-"+hashString).Bool()
	if flag {
		return false
	}

	i, _ := model.Rdb.Get(context, "xyz-"+hashString).Int()
	if i > 5 {
		model.Rdb.SetEx(context, "ban-"+hashString, true, 30*time.Second)
		return false
	}
	model.Rdb.Incr(context, "xyz-"+hashString)
	model.Rdb.Expire(context, "xyz-"+hashString, 50*time.Second)
	return true

}

// GetCaptcha godoc
// @Summary 获取验证码
// @Description 获取验证码
// @Tags login
// @Produce json
// @Success 200 {object} tools.ECode
// @Router /captcha [get]
func GetCaptcha(context *gin.Context) {
	if !CheckXYZ(context) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "您的手速真是太快了！",
		})
	}
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, tools.ECode{
		Data: captcha,
	})
}
