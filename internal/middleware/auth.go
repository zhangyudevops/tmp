package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"pack/internal/service"
)

// JWTAuth user auth
func JWTAuth(r *ghttp.Request) {
	token := r.Get("token")
	if _, err := service.Jwt().ParseToken(token.String()); err != nil {
		r.Response.WriteStatus(401)
	}

	r.Middleware.Next()
}
