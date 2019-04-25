package pod

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/graphql-go/graphql"
	"github.com/nico-ulbricht/bullauge/pkg/client"
)

type pod struct {
	Image     string `json:"image"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
}

var podType = graphql.NewObject(graphql.ObjectConfig{
	Name: "POD",
	Fields: graphql.Fields{
		"image": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"namespace": &graphql.Field{
			Type: graphql.String,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var Query = graphql.Field{
	Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(podType))),
	Args: graphql.FieldConfigArgument{
		"app": &graphql.ArgumentConfig{
			Description: "Name of the application",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"namespace": &graphql.ArgumentConfig{
			Description: "Namespace of the application",
			Type:        graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		podList, _ := client.GetPods()
		pods := convertPodList(podList)
		pods = filterPods(pods, &podFilterConfig{
			App:       p.Args["app"].(string),
			Namespace: p.Args["namespace"].(string),
		})
		return pods, nil
	},
}

func convertPodList(podList *v1.PodList) []*pod {
	pods := make([]*pod, len(podList.Items))
	for idx, podItem := range podList.Items {
		pods[idx] = &pod{
			Image:     podItem.Spec.Containers[0].Image,
			Name:      podItem.GetName(),
			Namespace: podItem.GetNamespace(),
			Status:    strings.ToLower(fmt.Sprintf("%s", podItem.Status.Phase)),
		}
	}

	return pods
}

type podFilterConfig struct {
	App       string
	Namespace string
}

//TODO: pass filter to k8s client instead of inmem filtering
func filterPods(pods []*pod, filterConfig *podFilterConfig) []*pod {
	filteredPods := []*pod{}
	for _, pod := range pods {
		if strings.HasPrefix(pod.Name, filterConfig.App) == false {
			continue
		}

		if filterConfig.Namespace != "" && strings.ToLower(pod.Namespace) == strings.ToLower(filterConfig.Namespace) == false {
			continue
		}

		filteredPods = append(filteredPods, pod)
	}

	return filteredPods
}
