package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
	"vote/app/model"
	"vote/app/tools"
)

// 逻辑层
type User struct {
	Name         string `json:"name"`
	Password     string `json:"password"`
	CaptchaId    string `json:"captcha_id"`
	CaptchaValue string `json:"captcha_value"`
}

// // GetLogin godoc
// // @Summary 获取用户登录页面
// // @Description 获取用户登录页面
// // @Tags login
// // @Produce html
// // @Success 200
// // @Router /login [get]
func GetLogin(context *gin.Context) {
	//
	context.HTML(http.StatusOK, "login.html", nil)
}

// DoLogin godoc
// @Summary 执行用户登录
// @Description 执行用户登录操作
// @Tags login
// @Accept json
// @Produce json
// @Param name body string true "用户名"
// @Param password body string true "密码"
// @Param captcha_id body string true "验证码ID"
// @Param captcha_value body string true "验证码值"
// @Success 200 {object} tools.ECode "成功响应"
// @Failure 400 {object} tools.ECode "请求失败"
// @Failure 401 {object} tools.ECode "未认证"
// @Router /login [post]
func DoLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}
	//if !tools.CaptchaVerify(tools.CaptchaData{
	//	CaptchaId: user.CaptchaId,
	//	Data:      user.CaptchaValue,
	//}) {
	//	context.JSON(http.StatusOK, tools.CacheErr)
	//	return
	//}
	ret := model.GetUserV1(user.Name)

	if ret.Id < 1 || ret.Password != encrypt(user.Password) {
		context.JSON(http.StatusOK, tools.UserErr)
		return
	}
	token, _ := model.GetJwt(ret.Id, user.Name)
	//context.SetCookie("name", user.Name, 3600, "/", "", true, false)
	//context.SetCookie("Id", fmt.Sprintf(strconv.FormatInt(ret.Id, 10)), 3600, "/", "", true, false)
	////
	//
	//sessionInfo := model.GetSession(context)
	//responseData := map[string]interface{}{
	//	"name": sessionInfo["name"],
	//	"id":   sessionInfo["id"],
	//}
	//fmt.Printf("Session Info: %+v\n", responseData)
	context.Set("token", token)
	context.JSON(http.StatusOK, tools.ECode{
		Message: "登陆成功",
		Data:    token,
	})
	return
}

// Logout godoc
// @Summary 执行用户退出
// @Description 执行用户退出
// @Tags login
// @Success 302
// @Router /logout [get]
func Logout(context *gin.Context) {
	//context.SetCookie("name", "", 3600, "/", "", true, false)
	//context.SetCookie("ID", "", 3600, "/", "", true, false)
	//_ = model.FlushSession(context)
	model.BlacklistToken(context)
	context.Redirect(http.StatusFound, "/login")
}

// 新创建一个结构体
type CUser struct {
	Name      string `json:"name"`
	PassWord  string `json:"password"`
	PassWord2 string `json:"password_2"`
}

// CreateUser godoc
// @Summary 创建一个新用户
// @Description 创建一个新用户
// @Tags login
// @Accept json
// @Produce json
// @Param name body string true "用户名"
// @Param password body string true "密码"
// @Param password2 body string true "确认密码"
// @Success 200 {object} tools.ECode
// @Router /logout [post]
func CreateUser(context *gin.Context) {
	var user CUser
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	if user.Name == "" || user.PassWord == "" || user.PassWord2 == "" {
		context.JSON(http.StatusOK, tools.ParamErr)
		return
	}

	if user.PassWord != user.PassWord2 {
		context.JSON(http.StatusOK, tools.PassWordErr)
		return
	}
	//并发安全
	if oldeUser := model.GetUserV1(user.Name); oldeUser.Id > 0 {
		context.JSON(http.StatusOK, tools.UserExits)
		return
	}

	namelen := len(user.Name)
	password := len(user.PassWord)
	if namelen > 16 || namelen < 8 || password < 8 || password > 16 {
		context.JSON(http.StatusOK, tools.UserOver)
		return
	}
	//判断密码是否为纯数字 正则表达式
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.PassWord) {
		context.JSON(http.StatusOK, tools.PassWordNum)
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
		context.JSON(http.StatusOK, tools.UserCreateFail)
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
