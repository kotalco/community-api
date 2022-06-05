package k8s

import (
	"github.com/kotalco/api/pkg/configs"
	"github.com/kotalco/api/pkg/logger"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var clientSetLock = &sync.Mutex{}
var KubernetesClientset *kubernetes.Clientset

// Clientset create k8s client once
func Clientset() *kubernetes.Clientset {
	var err error
	clientSetLock.Lock()
	defer clientSetLock.Unlock()
	if KubernetesClientset == nil {
		KubernetesClientset, err = NewClientset()
		if err != nil {
			logger.Panic("K8S_CLIENT_SET", err)
		}
	}

	return KubernetesClientset
}

// NewClientset returns client-go clientset
func NewClientset() (*kubernetes.Clientset, error) {
	config, err := configs.KubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
