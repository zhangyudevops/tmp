package kube

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"time"
)

var shareInformerFactory informers.SharedInformerFactory

func NewSharedInformerFactory(stopCh <-chan struct{}) (err error) {
	clients := NewClients()
	g.Log().Info(context.TODO(), "k8s客户端成功初始化")
	shareInformerFactory = informers.NewSharedInformerFactory(clients.clientSet, time.Hour*24*2)

	// group version
	groupVersionResourceVars := []schema.GroupVersionResource{
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "apps", Version: "v1", Resource: "daemonsets"},
		{Group: "", Version: "v1", Resource: "namespaces"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "pods"},
	}

	for _, v := range groupVersionResourceVars {
		_, err = shareInformerFactory.ForResource(v)
		if err != nil {
			return err
		}
	}

	shareInformerFactory.Start(stopCh)
	shareInformerFactory.WaitForCacheSync(stopCh)
	return
}

func Get() informers.SharedInformerFactory {
	return shareInformerFactory
}

func Setup(stopCh <-chan struct{}) (err error) {
	err = NewSharedInformerFactory(stopCh)
	if err != nil {
		g.Log().Error(context.TODO(), err)
		return err
	}

	g.Log().Info(context.TODO(), "k8s informer成功初始化")
	return
}
