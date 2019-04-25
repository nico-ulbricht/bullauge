package client

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"k8s.io/api/core/v1"
)

func GetLogs(name, namespace string) string {
	k8s := getClient()

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
