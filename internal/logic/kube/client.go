package kube

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Clients struct {
	clientSet     kubernetes.Interface
	dynamicClient dynamic.Interface
}

func NewConfig() (config *rest.Config) {
	kubeConfig, _ := g.Config().Get(context.TODO(), "kube.config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig.String())
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return
	}
	return
}

func NewClients() (client Clients) {
	var err error
	config := NewConfig()
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
