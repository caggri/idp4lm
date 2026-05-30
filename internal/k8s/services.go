package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateService creates a Kubernetes Service for a specific deployment type
func CreateService(namespace, svcType string) error {
	var port int32

	switch svcType {
	case "app":
		port = 80
	case "redis":
		port = 6379
	case "mysql":
		port = 3306
	default:
		return fmt.Errorf("unsupported service type: %s", svcType)
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-service", svcType),
			Namespace: namespace,
			Labels: map[string]string{
				"app":        svcType,
				"managed-by": "k8s-idp",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": svcType,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       port,
					TargetPort: intstr.FromInt32(port),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	_, err := Client.CoreV1().Services(namespace).Create(context.Background(), svc, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}
