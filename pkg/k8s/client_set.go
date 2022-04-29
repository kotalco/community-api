package k8s

import (
	"context"
	"github.com/kotalco/api/pkg/configs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type k8ClientSetService struct{}
type k8ClientSetServiceInterface interface {
}

var (
	K8ClientSetService k8ClientSetServiceInterface
)

func init() { K8ClientSetService = &k8ClientSetService{} }

func (k8ClientSets *k8ClientSetService) CreateWorkspace(name string) error {
	nsName := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := Clientset().CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (k8ClientSets *k8ClientSetService) GetWorkspace(name string) (*corev1.Namespace, error) {
	workspace, err := Clientset().CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (k8ClientSets *k8ClientSetService) UpdateWorkspace(name string) error {
	nsName := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := Clientset().CoreV1().Namespaces().Update(context.Background(), nsName, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
