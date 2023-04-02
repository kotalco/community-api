package shared

import (
	"context"
	"fmt"
	"github.com/gofiber/websocket/v2"
	restError "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	"github.com/kotalco/community-api/pkg/shared"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var metricsClientset = k8s.MetricsClientset()

type metricsResponseDto struct {
	Cpu    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

// Metrics returns a websocket that emits cpu and memory usage
func Metrics(c *websocket.Conn) {
	defer c.Close()

	namespace := c.Locals("namespace").(string)
	name := c.Params("name")
	podKey := types.NamespacedName{
		Namespace: namespace,
		Name:      fmt.Sprintf("%s-0", name),
	}
	stsKey := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}

	podStatus := make(chan error)
	metricsResponse := make(chan *metricsResponseDto)

	go fetchPodStatus(podKey, stsKey, podStatus)
	go fetchMetrics(podKey, metricsResponse)

	for {
		select {
		case err := <-podStatus:
			if err != nil {
				go logger.Info("METRICS_POD_NOTFOUND", err.Error())
				c.WriteJSON(shared.NewResponse(restError.NewNotFoundError(err.Error())))
				return
			}
		case response := <-metricsResponse:
			if response == nil {
				continue
			}

			//snap out of the infinte loop if the connection get closed
			if err := c.WriteJSON(response); err != nil {
				go logger.Info("METRICS_JSON_ERR", err.Error())
				return
			}
		}
	}
}

func fetchPodStatus(podKey types.NamespacedName, stsKey types.NamespacedName, status chan error) {
	stsCount := 0
	for {
		stsCount++
		fmt.Println(stsCount)
		pod := new(corev1.Pod)
		err := k8sClient.Get(context.Background(), podKey, pod)
		if err != nil {
			if apierrors.IsNotFound(err) {
				stsErr := k8sClient.Get(context.Background(), stsKey, &appsv1.StatefulSet{})
				if apierrors.IsNotFound(stsErr) {
					status <- fmt.Errorf("statefulset not found: %s", stsKey)
					return
				}
			}
			time.Sleep(3 * time.Second)
			continue
		}
		status <- nil
		return
	}
}

func fetchMetrics(podKey types.NamespacedName, response chan *metricsResponseDto) {
	cpuIdx := 0
	memoryIdx := 0

	metricsCount := 0
	for range time.Tick(1 * time.Second) {
		metricsCount++
		fmt.Println(metricsCount)
		metrics, err := metricsClientset.MetricsV1beta1().PodMetricses(podKey.Namespace).Get(context.Background(), podKey.Name, metav1.GetOptions{})
		if err != nil {
			go logger.Info("METRICS_API_ERR", err.Error())
			response <- nil
			continue
		}

		if len(metrics.Containers) > 0 {
			cpu := metrics.Containers[cpuIdx].Usage.Cpu().ScaledValue(resource.Milli)
			memory := metrics.Containers[memoryIdx].Usage.Memory().ScaledValue(resource.Mega)

			response <- &metricsResponseDto{
				Cpu:    cpu,
				Memory: memory,
			}
		}
	}
}
