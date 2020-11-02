package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/auth"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/task"
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

		if err := claims.Valid(); err != nil {
			c.JSON(http.StatusPermanentRedirect, Resp{
				Code:    code.Redirect,
				Message: "token expired, place login",
				Data:    config.Conf.Login.LoginPath,
			})
			c.Abort()
			return
		}

		c.Set("user", claims.Username)
		c.Next()
	}
}

func OnlyDba() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(string)
		isDba, err := task.IsDba(user)
		if err != nil {
			c.JSON(http.StatusOK, Resp{
				Code:    code.InternalErr,
				Message: fmt.Sprintf("check dba auth failed, %s", err.Error()),
			})
			c.Abort()
			return
		}

		if !isDba {
			c.JSON(http.StatusOK, Resp{
				Code:    code.InternalErr,
				Message: "only dba can operate cluster",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
