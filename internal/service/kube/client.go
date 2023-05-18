package kube

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"pack/internal/service"
)

type Clients struct {
	clientSet     kubernetes.Interface
	dynamicClient dynamic.Interface
}

func NewClients() (client Clients) {
	kubeConfig, _ := service.Config().ParseConfig(context.TODO(), "kube.config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return
	}

	client.clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return
	}

	client.dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return
	}

	return
}

func (c *Clients) ClientSet() kubernetes.Interface {
	return c.clientSet
}

func (c *Clients) DynamicClient() dynamic.Interface {
	return c.dynamicClient
}
