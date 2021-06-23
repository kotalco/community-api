package k8s

import (
	"os"
	"path/filepath"

	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Client() client.Client {

	var config *rest.Config
	var err error

	// if we're in k8s cluster, create in cluster config using service account
	// otherwise, create out of cluster config using kubeconfig at $HOME/.kube/config
	// TODO: don't panic
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
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
