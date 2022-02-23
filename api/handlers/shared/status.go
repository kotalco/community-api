package shared

import (
	"context"
	"errors"
	"fmt"
	"github.com/kotalco/api/pkg/k8s"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
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
		Namespace: "default",
		Name:      c.Params("name"),
	}

	pod := &corev1.Pod{}
	key := types.NamespacedName{
		Namespace: "default",
		Name:      fmt.Sprintf("%s-0", c.Params("name")),
	}

	for {
		err := k8s.Client().Get(context.Background(), stsKey, sts)
		stsNotFound := apierrors.IsNotFound(err)

		err = k8s.Client().Get(context.Background(), key, pod)
		if err != nil {
			if stsNotFound {
				return
			} else {
				if apierrors.IsNotFound(err) {
					err = errors.New("NotFound")
				}
				c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				continue
			}
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

		time.Sleep(1 * time.Second)
	}

}