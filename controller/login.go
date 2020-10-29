package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/auth"
)

func Login(ctx *gin.Context) Resp {
	f := "Login()-->"
	var authParam auth.Claims
	if err := ctx.BindJSON(&authParam); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), Code: code.ParamInvalid}
	}

	token, err := auth.Login(authParam.Username, authParam.Password)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s, login failed :%s", f, err.Error()), Code: code.ParamInvalid}
	}
	return Resp{Data: token}
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		claims, err := auth.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusPermanentRedirect, Resp{
				Code:    code.Redirect,
				Message: "check token failed, place login",
				Data:    config.Conf.Login.LoginPath,
			})
			c.Abort()
			return
		}

		c.Set("user", claims.Username)
		c.Next()
	}
}
