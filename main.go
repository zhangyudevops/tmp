package main

import (
	"github.com/gogf/gf/v2/os/gctx"
	"pack/internal/cmd"
	"pack/internal/service"
	"pack/internal/service/kube"
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
	service.CronSetUp()
	cmd.Main.Run(gctx.New())
}
