package tools

import (
	"github.com/mojocn/base64Captcha"
)

type CaptchaData struct {
	CaptchaId string `json:"captcha_id"`
	Data      string `json:"data"`
}

type driverString struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var digitDriver = base64Captcha.DriverDigit{
	Height:   50,  //生成图片高度
	Width:    150, //生成图片宽度
	Length:   5,   //验证码长度
	MaxSkew:  1,   //文字的倾斜度 越大倾斜越狠，越不容易看懂
	DotCount: 1,   //背景的点数，越大，字体越模糊
}

// 使用内存驱动
var store = base64Captcha.DefaultMemStore

func CaptchaGenerate() (CaptchaData, error) {
	var ret CaptchaData
	c := base64Captcha.NewCaptcha(&digitDriver, store)
	id, b64s, _, err := c.Generate()
	if err != nil {
		return ret, err
	}
	ret.CaptchaId = id
	ret.Data = b64s
	return ret, nil
}

func CaptchaVerify(data CaptchaData) bool {
	return store.Verify(data.CaptchaId, data.Data, true)
}
