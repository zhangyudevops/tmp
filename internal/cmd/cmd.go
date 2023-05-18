package cmd

import (
	"context"
	"pack/internal/controller/clean"
	"pack/internal/controller/docker"
	"pack/internal/controller/file"
	"pack/internal/controller/harbor"
	"pack/internal/controller/pack"
	"pack/internal/controller/path"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					file.File(),
					docker.Docker(),
					path.Path(),
					harbor.Harbor(),
					clean.Clean(),
					pack.Pack(),
				)
			})
			s.Run()
			return nil
		},
	}
)
