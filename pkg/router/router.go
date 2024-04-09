package router

import (
	"go-app/pkg/handler"
	"net/http"

	"github.com/BoyYangZai/go-server-lib/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		// 允许所有来源访问
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许以下方法访问
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许以下头部字段访问
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		// 允许客户端发送 Cookie
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 如果是 OPTIONS 请求，结束请求处理
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}

		// 继续处理其他请求
		ctx.Next()
	})

	user := router.Group("/user")
	{
		user.POST("/verify-code", handler.VerifyCode)
		user.POST("/registry", handler.Registry)
		user.POST("/login", handler.Login)
		user.GET("/list", handler.List).Use(jwt.AuthMiddleware())
	}

	auth_test := router.Group("/auth-test")
	{
		auth_test.Use(jwt.AuthMiddleware())
		auth_test.GET("/", handler.Submit)
	}

	router.Run(":8080")

	return router
}
