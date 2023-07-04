package v1

import "github.com/gogf/gf/v2/frame/g"

type ListNamespaceReq struct {
	g.Meta `path:"/k8s/get-namespaces" method:"get" tag:"k8s" summary:"list namespaces"`
}

type ListNamespaceRes struct {
	List interface{} `json:"list"`
}

type ListDeployReq struct {
	g.Meta    `path:"/k8s/get-deploys" method:"get" tag:"k8s" summary:"list deployments"`
	Namespace string `json:"namespace"`
}

type ListDeployRes struct {
	List interface{} `json:"list"`
}

type ListPodsReq struct {
	g.Meta    `path:"/k8s/get-deploy-pods" method:"get" tag:"k8s" summary:"list deployment's pods"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

type ListPodsRes struct {
	List interface{} `json:"list"`
}

type CopyFileToPodReq struct {
	g.Meta    `path:"/k8s/copy-file-to-pod" method:"post" tag:"k8s" summary:"copy file or directory to pod"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Container string `json:"container"`
	InPath    string `json:"inPath"`
	OutPath   string `json:"outPath"`
}

type CopyFileToPodRes struct {
}
