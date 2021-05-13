package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/code"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/service/admin"
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
			c.JSON(http.StatusOK, "please login")
			c.Abort()
			return
		}

		if err := claims.Valid(); err != nil {
			c.JSON(http.StatusOK, "place login")
			c.Abort()
			return
		}

		c.Set("user", claims.Username)
		c.Next()
	}
}

func OnlyDbaOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(string)
		isDba, err := task.IsDba(user)
		if err != nil {
			c.JSON(http.StatusOK, Resp{
				Code:    code.ParamInvalid,
				Message: fmt.Sprintf("check dba auth failed, %s", err.Error()),
			})
			c.Abort()
			return
		}

		isAdmin, err := admin.IsAdmin(user)
		if err != nil {
			c.JSON(http.StatusOK, Resp{
				Code:    code.ParamInvalid,
				Message: fmt.Sprintf("check admin auth failed, %s", err.Error()),
			})
			c.Abort()
			return
		}
		if !isDba && !isAdmin {
			c.JSON(http.StatusOK, Resp{
				Code:    code.InternalErr,
				Message: "only dba and admin can operate cluster and administrator module",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
