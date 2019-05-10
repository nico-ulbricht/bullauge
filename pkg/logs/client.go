package logs

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/nico-ulbricht/bullauge/pkg/client"
	v1 "k8s.io/api/core/v1"
)

func GetLogs(name, namespace string, limit int) ([]string, error) {
	k8s := client.GetClient()

	limit64 := int64(limit)
	request := k8s.CoreV1().Pods(namespace).GetLogs(name, &v1.PodLogOptions{
		TailLines: &limit64,
	})

	readCloser, err := request.Stream()
	if err != nil {
		return []string{}, err
	}

	defer readCloser.Close()
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	_, err = io.Copy(writer, readCloser)

	logs := filterEmpty(strings.Split(b.String(), "\n"))
	return logs, nil
}

func filterEmpty(logs []string) []string {
	filteredLogs := make([]string, 0)
	for _, log := range logs {
		if log == "" {
			continue
		}

		filteredLogs = append(filteredLogs, log)
	}

	return filteredLogs
}
