package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"vote/app/model"
	"vote/app/tools"
)

// AddVote godoc
// @Summary 添加投票
// @Description 添加新的投票主题和选项
// @Tags vote
// @Accept json
// @Produce json
// @Param title body string true "投票标题"
// @Param optStr body []string true "投票选项名称"
// @Success 200 {object} tools.ECode
// @Router /addVote [post]
type VoteAdd struct {
	Title  string   `json:"title"`
	OptStr []string `json:"optStr"`
}

func AddVote(context *gin.Context) {
	var voteadd VoteAdd
	//idStr := voteadd.title
	//var idStr = context.Param("title") //optStr := voteadd.optStr
	//optStr:= context.Param("optStr")
	if err := context.ShouldBind(&voteadd); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}
	fmt.Printf("有没有：%s", voteadd.Title)
	//构建结构体
	vote := model.Vote{
		Title:       voteadd.Title,
		Type:        0,
		Status:      0,
		CreatedTime: time.Now(),
	}
	if vote.Title == "" {
		context.JSON(http.StatusBadRequest, tools.ParamErr)
		return
	}
	oldVote := model.GetVoteByName(voteadd.Title)
	if oldVote.Id > 0 {
		context.JSON(http.StatusOK, tools.VoteAlready)
		return
	}
	opt := make([]model.VoteOpt, 0)
	for _, v := range voteadd.OptStr {
		opt = append(opt, model.VoteOpt{
			Name:        v,
			CreatedTime: time.Now(),
		})
	}
	if err := model.Addvote(vote, opt); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, tools.OK)
	return
}

// UpdateVote godoc
// @Summary 更新投票
// @Description 更新已有投票主题和选项
// @Tags vote
// @Accept json
// @Produce json
// @Param title body string true "投票标题"
// @Param optStr body []string true "投票选项名称"
// @Success 200 {object} tools.ECode
// @Router /updateVote [post]
func UpdateVote(context *gin.Context) {
	var voteadd VoteAdd
	if err := context.ShouldBind(&voteadd); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}
	token := context.GetHeader("Authorization")
	user, err := model.CheckJwt(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, tools.NotLogin)
		context.Abort()
		return
	}
	userid := user.Id
	//userid := model.UserId(context)

	//构建结构体
	vote := model.Vote{
		Title:       voteadd.Title,
		Type:        0,
		Status:      0,
		CreatedTime: time.Now(),
		UserId:      userid,
	}
	if vote.Title == "" {
		context.JSON(http.StatusBadRequest, tools.ParamErr)
		return
	}
	oldVote := model.GetVoteByName(voteadd.Title)
	fmt.Printf("%d\n", oldVote.Id)
	if oldVote.Id < 1 {
		context.JSON(http.StatusOK, tools.VoteNOAlready)
		return
	}
	opt := make([]model.VoteOpt, 0)
	for _, v := range voteadd.OptStr {
		opt = append(opt, model.VoteOpt{
			Name:        v,
			CreatedTime: time.Now(),
		})
	}
	if err := model.Updatevote(vote, opt); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, tools.OK)
	return
}

// DelVote godoc
// @Summary 删除投票
// @Description 根据投票ID删除投票主题和选项
// @Tags vote
// @Produce json
// @Param id query string true "投票ID"
// @Success 200 {object} tools.ECode
// @Failure 404 {string} string "未找到对应投票"
// @Router /delVote [delete]
func DelVote(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	//fmt.Printf("id接收到了没有%s\n", idStr)
	id, _ = strconv.ParseInt(idStr, 10, 64)
	vote := model.GetVoteCache(context, id)
	if vote.Vote.Id <= 0 {
		context.JSON(http.StatusOK, tools.OK)
		return
	}
	//fmt.Printf("这个有没有%d\n", vote.Vote.Id)
	if err := model.DelVote(id); err != true {
		context.JSON(http.StatusOK, tools.DelFail)
		return
	}
	context.JSON(http.StatusOK, tools.OK)
}

type ResultData struct {
	Title string
	Count int64
	Opt   []*ResultVoteOpt
}

type ResultVoteOpt struct {
	Name  string
	Count int64
}

// ResultVote godoc
// @Summary 获取投票结果
// @Description 根据投票ID获取投票结果数据
// @Tags vote
// @Produce json
// @Param id query int true "投票ID"
// @Success 200 {object} tools.ECode
// @Router /resultVote [get]
func ResultVote(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetVoteCache(context, id)
	data := ResultData{
		Title: ret.Vote.Title,
	}
	for _, v := range ret.Opt {
		data.Count = data.Count + v.Count
		tmp := ResultVoteOpt{
			Name:  v.Name,
			Count: v.Count,
		}
		data.Opt = append(data.Opt, &tmp)
	}
	context.JSON(http.StatusOK, tools.ECode{
		Data: data,
	})
}

// ResultInfo godoc
// @Summary 获取投票结果页面
// @Description 获取展示投票结果的页面
// @Tags vote
// @Produce html
// @Success 200
// @Router /resultInfo [get]
func ResultInfo(context *gin.Context) {
	context.HTML(http.StatusOK, "result.html", nil)
}
