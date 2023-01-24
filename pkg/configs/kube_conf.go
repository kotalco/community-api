package configs

import (
	"go/build"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// KubeConfig returns REST config based on the environment
func KubeConfig() (*rest.Config, error) {

	// if we're in k8s cluster, create in cluster config using service account
	// otherwise, create out of cluster config using kubeconfig at $HOME/.kube/config
	if os.Getenv("MOCK") == "true" {
		log.Println("creating k8s client using test environment ...")
		testEnv := envtest.Environment{
			CRDDirectoryPaths: []string{
				filepath.Join(build.Default.GOPATH, "pkg", "mod", "github.com", "kotalco", "kotal@v0.0.0-20230124184344-65d52739f9c5", "config", "crd", "bases"),
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
