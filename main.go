package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"pack/internal/cmd"
	"pack/internal/logic/kube"
	"pack/internal/service"
)

func main() {
	stopCh := make(chan struct{})
	// initial docker client
	service.DockerSetUP()
	// initial harbor api http client
	service.HarborSetUp()
	// initial k8s share informer client
	_ = kube.Setup(stopCh)
	// set up cron job
	//service.CronSetUp()
	cmd.Main.Run(gctx.New())
}
