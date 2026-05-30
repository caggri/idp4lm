package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateQuota creates a ResourceQuota for the namespace
func CreateQuota(namespace string) error {
	quota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "resource-quota",
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceRequestsCPU:            resource.MustParse("2"), // TODO:  Parameterize all
				corev1.ResourceRequestsMemory:         resource.MustParse("4Gi"),
				corev1.ResourceLimitsCPU:              resource.MustParse("6"),
				corev1.ResourceLimitsMemory:           resource.MustParse("8Gi"),
				corev1.ResourcePods:                   resource.MustParse("10"),
				corev1.ResourceServices:               resource.MustParse("5"),
				corev1.ResourceConfigMaps:             resource.MustParse("10"),
				corev1.ResourceSecrets:                resource.MustParse("10"),
				corev1.ResourcePersistentVolumeClaims: resource.MustParse("5"),
				corev1.ResourceRequestsStorage:        resource.MustParse("100Gi"),
			},
		},
	}

	_, err := Client.CoreV1().ResourceQuotas(namespace).Create(context.Background(), quota, metav1.CreateOptions{})
	return err
}
