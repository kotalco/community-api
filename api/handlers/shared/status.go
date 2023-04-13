package shared

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/logger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	k8sClient       = k8s.NewClientService()
	k8sClientset, _ = k8s.NewClientset()
)

// Status returns a websocket that emits logs from pod
// Possible values are: NotFound, Pending, PodInitializing, ContainerCreating, Running, Error, Terminating
func Status(c *websocket.Conn) {
	defer c.Close()

	if os.Getenv("MOCK") == "true" {
		statuses := []string{
			"NotFound",
			"Pending",
			"PodInitializing",
			"ContainerCreating",
			"Running",
			"Error",
			"Terminating",
		}

		for {
			i := rand.Intn(len(statuses))
			c.WriteMessage(websocket.TextMessage, []byte(statuses[i]))
			time.Sleep(time.Second)
		}
	}

	ns := c.Locals("namespace").(string)
	name := c.Params("name")
	selector := fmt.Sprintf("app.kubernetes.io/managed-by=kotal-operator,app.kubernetes.io/instance=%s", name)

	watch, err := k8sClientset.CoreV1().Pods(ns).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		go logger.Info("STATUS_STREAM", err.Error())
		return
	}

	for event := range watch.ResultChan() {

		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			return
		}

		phase := string(pod.Status.Phase)

		if pod.DeletionTimestamp != nil {
			phase = "Terminating"

			// if pod is being terminated, check owner sts is found or not
			go func() {
				time.Sleep(3 * time.Second)
				_, err := k8sClientset.AppsV1().StatefulSets(ns).Get(context.Background(), name, metav1.GetOptions{})
				if err != nil && apierrors.IsNotFound(err) {
					watch.Stop()
				}
			}()
		}

		if len(pod.Status.ContainerStatuses) != 0 {
			if pod.Status.ContainerStatuses[0].State.Waiting != nil {
				phase = pod.Status.ContainerStatuses[0].State.Waiting.Reason
			}
		}

		if err := c.WriteMessage(websocket.TextMessage, []byte(phase)); err != nil {
			return
		}
	}

}
