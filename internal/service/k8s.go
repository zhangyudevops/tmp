package service

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"pack/internal/logic/kube"
	"pack/internal/model"
	"strings"
)

type sK8S struct{}

func K8S() *sK8S {
	return &sK8S{}
}

func (s *sK8S) GetNamespaces(ctx context.Context) (namespaces []*model.Namespace, err error) {
	// get namespaces
	ret, err := kube.Get().Core().V1().Namespaces().Lister().List(labels.Everything())
	if err != nil {
		return
	}

	var ns *model.Namespace
	for _, namespace := range ret {
		ns = &model.Namespace{
			Name:            namespace.Name,
			CreateTimeStamp: namespace.CreationTimestamp.Unix(),
			Status:          string(namespace.Status.Phase),
		}
		namespaces = append(namespaces, ns)
	}

	return
}

func (s *sK8S) GetDeploys(ctx context.Context, namespace string) (deploys []*model.Deployment, err error) {
	ret, err := kube.Get().Apps().V1().Deployments().Lister().
		Deployments(namespace).List(labels.Everything())
	if err != nil {
		return
	}

	var (
		de        *model.Deployment
		available string
	)
	for _, deploy := range ret {
		for _, condition := range deploy.Status.Conditions {
			if condition.Type == "Available" {
				available = string(condition.Status)
			}
		}
		de = &model.Deployment{
			Name:            deploy.Name,
			CreateTimeStamp: deploy.CreationTimestamp.Unix(),
			Namespace:       deploy.Namespace,
			Replicas:        deploy.Status.Replicas,
			UpdateReplicas:  deploy.Status.UpdatedReplicas,
			ReadyReplicas:   deploy.Status.ReadyReplicas,
			Available:       available,
			Labels:          deploy.Labels,
		}

		deploys = append(deploys, de)
	}

	return
}

// GetDeployPods 获取指定namespace下deployment的所有pods信息
func (s *sK8S) GetDeployPods(ctx context.Context, namespace string, label map[string]string) (pods []*model.Pod, err error) {
	// 组装labelSelector
	labelSet := labels.Set(label)
	labelSelector := labels.SelectorFromSet(labelSet)

	ret, err := kube.Get().Core().V1().Pods().Lister().Pods(namespace).List(labelSelector)
	if err != nil {
		return
	}

	var po *model.Pod
	for _, pod := range ret {
		po = &model.Pod{
			Name:            pod.Name,
			CreateTimeStamp: pod.CreationTimestamp.Unix(),
			Namespace:       pod.Namespace,
			Status:          string(pod.Status.Phase),
			HostIp:          pod.Status.HostIP,
			PodIp:           pod.Status.PodIP,
			NodeName:        pod.Spec.NodeName,
		}
		pods = append(pods, po)
	}

	err = s.CreateOrUpdateFromYamlFile(ctx, "")
	if err != nil {
		g.Dump(err)
	}
	return
}

// parseYamlFile 解析读取的yaml文件
func (s *sK8S) parseYamlFile(ctx context.Context, fileName string) (structures []*unstructured.Unstructured) {
	// test
	//fileName = "/Users/zhangsan/Documents/devops/go/pack/test/redis-sts.yaml"

	// 读取yaml文件内容
	yamlData, err := os.ReadFile(fileName)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	// yaml文件拆分多份
	yamlDecoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(string(yamlData)), 4096)

	for {
		structure := &unstructured.Unstructured{}
		err = yamlDecoder.Decode(structure)
		if err != nil {
			if err != io.EOF {
				g.Log().Error(ctx, err)
			} else {
				break
			}
		}
		structures = append(structures, structure)
	}

	return
}

// createDeployFromYamlFile 根据yaml文件内容创建deployment
func (s *sK8S) createDeployFromYamlFile(ctx context.Context, deploy *appsv1.Deployment) (err error) {
	_, err = kube.ClientSets.ClientSet().AppsV1().Deployments(deploy.Namespace).
		Create(ctx, deploy, metav1.CreateOptions{})
	if err != nil {
		return
	}

	return
}

// createStatefulSetFromYamlFile 根据yaml文件内容创建statefulSet
func (s *sK8S) createStatefulSetFromYamlFile(ctx context.Context, sts *appsv1.StatefulSet) (err error) {
	_, err = kube.ClientSets.ClientSet().AppsV1().StatefulSets(sts.Namespace).
		Create(ctx, sts, metav1.CreateOptions{})
	if err != nil {
		return
	}

	return
}

