package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vote/app/model"
	"vote/app/tools"
)

func Index(context *gin.Context) {
	ret := model.GetVotes()
	context.HTML(http.StatusOK, "index.html", gin.H{"vote": ret})
}

func GetVotes(context *gin.Context) {
	ret := model.GetVotes()
	context.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
}

func GetVoteInfo(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetVote(id)
	context.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
}

func DoVote(context *gin.Context) {
	userIDstr, _ := context.Cookie("Id")
	voteIdstr, _ := context.GetPostForm("vote_id")
	optstr, _ := context.GetPostFormArray("opt[]")

	UserID, _ := strconv.ParseInt(userIDstr, 10, 64)
	voteID, _ := strconv.ParseInt(voteIdstr, 10, 64)

	//	前置查询 悲观锁，乐观锁，分布式锁，消息队列
	old := model.GetVoteHistory(UserID, voteID)
	if len(old) >= 1 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "您已经投过票",
		})
	}

	opt := make([]int64, 0)
	for _, v := range optstr {
		opyId, _ := strconv.ParseInt(v, 10, 64)
		opt = append(opt, opyId)
	}

	model.DoVote(UserID, voteID, opt)
	context.JSON(http.StatusOK, tools.ECode{
		Message: "投票完成",
	})
}
