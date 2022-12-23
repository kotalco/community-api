package shared

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
)

// Logger returns a websocket that emits logs from pod
func Logger(c *websocket.Conn) {
	defer c.Close()

	if os.Getenv("MOCK") == "true" {
		var i int
		for {
			i++
			if i == 10 {
				return
			}
			msg := fmt.Sprintf("%s \n", time.Now().Local())
			c.WriteMessage(websocket.TextMessage, []byte(msg))
			time.Sleep(time.Second)
		}
	}

	lines := int64(100)

	podLogOptions := corev1.PodLogOptions{
		Follow:    true,
		TailLines: &lines,
	}

	podLogRequest := k8s.Clientset().CoreV1().Pods(c.Locals("namespace").(string)).GetLogs(fmt.Sprintf("%s-0", c.Params("name")), &podLogOptions)

	stream, err := podLogRequest.Stream(context.TODO())
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer stream.Close()

	for {
		buf := make([]byte, 1024)
		numBytes, err := stream.Read(buf)
		if err != nil {
			break
		}

		if numBytes == 0 {
			continue
		}

		c.WriteMessage(websocket.TextMessage, buf[:numBytes])
	}
}
