package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 跨域中间件配置

// Cors 跨域处理中间件
func Cors() gin.HandlerFunc {
	// 返回一个处理跨域请求的中间件
	return func(c *gin.Context) {
		// 获取请求的方法
		method := c.Request.Method

		// 设置允许所有来源跨域
		c.Header("Access-Control-Allow-Origin", "*")

		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")

		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")

		// 设置在浏览器中可访问的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		// 允许发送Cookie
		c.Header("Access-Control-Allow-Credentials", "true")

		// 如果请求方法是OPTIONS，处理预检请求
		if method == "OPTIONS" {
			// 中止请求，返回HTTP状态码 http.StatusNoContent
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求，执行链中的下一个处理程序
		c.Next()
	}
}
