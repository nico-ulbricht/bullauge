package client

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/nico-ulbricht/bullauge/pkg/client"
	"k8s.io/api/core/v1"
)

func GetLogs(name, namespace string) string {
	k8s := client.GetClient()

	request := k8s.CoreV1().Pods(namespace).GetLogs(name, &v1.PodLogOptions{})
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
