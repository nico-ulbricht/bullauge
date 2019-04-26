package pod

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/nico-ulbricht/bullauge/pkg/logs"
	v1 "k8s.io/api/core/v1"
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
			Description: "Image that is being used by the POD",
			Type:        graphql.String,
		},
		"logs": &graphql.Field{
			Description: "Logs of the POD",
			Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					DefaultValue: 10,
					Description:  "Maximum number of logs to return",
					Type:         graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pod := p.Source.(*pod)
				if pod.Status == "containercreating" || pod.Status == "pending" {
					return []string{}, nil
				}

				lineLimit := p.Args["limit"].(int)
				logs := logs.GetLogs(pod.Name, pod.Namespace, lineLimit)
				return logs, nil
			},
		},
		"name": &graphql.Field{
			Description: "Name of the POD",
			Type:        graphql.String,
		},
		"namespace": &graphql.Field{
			Description: "Namespace the POD is running in",
			Type:        graphql.String,
		},
		"status": &graphql.Field{
			Description: "Current status of the POD",
			Type:        graphql.String,
		},
	},
})

var Query = graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(podType))),
	Description: "List of PODs on the Cluster",
	Args: graphql.FieldConfigArgument{
		"app": &graphql.ArgumentConfig{
			Description: "Prefix of the POD name",
			Type:        graphql.String,
		},
		"namespace": &graphql.ArgumentConfig{
			Description: "Namespace of the POD",
			Type:        graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		podList, _ := GetPods(p.Args["namespace"].(string))
		pods := convertPodList(podList)
		if p.Args["app"] != nil {
			pods = filterPods(pods, p.Args["app"].(string))
		}

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

//TODO: pass filter to k8s client instead of inmem filtering
func filterPods(pods []*pod, filterPrefix string) []*pod {
	filteredPods := []*pod{}
	for _, pod := range pods {
		if strings.HasPrefix(pod.Name, filterPrefix) == false {
			continue
		}

		filteredPods = append(filteredPods, pod)
	}

	return filteredPods
}
