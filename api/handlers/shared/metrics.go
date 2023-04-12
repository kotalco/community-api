package shared

import (
	"context"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type metricsResponseDto struct {
	Cpu    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

var metricsClientset = k8s.MetricsClientset()

// Metrics returns a websocket that emits cpu and memory usage
func Metrics(c *websocket.Conn) {
	defer c.Close()

	name := c.Params("name")
	ns := c.Locals("namespace").(string)

	// Create a new Metrics Query object to get the CPU and Memory usage of the pod
	selector := fmt.Sprintf("app.kubernetes.io/managed-by=kotal-operator,app.kubernetes.io/instance=%s", name)

	// Watch the Metrics of the pod using the Metrics API
	watcher, err := metricsClientset.MetricsV1beta1().PodMetricses(ns).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		go logger.Info("METRICS_STREAM", err.Error())
		return
	}

	for event := range watcher.ResultChan() {
		if event.Type == "ERROR" {
			err := event.Object.(*v1beta1.PodMetrics)
			go logger.Info("METRICS_STREAM", fmt.Sprintf("Error watching Metrics: %v", err))
			return
		}

		metrics := event.Object.(*v1beta1.PodMetrics)
		response := new(metricsResponseDto)

		response.Cpu = metrics.Containers[0].Usage.Cpu().ScaledValue(resource.Milli)
		response.Memory = metrics.Containers[0].Usage.Memory().ScaledValue(resource.Mega)

		if err := c.WriteJSON(response); err != nil {
			return
		}
	}
}
