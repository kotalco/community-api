package k8s

import (
	"github.com/kotalco/community-api/pkg/configs"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"sync"
)

var metricsClientset *metrics.Clientset
var metricsClientsetOnce sync.Once

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

// NewMetricsClientset returns metrics client
func NewMetricsClientset() (*metrics.Clientset, error) {
	config, err := configs.KubeConfig()
	if err != nil {
		return nil, err
	}

	return metrics.NewForConfig(config)
}
