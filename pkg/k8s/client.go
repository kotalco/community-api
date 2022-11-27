package k8s

import (
	"context"
	"sync"

	"github.com/kotalco/community-api/pkg/configs"
	"github.com/kotalco/community-api/pkg/logger"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
	polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	clientLock              = &sync.Mutex{}
	controllerRuntimeClient client.Client
	RunTimeScheme           = runtime.NewScheme()
)

func newClient() client.Client {
	var err error
	clientLock.Lock()
	defer clientLock.Unlock()
	if controllerRuntimeClient == nil {
		controllerRuntimeClient, err = newRuntimeClient()
		if err != nil {
			logger.Panic("K8S_CLIENT", err)
		}
	}

	return controllerRuntimeClient
}

// newRuntimeClient creates new controller-runtime k8s client
func newRuntimeClient() (client.Client, error) {
	config, err := configs.KubeConfig()
	if err != nil {
		return nil, err
	}

	clientgoscheme.AddToScheme(RunTimeScheme)
	ethereumv1alpha1.AddToScheme(RunTimeScheme)
	ethereum2v1alpha1.AddToScheme(RunTimeScheme)
	ipfsv1alpha1.AddToScheme(RunTimeScheme)
	filecoinv1alpha1.AddToScheme(RunTimeScheme)
	chainlinkv1alpha1.AddToScheme(RunTimeScheme)
	polkadotv1alpha1.AddToScheme(RunTimeScheme)
	nearv1alpha1.AddToScheme(RunTimeScheme)

	opts := client.Options{Scheme: RunTimeScheme}

	return client.New(config, opts)
}

type k8sClientService struct{}
type ObjectKey = types.NamespacedName

type K8sClientServiceInterface interface {
	client.Reader
	client.Writer
}

func NewClientService() K8sClientServiceInterface {
	k8client := k8sClientService{}
	return k8client
}

// Get retrieves an obj for the given object key from the Kubernetes Cluster.
// obj must be a struct pointer so that obj can be updated with the response
// returned by the Server.
func (k8sClient k8sClientService) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	return newClient().Get(ctx, key, obj)
}

// List retrieves list of objects for a given namespace and list options. On a
// successful call, Items field in the list will be populated with the
// result returned from the server.
func (k8sClient k8sClientService) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return newClient().List(ctx, list, opts...)

}

// Create saves the object obj in the Kubernetes cluster.
func (k8sClient k8sClientService) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return newClient().Create(ctx, obj, opts...)
}

// Delete deletes the given obj from Kubernetes cluster.
func (k8sClient k8sClientService) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return newClient().Delete(ctx, obj, opts...)
}

// Update updates the given obj in the Kubernetes cluster. obj must be a
// struct pointer so that obj can be updated with the content returned by the Server.
func (k8sClient k8sClientService) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return newClient().Update(ctx, obj, opts...)
}

// Patch patches the given obj in the Kubernetes cluster. obj must be a
// struct pointer so that obj can be updated with the content returned by the Server.
func (k8sClient k8sClientService) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return newClient().Patch(ctx, obj, patch, opts...)
}

// DeleteAllOf deletes all objects of the given type matching the given options.
func (k8sClient k8sClientService) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return newClient().DeleteAllOf(ctx, obj, opts...)
}
