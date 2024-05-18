package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
	jwt "ququiz.org/lintang/quiz-query-service/biz/mw"
)

func Protected() []app.HandlerFunc {
	mwJwt := jwt.GetJwtMiddleware()
	mwJwt.MiddlewareInit()
	return []app.HandlerFunc{
		mwJwt.MiddlewareFunc(),
	}

}
