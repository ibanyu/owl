package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
)

func Login(userName, Pwd string) (string, error) {
	if err := LdapCheck(userName, Pwd); err != nil {
		return "", err
	}

	return GenerateToken(userName, Pwd)
}

type Claims struct {
	Username           string `json:"username"  binding:"required"`
	Password           string `json:"password"  binding:"required"`
	jwt.StandardClaims `json:"-"`
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(config.Conf.Login.TokenEffectiveHour) * time.Hour)
	claims := Claims{username, password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ipalfish-db-injection",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(config.Conf.Login.TokenSecret)
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Conf.Login.TokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
		return nil, fmt.Errorf("parse token failed, not a claims ins")
	}
	return nil, fmt.Errorf("get nil token claims")
}
