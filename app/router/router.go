package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"vote/app/logic"
	"vote/app/model"
	"vote/app/tools"
	_ "vote/docs"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	//相关的路径，放在这里
	//r.Use(tools.Cors())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	index := r.Group("")
	index.Use(checkUser)
	//{
	//	//vote相关
	//	index.GET("/index", logic.Index)
	//
	//	index.POST("/vote", logic.DoVote)
	//
	//	index.POST("/vote/add", logic.AddVote)
	//	index.POST("/vote/update", logic.UpdateVote)
	//	index.POST("/vote/del", logic.DelVote)
	//
	//	index.GET("/result", logic.ResultInfo)
	//	index.GET("/result/info", logic.ResultVote)
	//}

	// 改造为resultful接口
	{ //读
		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)

		index.POST("/vote", logic.AddVote)
		index.PUT("/vote", logic.UpdateVote)
		index.DELETE("/vote", logic.DelVote)

		index.GET("/vote/result", logic.ResultVote)

		index.POST("/do_vote", logic.DoVote)

	}
	r.GET("/", logic.Index)

	{
		r.GET("/login", logic.GetLogin)
		//登录
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)

		//注册
		r.POST("/user/create", logic.CreateUser)

	}
	//验证码
	{
		r.GET("/captcha", logic.GetCaptcha)

		r.POST("/captcha/verify", func(context *gin.Context) {
			var param tools.CaptchaData
			if err := context.ShouldBind(&param); err != nil {
				context.JSON(http.StatusOK, tools.ParamErr)
				return
			}

			fmt.Printf("参数为：%+v", param)
			if !tools.CaptchaVerify(param) {
				context.JSON(http.StatusOK, tools.ECode{
					Code:    10008,
					Message: "验证失败",
				})
				return
			}
			context.JSON(http.StatusOK, tools.OK)
		})
	}

	if err := r.Run(":8080"); err != nil {
		panic("gin启动失败")
	}
}

func checkUser(context *gin.Context) {
	var name string
	var id int64
	//values := model.GetSession(context)
	token := context.GetHeader("Authorization")
	fmt.Printf("token是什么?? 有没有:%s\n", token)
	user, err := model.CheckJwt(token)
	fmt.Printf("user看看打印出来了没有:%s\n", user)
	if err != nil {
		context.JSON(http.StatusUnauthorized, tools.NotLogin)
		context.Abort()
		return
	}

	name = user.Name
	id = user.Id
	//if v, ok := values["name"]; ok {
	//	name = v.(string)
	//}
	//if v, ok := values["id"]; ok {
	//	id = v.(int64)
	//}

	if name == "" || id < 0 {
		context.JSON(http.StatusUnauthorized, tools.NotLogin)
		context.Abort()
	}
	context.Next()
}
