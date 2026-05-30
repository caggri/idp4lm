package k8s

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// CreateDeployment creates a Kubernetes Deployment. If it exists, updates the replica count.
// TODO: Add service after first Deployment
func CreateDeployment(namespace, depType string, replicas int32) error {
	var containerName, image string
	var port int32
	var envs []corev1.EnvVar

	switch depType {
	case "app":
		containerName = "nginx"
		image = "nginx:latest"
		port = 80
	case "redis":
		containerName = "redis"
		image = "redis:latest"
		port = 6379
	case "mysql":
		containerName = "mysql"
		image = "mysql:latest"
		port = 3306
		// Read password from env
		pwd := os.Getenv("MYSQL_ROOT_PASSWORD")
		if pwd == "" {
			pwd = "password123" // Testing
		}
		envs = []corev1.EnvVar{
			{Name: "MYSQL_ROOT_PASSWORD", Value: pwd},
		}
	default:
		return fmt.Errorf("unsupported deployment type: %s", depType)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-deployment", depType),
			Namespace: namespace,
			Labels: map[string]string{
				"app":        depType,
				"managed-by": "k8s-idp",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": depType},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        depType,
						"managed-by": "k8s-idp",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: image,
							Ports: []corev1.ContainerPort{{ContainerPort: port}},
							Env:   envs,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("250m"),  // TODO: parameterize them
									corev1.ResourceMemory: resource.MustParse("128Mi"), // TODO: parameterize them
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),  // TODO: parameterize them
									corev1.ResourceMemory: resource.MustParse("768Mi"), // TODO: parameterize them
								},
							},
						},
					},
				},
			},
		},
	}

	appsClient := Client.AppsV1()
	ctx := context.Background()

	_, err := appsClient.Deployments(namespace).Create(ctx, dep, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			// BUG FIX: Instead of throwing 409 Conflict, we patch/update the existing deployment's replicas
			existing, getErr := appsClient.Deployments(namespace).Get(ctx, dep.Name, metav1.GetOptions{})
			if getErr != nil {
				return fmt.Errorf("failed to get existing deployment: %v", getErr)
			}
			existing.Spec.Replicas = ptr.To(replicas)
			_, updateErr := appsClient.Deployments(namespace).Update(ctx, existing, metav1.UpdateOptions{})
			if updateErr != nil {
				return fmt.Errorf("failed to update deployment replicas: %v", updateErr)
			}
			return nil
		}
		return err
	}

	return nil
}
