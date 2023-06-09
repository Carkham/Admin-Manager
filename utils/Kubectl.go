package utils

import (
	"admin/conf"
	"admin/model"
	"context"
	"encoding/json"
	"fmt"
	appV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	networkV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"path"
	"strconv"
	"time"
)

var KubeClient *kubernetes.Clientset
var basePath string

const (
	MountRootKey      = "INIT_MOUNT_ROOT"
	SourceTypeKey     = "INIT_FETCH_TYPE"
	SourceLocationKey = "INIT_SRC_LOC"
	WorkDirKey        = "INIT_WORKDIR"
	FuncIDLabelKey    = "funcID"
	CodeFolder        = "code"
	DataFolder        = "data"
	CpuQuotaKey       = "cpu"
	MemQuotaKey       = "memory"
	GpuQuotaKey       = "nvidia.com/gpu"
	DeployNameKey     = "app"
	PrefixPathTypeStr = "Prefix"
)

var PrefixType = networkV1.PathType(PrefixPathTypeStr)

func KubeClientInit(projectPath string) {
	basePath = conf.Config.FileSystem.RootPath
	config, err := clientcmd.BuildConfigFromFlags("", projectPath+"/conf/admin.conf")
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	KubeClient = clientset
}

func PrepareInitContainer(funcInfo *model.FuncInfo) []coreV1.Container {
	return []coreV1.Container{
		{
			Name:  "project-initializer",
			Image: "realssd/pod_initializer",
			VolumeMounts: []coreV1.VolumeMount{
				{
					Name:      "workspace",
					MountPath: "/workspace",
				},
				{
					Name:      "ucode",
					MountPath: "/ucode",
				},
			},
			Env: []coreV1.EnvVar{
				{
					Name:  MountRootKey,
					Value: "/mnt",
				},
				{
					Name:  SourceTypeKey,
					Value: funcInfo.SourceType,
				},
				{
					Name:  SourceLocationKey,
					Value: funcInfo.SourceLocation,
				},
				{
					Name:  WorkDirKey,
					Value: "/workspace",
				},
			},
		},
	}
}

