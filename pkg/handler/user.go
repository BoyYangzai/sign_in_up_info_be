package handler

import (
	"go-app/pkg/service"
	"net/http"
	"time"

	"github.com/BoyYangZai/go-server-lib/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type VerifyCodeRquest struct {
	Email string `json:"email"`
}

func VerifyCode(c *gin.Context) {
	var requestBody VerifyCodeRquest

	// 通过 ShouldBindJSON 解析 JSON 请求体到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code := generateVerificationCode(6)
	user := "1484502768@qq.com"
	password := "jyderqttsuyyiagf"
	host := "smtp.qq.com:587"
	to := requestBody.Email
	subject := "verifycode:"
	body := `
	<html>
	<body>
	<h3>
	` + code +
		`
	</h3>
	</body>
	</html>
	`
	println("send email")
	err := SendMail(user, password, host, to, subject, body, "html")
	if err != nil {
		println("send mail error!")
		println(err)
	} else {
		println("send mail success!")
	}

	service.UpdateVerifyCode(to, code)
	c.JSON(http.StatusOK, gin.H{
		"msg": "verifyCode sent",
	})
}

type RegistryRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Roles    string `json:"roles"`
	Position string `json:"position"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
}

func Registry(c *gin.Context) {
	var requestBody RegistryRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//如果username已经存在
	if service.CheckUsernameIsExisted(requestBody.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	//检查所有参数都不能为空
	if requestBody.Username == "" || requestBody.Password == "" || requestBody.Roles == "" || requestBody.Position == "" || requestBody.Age == "" || requestBody.Gender == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all parameters must not be empty"})
		return
	}

	user := service.User{
		Username:    requestBody.Username,
		Password:    requestBody.Password,
		Roles:       requestBody.Roles, // 将拆分后的字符串切片赋值给 user.Roles
		Position:    requestBody.Position,
		Age:         requestBody.Age,
		Gender:      requestBody.Gender,
		CreatedTime: time.Now(),
	}
	service.InitUser(user)
	c.JSON(http.StatusOK, gin.H{
		"msg": "registry success",
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var requestBody LoginRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isMatched, user := service.MatchEmailAndKey(requestBody.Username, requestBody.Password, "Password")
	if !isMatched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is wrong"})
		return
	}
	jwt.Auth(c, isMatched, user.Username, user.ID)

}

func List(c *gin.Context) {
	users := service.List()
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
