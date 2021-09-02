package k8s

import (
	"go/build"
	"log"
	"os"
	"path/filepath"
	"sync"

	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var k8sClient client.Client
var once sync.Once

// Client create k8s client once
func Client() client.Client {
	once.Do(func() {
		k8sClient = NewClient()
	})

	return k8sClient
}

// NewClient creates new k8s client
func NewClient() client.Client {

	var config *rest.Config
	var err error

	// if we're in k8s cluster, create in cluster config using service account
	// otherwise, create out of cluster config using kubeconfig at $HOME/.kube/config
	// TODO: don't panic
	if os.Getenv("MOCK") == "true" {
		log.Println("creating k8s client using test environment ...")
		testEnv := envtest.Environment{
			CRDDirectoryPaths: []string{
				filepath.Join(build.Default.GOPATH, "pkg", "mod", "github.com", "kotalco", "kotal@v0.0.0-20210817190935-979d6e70b8e5", "config", "crd", "bases"),
			},
			ErrorIfCRDPathMissing: true,
		}
		config, err = testEnv.Start()
		if err != nil {
			panic(err.Error())
		}
	} else if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		log.Println("creating k8s client using in-cluster config ...")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		log.Println("creating k8s client using out-of-cluster config ...")
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	ethereumv1alpha1.AddToScheme(scheme)
	ethereum2v1alpha1.AddToScheme(scheme)
	ipfsv1alpha1.AddToScheme(scheme)

	opts := client.Options{Scheme: scheme}

	c, err := client.New(config, opts)
	if err != nil {
		// TODO: don't panic
		panic(err.Error())
	}

	return c
}
