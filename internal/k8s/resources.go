package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodInfo struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
}

type DeploymentInfo struct {
	Name     string `json:"name"`
	Ready    bool   `json:"ready"`
	Replicas int32  `json:"replicas"`
}

type ManagedResources struct {
	Pods        []PodInfo        `json:"pods"`
	Services    []string         `json:"services"`
	Deployments []DeploymentInfo `json:"deployments"`
}

// GetManagedResources lists all k8s-idp managed resources in a namespace
func GetManagedResources(namespace string) (*ManagedResources, error) {
	ctx := context.Background()
	labelSelector := "managed-by=k8s-idp"
	listOpts := metav1.ListOptions{LabelSelector: labelSelector}

	// 1. Get Pods
	pods, err := Client.CoreV1().Pods(namespace).List(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	podInfoList := make([]PodInfo, 0)
	for _, pod := range pods.Items {
		ready := true
		if len(pod.Status.ContainerStatuses) == 0 {
			ready = false
		} else {
			for _, status := range pod.Status.ContainerStatuses {
				if !status.Ready {
					ready = false
					break
				}
			}
		}
		podInfoList = append(podInfoList, PodInfo{Name: pod.Name, Ready: ready})
	}

	// 2. Get Services
	services, err := Client.CoreV1().Services(namespace).List(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	serviceNames := make([]string, 0)
	for _, svc := range services.Items {
		serviceNames = append(serviceNames, svc.Name)
	}

	// 3. Get Deployments
	deployments, err := Client.AppsV1().Deployments(namespace).List(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	deploymentInfoList := make([]DeploymentInfo, 0)
	for _, dep := range deployments.Items {
		ready := false
		var replicas int32 = 1
		if dep.Spec.Replicas != nil {
			replicas = *dep.Spec.Replicas
			if dep.Status.ReadyReplicas == replicas {
				ready = true
			}
		}
		deploymentInfoList = append(deploymentInfoList, DeploymentInfo{Name: dep.Name, Ready: ready, Replicas: replicas})
	}

	return &ManagedResources{
		Pods:        podInfoList,
		Services:    serviceNames,
		Deployments: deploymentInfoList,
	}, nil
}
