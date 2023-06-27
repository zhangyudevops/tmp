package v1

import "github.com/gogf/gf/v2/frame/g"

type RegisterUserReq struct {
	g.Meta   `path:"/user/register" method:"post"  tag:"user" summary:"register user"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
}

type RegisterUserRes struct{}
