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

	opts := corev1.PodLogOptions{
		Follow:    true,
		TailLines: &lines,
	}

	ns := c.Locals("namespace").(string)
	name := fmt.Sprintf("%s-0", c.Params("name"))
	logs := k8s.Clientset().CoreV1().Pods(ns).GetLogs(name, &opts)

	stream, err := logs.Stream(context.TODO())
	if stream != nil {
		defer stream.Close()
	}
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	for {
		buf := make([]byte, 1024)

		numBytes, err := stream.Read(buf)
		if err != nil {
			return
		}

		if numBytes == 0 {
			time.Sleep(time.Second)
			continue
		}

		if err := c.WriteMessage(websocket.TextMessage, buf[:numBytes]); err != nil {
			return
		}
	}
}
