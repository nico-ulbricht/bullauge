package client

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var k8s *kubernetes.Clientset

func GetClient() *kubernetes.Clientset {
	if k8s != nil {
		return k8s
	}

	var c *kubernetes.Clientset
	switch os.Getenv("K8S_CONNECTION_METHOD") {
	case "local":
		c = getLocalClient()
		break
	case "remote":
		c = getRemoteClient()
		break
	default:
		log.Fatalf("K8S_CONNECTION_METHOD neither local nor remote.")
	}

	k8s = c
	return k8s
}

func getLocalClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to created Cluster Configuration. Error: %v", err)
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes Client. Error: %v", err)
	}

	return c
}

func getRemoteClient() *kubernetes.Clientset {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Failed to created Cluster Configuration. Error: %v", err)
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes Client. Error: %v", err)
	}

	return c
}
