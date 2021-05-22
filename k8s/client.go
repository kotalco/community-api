package k8s

import (
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Client() client.Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		// TODO: don't panic
		panic(err.Error())
	}

	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	ethereumv1alpha1.AddToScheme(scheme)

	opts := client.Options{Scheme: scheme}

	c, err := client.New(config, opts)
	if err != nil {
		// TODO: don't panic
		panic(err.Error())
	}

	return c
}
