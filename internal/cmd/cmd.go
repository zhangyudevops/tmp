package cmd

import (
	"context"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"pack/internal/controller/clean"
	"pack/internal/controller/docker"
	"pack/internal/controller/file"
	"pack/internal/controller/harbor"
	"pack/internal/controller/pack"
	"pack/internal/controller/path"
	"pack/internal/controller/user"
	"pack/internal/dao"
	"pack/internal/model/entity"
	"pack/utility/util"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// 启动gToken
			gfToken := &gtoken.GfToken{
				ServerName:       "pack",
				LoginPath:        "/login",
				LoginBeforeFunc:  loginFunc,
				LoginAfterFunc:   loginAfter,
				LogoutPath:       "/user/logout",
				AuthExcludePaths: g.SliceStr{"/user/register"}, // 不拦截路径
				MultiLogin:       true,
			}

			// 认证接口
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse, CORS)
				err := gfToken.Middleware(ctx, group)
				if err != nil {
					panic(err)
				}
				group.Bind(
					file.File(),
					docker.Docker(),
					path.Path(),
					harbor.Harbor(),
					clean.Clean(),
					pack.Pack(),
					user.User(),
				)
			})
			s.Run()
			return nil
		},
	}
)

// 登陆函数
func loginFunc(r *ghttp.Request) (string, interface{}) {
	username := r.Get("username").String()
	password := r.Get("password").String()
	ctx := context.TODO()

	if username == "" || password == "" {
		r.Response.WriteJson(gtoken.Fail("账号或密码错误."))
		r.ExitAll()
	}

	// query userinfo from user table
	userInfo := entity.User{}
	err := dao.User.Ctx(ctx).Where("username = ?", username).Scan(&userInfo)
	if err != nil {
		r.Response.WriteJson(gtoken.Fail("账号或密码错误."))
		r.ExitAll()
	}

	if util.HashPassword(password, userInfo.Salt) != userInfo.Password {
		r.Response.WriteJson(gtoken.Fail("账号或密码错误."))
		r.ExitAll()
	}

	// 唯一标识，扩展参数user data
	return username, userInfo
}

// loginAfter 用户登录后返回数据
func loginAfter(r *ghttp.Request, respData gtoken.Resp) {
	if !respData.Success() {
		respData.Code = 0
		r.Response.WriteJson(respData)
	} else {
		respData.Code = 0
		r.Response.WriteJson(respData)
	}
	return
}

// CORS 跨域
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
