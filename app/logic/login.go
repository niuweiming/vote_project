package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"vote/app/model"
	"vote/app/tools"
)

// 逻辑层
type User struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

func PostLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: err.Error(), //这里有风险
		})
	}
	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "验证码校验失败！", //这里有风险
		})
		return
	}

	ret := model.GetUser(user.Name)

	if ret.Id < 1 || ret.Password != encrypt(user.Password) {
		context.JSON(http.StatusOK, tools.UserErr)
		return
	}
	context.SetCookie("name", user.Name, 3600, "/", "", true, false)
	context.SetCookie("Id", fmt.Sprintf(strconv.FormatInt(ret.Id, 10)), 3600, "/", "", true, false)
	//

	_ = model.SetSession(context, user.Name, ret.Id)

	context.JSON(http.StatusOK, tools.ECode{
		Message: "登陆成功",
	})
	return
}

func Logout(context *gin.Context) {
	//context.SetCookie("name", "", 3600, "/", "", true, false)
	//context.SetCookie("ID", "", 3600, "/", "", true, false)
	_ = model.FlushSession(context)
	context.Redirect(http.StatusFound, "/login")

}

// 新创建一个结构体
type CUser struct {
	Name      string `json:"name"`
	PassWord  string `json:"password"`
	PassWord2 string `json:"password_2"`
}

func CreateUser(context *gin.Context) {
	var user CUser
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}

	if user.Name == "" || user.PassWord == "" || user.PassWord2 == "" {
		context.JSON(http.StatusOK, tools.ParamErr)
		return
	}

	if user.PassWord != user.PassWord2 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不同",
		})
		return
	}
	//并发安全
	if oldeUser := model.GetUser(user.Name); oldeUser.Id > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户已经存在",
		})
		return
	}

	namelen := len(user.Name)
	password := len(user.PassWord)
	if namelen > 16 || namelen < 8 || password < 8 || password > 16 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "用户名或密码大于8位小于16位",
		})
		return
	}
	//判断密码是否为纯数字 正则表达式
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.PassWord) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "密码不能为纯数字",
		})
		return
	}
	newUser := model.User{
		Name:        user.Name,
		Password:    encrypt(user.PassWord),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Uuid:        tools.GetUUID(),
	}

	if err := model.CreateUser(&newUser); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10007,
			Message: "新用户创建失败！",
		})
		return
	}

	context.JSON(http.StatusOK, tools.OK)
	return

}

func encrypt(pwd string) string {
	newPwd := pwd + "ydsy"
	hash := md5.New()
	hash.Write([]byte(newPwd))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf("加密后的密码：%s\n", hashString)
	return hashString
}