// CreateConfigMapFromYamlFile 根据yaml文件内容创建configmap
func (s *sK8S) createConfigMapFromYamlFile(ctx context.Context, cm *corev1.ConfigMap) (err error) {
	_, err = kube.ClientSets.ClientSet().CoreV1().ConfigMaps(cm.Namespace).
		Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		return
	}

	return
}

// CreateServiceFromYamlFile 从解析后的内容创建service
func (s *sK8S) createServiceFromYamlFile(ctx context.Context, service *corev1.Service) (err error) {
	_, err = kube.ClientSets.ClientSet().CoreV1().Services(service.Namespace).
		Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return
	}

	return
}

// createPodFromYamlFile 创建pod
func (s *sK8S) createPodFromYamlFile(ctx context.Context, pod *corev1.Pod) (err error) {
	_, err = kube.ClientSets.ClientSet().CoreV1().Pods(pod.Namespace).
		Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return
	}

	return
}

// selectSourceType 根据kind类型确定创建对应资源
func (s *sK8S) selectSourceTypeCreateOrUpdate(ctx context.Context, structure *unstructured.Unstructured) (err error) {
	converter := runtime.DefaultUnstructuredConverter
	switch structure.GetKind() {
	case "StatefulSet":
		statefulSet := &appsv1.StatefulSet{}
		if err = converter.FromUnstructured(structure.Object, statefulSet); err != nil {
			return
		}
		if s.getSourceResult(ctx, structure.GetNamespace(), "StatefulSet", structure.GetName()) != true {
			if err = s.createStatefulSetFromYamlFile(ctx, statefulSet); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s不存在，创建成功", structure.GetName())
		} else {
			if err = s.updateStatefulSetFromYamlFile(ctx, statefulSet); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s不存在，创建成功", structure.GetName())
		}
	case "Deployment":
		deployment := &appsv1.Deployment{}
		if err = converter.FromUnstructured(structure.Object, deployment); err != nil {
			return
		}
		if s.getSourceResult(ctx, structure.GetNamespace(), "Deployment", structure.GetName()) != true {
			if err = s.createDeployFromYamlFile(ctx, deployment); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s不存在，创建成功", structure.GetName())
		} else {
			if err = s.updateDeployFromYamlFile(ctx, deployment); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s已经存在，升级成功", structure.GetName())
		}
	case "ConfigMap":
		configMap := &corev1.ConfigMap{}
		if err = converter.FromUnstructured(structure.Object, configMap); err != nil {
			return
		}
		if s.getSourceResult(ctx, structure.GetNamespace(), "ConfigMap", structure.GetName()) != true {
			if err = s.createConfigMapFromYamlFile(ctx, configMap); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s不存在，创建成功", structure.GetName())
		} else {
			if err = s.updateConfigMapFromYamlFile(ctx, configMap); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s已经存在，升级成功", structure.GetName())
		}
	case "Service":
		service := &corev1.Service{}
		if err = converter.FromUnstructured(structure.Object, service); err != nil {
			return
		}
		if s.getSourceResult(ctx, structure.GetNamespace(), "Service", structure.GetName()) != true {
			if err = s.createServiceFromYamlFile(ctx, service); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s不存在，创建成功", structure.GetName())
		} else {
			if err = s.updateServiceFromYamlFile(ctx, service); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s已经存在，升级成功", structure.GetName())
		}
	case "Pod":
		pod := &corev1.Pod{}
		if err = converter.FromUnstructured(structure.Object, pod); err != nil {
			return
		}
		if s.getSourceResult(ctx, structure.GetNamespace(), "Pod", structure.GetName()) != true {
			if err = s.createPodFromYamlFile(ctx, pod); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s不存在，创建成功", structure.GetName())
		} else {
			if err = s.updatePodFromYamlFile(ctx, pod); err != nil {
				return
			}
			g.Log().Debugf(ctx, "%s已经存在，升级成功", structure.GetName())
		}
	default:
		return errors.New("类型不存在")
	}

	return
}

// CreateOrUpdateFromYamlFile 读取yaml文件，并解析对应yaml，并创建
// 先判断对应资源是否存在，存在则升级，不存在则创建
func (s *sK8S) CreateOrUpdateFromYamlFile(ctx context.Context, fileName string) (err error) {
	fileName = "/Users/zhangsan/Documents/devops/go/pack/test/busybox-pod.yaml"
	structure := s.parseYamlFile(ctx, fileName)
	for _, st := range structure {
		if st.GetNamespace() == "" {
			st.SetNamespace("default")
		}
		if err = s.selectSourceTypeCreateOrUpdate(ctx, st); err != nil {
			return
		}
	}

	return
}

