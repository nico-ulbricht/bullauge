package client

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(namespace string) (*v1.PodList, error) {
	k8s := getClient()

	var pods *v1.PodList
	pods, err := k8s.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pods, nil
}