func CreateServiceAndIngress(funcInfo *model.FuncInfo) error {
	depLabel := fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel)
	fidStr := strconv.FormatInt(funcInfo.FunctionID, 10)
	uidStr := strconv.FormatInt(funcInfo.UserID, 10)
	selectLabels := map[string]string{
		"app":    depLabel,
		"userID": uidStr,
		"funcID": fidStr,
	}
	svc := coreV1.Service{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel),
			Labels: selectLabels,
		},
		Spec: coreV1.ServiceSpec{
			Selector: map[string]string{
				"app": fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel),
			},
			Ports: []coreV1.ServicePort{
				{
					Name:       "tcp",
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
					Protocol:   "TCP",
				},
			},
		},
	}
	_, err := KubeClient.CoreV1().Services(coreV1.NamespaceDefault).Create(
		context.Background(),
		&svc,
		metaV1.CreateOptions{},
	)

	if err != nil {
		return err
	}

	ing := networkV1.Ingress{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:   fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel),
			Labels: map[string]string{},
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
			},
		},
		Spec: networkV1.IngressSpec{
			Rules: []networkV1.IngressRule{
				{
					"", networkV1.IngressRuleValue{
						HTTP: &networkV1.HTTPIngressRuleValue{
							Paths: []networkV1.HTTPIngressPath{
								{
									PathType: &PrefixType,
									Path:     fmt.Sprintf("/call/%s/%s", funcInfo.UserName, funcInfo.FuncLabel),
									Backend: networkV1.IngressBackend{
										Service: &networkV1.IngressServiceBackend{
											Name: fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel),
											Port: networkV1.ServiceBackendPort{
												Number: int32(8080),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = KubeClient.NetworkingV1().Ingresses(coreV1.NamespaceDefault).Create(
		context.Background(),
		&ing,
		metaV1.CreateOptions{},
	)

	if err != nil {
		return err
	}
	return nil
}

func DeleteServiceAndIngress(name string) error {

	_ = KubeClient.NetworkingV1().Ingresses(coreV1.NamespaceDefault).Delete(
		context.Background(),
		name,
		metaV1.DeleteOptions{},
	)

	_ = KubeClient.CoreV1().Services(coreV1.NamespaceDefault).Delete(
		context.Background(),
		name,
		metaV1.DeleteOptions{},
	)
	return nil
}

func PrepareMainContainer(funcInfo *model.FuncInfo) []coreV1.Container {
	cpuLimit, _ := resource.ParseQuantity(fmt.Sprintf("%dm", funcInfo.CPUQuotaM[1]))
	memLimit, _ := resource.ParseQuantity(fmt.Sprintf("%dMi", funcInfo.MemQuotaMi[1]))
	limit := coreV1.ResourceList{
		CpuQuotaKey: cpuLimit,
		MemQuotaKey: memLimit,
	}

	cpuRequest, _ := resource.ParseQuantity(fmt.Sprintf("%dm", funcInfo.CPUQuotaM[0]))
	memRequest, _ := resource.ParseQuantity(fmt.Sprintf("%dMi", funcInfo.MemQuotaMi[0]))
	request := coreV1.ResourceList{
		CpuQuotaKey: cpuRequest,
		MemQuotaKey: memRequest,
	}

	if funcInfo.GPUQuota > 0 {
		gpuQuota, _ := resource.ParseQuantity(strconv.Itoa(funcInfo.GPUQuota))
		limit[GpuQuotaKey] = gpuQuota
	}

	return []coreV1.Container{
		{
			Name:  "service-main",
			Image: funcInfo.ImageName,
			VolumeMounts: []coreV1.VolumeMount{
				{
					Name:      "workspace",
					MountPath: "/workspace",
				},
				{
					Name:      "udata",
					MountPath: "/udata",
				},
			},
			Resources: coreV1.ResourceRequirements{
				Limits:   limit,
				Requests: request,
			},
		},
	}
}

func PreparePodTemplateSpec(funcInfo *model.FuncInfo, labels map[string]string) coreV1.PodTemplateSpec {
	ucodePath := path.Join(basePath, CodeFolder, strconv.FormatInt(funcInfo.UserID, 10))
	udataPath := path.Join(basePath, DataFolder, strconv.FormatInt(funcInfo.UserID, 10))
	nfsServer := conf.Config.FileSystem.NFSAddr
	return coreV1.PodTemplateSpec{
		ObjectMeta: metaV1.ObjectMeta{
			Labels: labels,
		},
		Spec: coreV1.PodSpec{
			InitContainers: PrepareInitContainer(funcInfo),
			Containers:     PrepareMainContainer(funcInfo),
			ResourceClaims: []coreV1.PodResourceClaim{},
			RestartPolicy:  "Always",
			Volumes: []coreV1.Volume{
				{
					"workspace",
					coreV1.VolumeSource{
						EmptyDir: &coreV1.EmptyDirVolumeSource{},
					},
				},
				{
					"ucode",
					coreV1.VolumeSource{
						NFS: &coreV1.NFSVolumeSource{
							Server:   nfsServer,
							Path:     ucodePath,
							ReadOnly: true,
						},
					},
				},
				{
					"udata",
					coreV1.VolumeSource{
						NFS: &coreV1.NFSVolumeSource{
							Server: nfsServer,
							Path:   udataPath,
						},
					},
				},
			},
		},
	}
}

func PrepareCreateDeployment(funcInfo *model.FuncInfo) *appV1.Deployment {
	depLabel := fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel)
	fidStr := strconv.FormatInt(funcInfo.FunctionID, 10)
	uidStr := strconv.FormatInt(funcInfo.UserID, 10)

	selectLabels := map[string]string{
		"app":    depLabel,
		"userID": uidStr,
		"funcID": fidStr,
	}

	ans := appV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:   depLabel,
			Labels: selectLabels,
		},
		Spec: appV1.DeploymentSpec{
			Replicas: &funcInfo.PodCount,
			Selector: &metaV1.LabelSelector{
				MatchLabels: selectLabels,
			},
			Template: PreparePodTemplateSpec(funcInfo, selectLabels),
		},
	}
	return &ans
}

// StopDeployment 函数必需UserName、FuncLabel、TrigType三个参数
func StopDeployment(funcInfo *model.FuncInfo) error {
	depName := fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel)
	switch funcInfo.TrigType {
	case "CRON":
		cronCli := KubeClient.BatchV1().CronJobs(coreV1.NamespaceDefault)
		err := cronCli.Delete(context.Background(), depName, metaV1.DeleteOptions{})
		if err != nil {
			log.Printf("[Create CronJob] Create CronJob Error: %s", err.Error())
			return err
		}
	default:
		depCli := KubeClient.AppsV1().Deployments(coreV1.NamespaceDefault)
		err := depCli.Delete(context.Background(), depName, metaV1.DeleteOptions{})
		if err != nil {
			log.Printf("[Create Deployment] Create Deployment Error: %s", err.Error())
			return err
		}
		err = DeleteServiceAndIngress(depName)
		if err != nil {
			return err
		}
	}

	return nil
}

func StartDeployment(funcInfo *model.FuncInfo) (depName string, err error) {

	switch funcInfo.TrigType {
	case "CRON":
		cronCli := KubeClient.BatchV1().CronJobs(coreV1.NamespaceDefault)
		cron := PrepareCronJob(funcInfo)
		result, err := cronCli.Create(context.Background(), cron, metaV1.CreateOptions{})
		if err != nil {
			log.Printf("[Create CronJob] Create CronJob Error: %s", err.Error())
			return "", err
		}
		depName = result.Name
	default:
		depCli := KubeClient.AppsV1().Deployments(coreV1.NamespaceDefault)
		dep := PrepareCreateDeployment(funcInfo)
		result, err := depCli.Create(context.Background(), dep, metaV1.CreateOptions{})
		if err != nil {
			log.Printf("[Create Deployment] Create Deployment Error: %s", err.Error())
			return "", err
		}
		err = CreateServiceAndIngress(funcInfo)
		depName = result.Name
		if err != nil {
			return depName, err
		}
	}

	return
}

func PrepareCronJob(funcInfo *model.FuncInfo) *batchV1.CronJob {
	depLabel := fmt.Sprintf("%s-%s", funcInfo.UserName, funcInfo.FuncLabel)
	fidStr := strconv.FormatInt(funcInfo.FunctionID, 10)
	uidStr := strconv.FormatInt(funcInfo.UserID, 10)
	selectLabels := map[string]string{
		"app":    depLabel,
		"userID": uidStr,
		"funcID": fidStr,
	}
	ans := batchV1.CronJob{
		ObjectMeta: metaV1.ObjectMeta{
			Name:   depLabel,
			Labels: selectLabels,
		},
		Spec: batchV1.CronJobSpec{
			Schedule: funcInfo.TimeStr,
			TimeZone: &conf.Config.Service.TimeZone,
			JobTemplate: batchV1.JobTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: selectLabels,
				},
				Spec: batchV1.JobSpec{
					Parallelism: &funcInfo.PodCount,
					Completions: &funcInfo.PodCount,
					Template:    PreparePodTemplateSpec(funcInfo, selectLabels),
				},
			},
		},
	}
	return &ans
}

func GetPodInfoList(funcID int64) (list map[string][]model.PodInfo, err error) {
	list = make(map[string][]model.PodInfo)
	podList, err := KubeClient.CoreV1().Pods(coreV1.NamespaceDefault).List(
		context.Background(),
		metaV1.ListOptions{
			LabelSelector: fmt.Sprintf("funcID=%d", funcID),
		},
	)

	if err != nil {
		log.Printf("[Pod List] Get Pod List Error: %s", err.Error())
		return
	}

	for _, v := range podList.Items {
		nodeName := v.Spec.NodeName
		if _, ok := list[nodeName]; ok {
			list[nodeName] = append(list[nodeName], model.PodInfo{
				Name: v.Name,
				Node: nodeName,
			})
		} else {
			list[nodeName] = []model.PodInfo{{
				Name: v.Name,
				Node: nodeName,
			}}
		}
	}

	return
}

func GetNodeList() (nodeMap map[string]*model.NodeInfo, err error) {
	nodeMap = make(map[string]*model.NodeInfo)
	var metricResp model.NodeMetricResp
	nodes, err := KubeClient.CoreV1().Nodes().List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		log.Printf("[Node List] Get Node List Error: %s", err.Error())
		return
	}

	cli := KubeClient.RESTClient().Get()
	resp, err := cli.RequestURI("/apis/metrics.k8s.io/v1beta1/nodes").DoRaw(context.Background())

	if err != nil {
		log.Printf("[Node List] Node Metric Api Error: %s", err.Error())
	}

	err = json.Unmarshal(resp, &metricResp)

	if err != nil {
		log.Printf("[Node List] Unmarshal Node Metric Api Resp Error: %s", err.Error())
		return
	}

	for _, v := range nodes.Items {
		status := "Unknown"

		for _, c := range v.Status.Conditions {
			if c.Status == "True" {
				status = string(c.Type)
			}
		}

		nodeMap[v.Name] = &model.NodeInfo{
			Name:     v.Name,
			Status:   status,
			Optional: !v.Spec.Unschedulable,
			Age:      fmt.Sprintf("%d hour", (time.Now().Unix()-v.CreationTimestamp.Unix())/3600),
			Version:  v.Status.NodeInfo.KubeletVersion,
			CpuTotal: int(v.Status.Allocatable.Cpu().MilliValue()),
			MemTotal: int(v.Status.Allocatable.Memory().Value() / (1 << 20)),
			GpuTotal: int(v.Status.Allocatable.Name(GpuQuotaKey, resource.DecimalSI).Value()),
			CpuUse:   0,
			MemUse:   0,
			GpuUse:   0,
		}
	}

	podList, err := KubeClient.CoreV1().Pods(coreV1.NamespaceDefault).List(
		context.Background(),
		metaV1.ListOptions{},
	)

	if err != nil {
		log.Printf("[Pod List] Get Pod List Error: %s", err.Error())
	} else {
		for _, pod := range podList.Items {
			for _, c := range pod.Spec.Containers {
				nodeMap[pod.Spec.NodeName].CpuUse += int(c.Resources.Requests.Cpu().MilliValue())
				nodeMap[pod.Spec.NodeName].MemUse += int(c.Resources.Requests.Memory().Value() / (1 << 20))
				nodeMap[pod.Spec.NodeName].GpuUse += int(c.Resources.Requests.Name(GpuQuotaKey, resource.DecimalSI).Value())
				nodeMap[pod.Spec.NodeName].PodNum++
			}
		}
	}

	return
}

func GetMetricsList(fid int64) (metricMap map[string]*model.PodMetricInfo, err error) {
	metricMap = make(map[string]*model.PodMetricInfo)
	var metricResp = model.PodMetricResp{}
	cli := KubeClient.RESTClient().Get()
	resp, err := cli.RequestURI("/apis/metrics.k8s.io/v1beta1/namespaces/default/pods").DoRaw(context.Background())

	if err != nil {
		log.Printf("[Pod Metrics] Metric Api Error: %s", err.Error())
	}

	err = json.Unmarshal(resp, &metricResp)

	if err != nil {
		log.Printf("[Pod Metrics] Unmarshal Api Resp Error: %s", err.Error())
	}

	podList, err := KubeClient.CoreV1().Pods(coreV1.NamespaceDefault).List(
		context.Background(),
		metaV1.ListOptions{
			LabelSelector: fmt.Sprintf("funcID=%d", fid),
		},
	)

	if err != nil {
		log.Printf("[Pod List] Get Pod List Error: %s", err.Error())
		return
	}

	for _, pod := range podList.Items {
		var gpuUsage int
		for _, c := range pod.Spec.Containers {
			gpuUseSingle := c.Resources.Requests.Name(GpuQuotaKey, resource.DecimalSI)
			if i, exist := gpuUseSingle.AsInt64(); exist {
				gpuUsage += int(i)
			}
		}
		metricMap[pod.Name] = &model.PodMetricInfo{
			NodeName:   pod.Spec.NodeName,
			PodName:    pod.Name,
			GpuUsage:   gpuUsage,
			FunctionID: fid,
			DeployName: pod.Labels[DeployNameKey],
			State:      fmt.Sprintf("%s:%s", string(pod.Status.Phase), pod.Status.Reason),
		}

	}

	for _, pod := range metricResp.Items {
		name := pod.Metadata.Name
		if _, ok := metricMap[name]; !ok {
			continue
		}
		cpu := 0
		mem := 0
		for _, c := range pod.Containers {
			cpuQuantity, err := resource.ParseQuantity(c.Usage.CPU)
			if err == nil {
				cpu += int(cpuQuantity.MilliValue())
			}
			memQuantity, err := resource.ParseQuantity(c.Usage.Memory)
			if err == nil {
				mem += int(memQuantity.Value() / 1024)
			}
		}
		metricMap[name].CpuUsage = cpu
		metricMap[name].MemUsage = mem
	}

	return
}

func GetDeploymentList() (depMap map[string]model.DeploymentInfo, err error) {
	depMap = make(map[string]model.DeploymentInfo)
	res, err := KubeClient.AppsV1().Deployments(coreV1.NamespaceDefault).List(context.Background(), metaV1.ListOptions{})

	if err != nil {
		log.Printf("[Deployment List] Get Deployment List Error: %s", err.Error())
		return
	}

	for _, v := range res.Items {
		depMap[v.Labels[FuncIDLabelKey]] = model.DeploymentInfo{
			Name:      v.Name,
			FuncIDStr: v.Labels[FuncIDLabelKey],
			Replicas:  int(*v.Spec.Replicas),
			Status:    fmt.Sprintf("%d/%d Ready", v.Status.ReadyReplicas, v.Status.Replicas),
		}
	}
	resCron, err := KubeClient.BatchV1().CronJobs(coreV1.NamespaceDefault).List(context.Background(), metaV1.ListOptions{})

	if err != nil {
		log.Printf("[Deployment List] Get Deployment List Error: %s", err.Error())
		return
	}

	for _, v := range resCron.Items {
		depMap[v.Labels[FuncIDLabelKey]] = model.DeploymentInfo{
			Name:      v.Name,
			FuncIDStr: v.Labels[FuncIDLabelKey],
			Replicas:  int(*v.Spec.JobTemplate.Spec.Parallelism),
			Status:    fmt.Sprintf("%d/%d Active Now", len(v.Status.Active), v.Spec.JobTemplate.Spec.Completions),
		}
	}
	return
}
