package logs

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/nico-ulbricht/bullauge/pkg/client"
	"k8s.io/api/core/v1"
)

func GetLogs(name, namespace string, limit int) string {
	k8s := client.GetClient()

	limit64 := int64(limit)
	request := k8s.CoreV1().Pods(namespace).GetLogs(name, &v1.PodLogOptions{
		TailLines: &limit64,
	})

	readCloser, err := request.Stream()
	if err != nil {
		log.Fatalf("Error requesting logs. Error: %v", err)
	}

	defer readCloser.Close()
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	_, err = io.Copy(writer, readCloser)

	return b.String()
}
