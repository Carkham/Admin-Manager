package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var KubeClient *kubernetes.Clientset

func KubeClientInit(projectPath string) {

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
