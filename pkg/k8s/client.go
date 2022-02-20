package k8s

import (
	"go/build"
	"log"
	"os"
	"path/filepath"
	"sync"

	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
	polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var ControllerRuntimeClient client.Client
var KubernetesClientset *kubernetes.Clientset
var metricsClientset *metrics.Clientset

var clientOnce sync.Once
var clientsetOnce sync.Once
var metricsClientsetOnce sync.Once

// Client create k8s client once
func Client() client.Client {
	var err error
	clientOnce.Do(func() {
		ControllerRuntimeClient, err = NewRuntimeClient()
		if err != nil {
			// TODO: Don't panic
			panic(err)
		}
	})
	return ControllerRuntimeClient
}

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

// MetricsClientset create k8s metrics client once
func MetricsClientset() *metrics.Clientset {
	var err error
	metricsClientsetOnce.Do(func() {
		metricsClientset, err = NewMetricsClientset()
		if err != nil {
			// TODO: Don't panic
			panic(err)
		}
	})
	return metricsClientset
}

// Config returns REST config based on the environment
func Config() (*rest.Config, error) {

	// if we're in k8s cluster, create in cluster config using service account
	// otherwise, create out of cluster config using kubeconfig at $HOME/.kube/config
	if os.Getenv("MOCK") == "true" {
		log.Println("creating k8s client using test environment ...")
		testEnv := envtest.Environment{
			CRDDirectoryPaths: []string{
				filepath.Join(build.Default.GOPATH, "pkg", "mod", "github.com", "kotalco", "kotal@v0.0.0-20220117145334-7bdbeb90323c", "config", "crd", "bases"),
			},
			ErrorIfCRDPathMissing: true,
		}
		return testEnv.Start()
	} else if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		log.Println("creating k8s client using in-cluster config ...")
		return rest.InClusterConfig()
	} else {
		log.Println("creating k8s client using out-of-cluster config ...")
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

}

// NewRuntimeClient creates new controller-runtime k8s client
func NewRuntimeClient() (client.Client, error) {

	config, err := Config()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	ethereumv1alpha1.AddToScheme(scheme)
	ethereum2v1alpha1.AddToScheme(scheme)
	ipfsv1alpha1.AddToScheme(scheme)
	filecoinv1alpha1.AddToScheme(scheme)
	chainlinkv1alpha1.AddToScheme(scheme)
	polkadotv1alpha1.AddToScheme(scheme)
	nearv1alpha1.AddToScheme(scheme)

	opts := client.Options{Scheme: scheme}

	return client.New(config, opts)
}

// NewClientset returns client-go clientset
func NewClientset() (*kubernetes.Clientset, error) {
	config, err := Config()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// NewMetricsClientset returns metrics client
func NewMetricsClientset() (*metrics.Clientset, error) {
	config, err := Config()
	if err != nil {
		return nil, err
	}

	return metrics.NewForConfig(config)
}
