package shared

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/pkg/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

var (
	k8sClient = k8s.NewClientService()
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

	sts := &appsv1.StatefulSet{}
	stsKey := types.NamespacedName{
		Namespace: c.Locals("namespace").(string),
		Name:      c.Params("name"),
	}

	pod := &corev1.Pod{}
	key := types.NamespacedName{
		Namespace: c.Locals("namespace").(string),
		Name:      fmt.Sprintf("%s-0", c.Params("name")),
	}

	for {

		err := k8sClient.Get(context.Background(), key, pod)
		if err != nil {

			// is the pod error due to sts has been deleted ?
			stsErr := k8sClient.Get(context.Background(), stsKey, sts)
			if apierrors.IsNotFound(stsErr) {
				return
			}

			if apierrors.IsNotFound(err) {
				err = errors.New("NotFound")
			}

			c.WriteMessage(websocket.TextMessage, []byte(err.Error()))

			time.Sleep(3 * time.Second)
			continue
		}

		phase := string(pod.Status.Phase)
		if pod.DeletionTimestamp != nil {
			phase = "Terminating"
		}
		if len(pod.Status.ContainerStatuses) != 0 {
			if pod.Status.ContainerStatuses[0].State.Waiting != nil {
				phase = pod.Status.ContainerStatuses[0].State.Waiting.Reason
			}
		}

		c.WriteMessage(websocket.TextMessage, []byte(phase))

		time.Sleep(time.Second)
	}

}
