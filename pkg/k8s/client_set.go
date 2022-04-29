package k8s

import (
	"github.com/kotalco/api/pkg/configs"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var KubernetesClientset *kubernetes.Clientset
var clientsetOnce sync.Once

// Clientset create k8s client once
func Clientset() *kubernetes.Clientset {
	var err error
	clientsetOnce.Do(func() {
		KubernetesClientset, err = NewClientset()
		if err != nil {
			// TODO: Don't panic
			panic(err)
		}
	})
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
