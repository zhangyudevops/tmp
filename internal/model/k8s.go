package model

type Namespace struct {
	Name            string `json:"name"`
	CreateTimeStamp int64  `json:"createTimeStamp"`
	Status          string `json:"status"`
}

type Deployment struct {
	Name            string            `json:"name"`
	CreateTimeStamp int64             `json:"createTimeStamp"`
	Namespace       string            `json:"namespace"`
	Replicas        int32             `json:"replicas"`
	UpdateReplicas  int32             `json:"updateReplicas"`
	ReadyReplicas   int32             `json:"readyReplicas"`
	Available       string            `json:"available"`
	Labels          map[string]string `json:"labels"`
}

type Pod struct {
	Name            string `json:"name"`
	CreateTimeStamp int64  `json:"createTimeStamp"`
	Namespace       string `json:"namespace"`
	Status          string `json:"status"`
	HostIp          string `json:"hostIp"`
	PodIp           string `json:"podIp"`
	NodeName        string `json:"nodeName"`
}

type Execute struct {
	Namespace     string `json:"namespace"`
	PodName       string `json:"podName"`
	ContainerName string `json:"containerName"`
}
