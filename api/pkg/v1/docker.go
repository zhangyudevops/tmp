package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DockerTestReq struct {
	g.Meta `path:"/docker/test" method:"post"`
	Image  string `json:"image" dc:"test"`
}

type DockerTestRes struct {
	List interface{} `json:"list" dc:"test"`
}