// updateServiceFromYamlFile 从解析后的内容升级service
func (s *sK8S) updateServiceFromYamlFile(ctx context.Context, service *corev1.Service) (err error) {
	_, err = kube.ClientSets.ClientSet().CoreV1().Services(service.Namespace).
		Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return
	}

	return
}

// UpdatePodFromYamlFile 升级pod
func (s *sK8S) updatePodFromYamlFile(ctx context.Context, pod *corev1.Pod) (err error) {
	_, err = kube.ClientSets.ClientSet().CoreV1().Pods(pod.Namespace).
		Update(ctx, pod, metav1.UpdateOptions{})
	if err != nil {
		return
	}

	return
}

// updateDeployFromYamlFile 根据yaml文件内容升级deployment
func (s *sK8S) updateDeployFromYamlFile(ctx context.Context, deploy *appsv1.Deployment) (err error) {
	_, err = kube.ClientSets.ClientSet().AppsV1().Deployments(deploy.Namespace).
		Update(ctx, deploy, metav1.UpdateOptions{})
	if err != nil {
		return
	}

	return
}

// updateStatefulSetFromYamlFile 根据yaml文件内容升级statefulSet
func (s *sK8S) updateStatefulSetFromYamlFile(ctx context.Context, sts *appsv1.StatefulSet) (err error) {
	_, err = kube.ClientSets.ClientSet().AppsV1().StatefulSets(sts.Namespace).
		Update(ctx, sts, metav1.UpdateOptions{})
	if err != nil {
		return
	}

	return
}

// updateConfigMapFromYamlFile 根据yaml文件内容升级configmap
func (s *sK8S) updateConfigMapFromYamlFile(ctx context.Context, cm *corev1.ConfigMap) (err error) {
	_, err = kube.ClientSets.ClientSet().CoreV1().ConfigMaps(cm.Namespace).
		Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		return
	}

	return
}

// getSourceResult 验证当前服务是否存在
func (s *sK8S) getSourceResult(ctx context.Context, namespace, kind, name string) (result bool) {
	switch kind {
	case "Service":
		_, err := kube.ClientSets.ClientSet().CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	case "Deployment":
		_, err := kube.ClientSets.ClientSet().AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		return true
	case "StatefulSet":
		_, err := kube.ClientSets.ClientSet().AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		return true
	case "ConfigMap":
		_, err := kube.ClientSets.ClientSet().CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		return true
	case "Pod":
		_, err := kube.ClientSets.ClientSet().CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		return true
	}

	return
}

// copyExecutor
// args: in 需要拷贝的目录路径或者文件路径
// args: out pod里面的目录路径
func (s *sK8S) copyExecutor(ctx context.Context, containerInfo *model.Execute, in, out string) (executor *rest.Request) {
	execRequest := kube.ClientSets.ClientSet().CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(containerInfo.Namespace).
		Name(containerInfo.PodName).
		SubResource("exec")

	// 根据需要拷贝的文件或者目录组装参数
	if gfile.IsDir(in) {
		execRequest.Param("container", containerInfo.ContainerName)
		execRequest.Param("command", "cp")
		execRequest.Param("command", "-r")
		execRequest.Param("command", in+"/*")
		execRequest.Param("command", out)
	} else {
		execRequest.Param("container", containerInfo.ContainerName)
		execRequest.Param("command", "cp")
		execRequest.Param("command", in)
		execRequest.Param("command", out)
	}

	return execRequest
}

// CopyFileToPod
// args: in 需要拷贝的目录路径或者文件路径
// args: out pod里面的目录路径
func (s *sK8S) CopyFileToPod(ctx context.Context, containerInfo *model.Execute, in, out string) (err error) {
	execRequest := s.copyExecutor(ctx, containerInfo, in, out)

	// 设置执行选项
	execOptions := &corev1.PodExecOptions{
		Command: []string{"sh"},
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}

	// 序列化
	execRequest.VersionedParams(
		execOptions,
		scheme.ParameterCodec,
	)

	// 执行
	exec, err := remotecommand.NewSPDYExecutor(kube.NewConfig(), "POST", execRequest.URL())
	if err != nil {
		return
	}
	if err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	}); err != nil {
		return
	}

	return
}
