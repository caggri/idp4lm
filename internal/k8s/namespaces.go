package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceResult represents the result of namespace creation
type NamespaceResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Exists  bool   `json:"exists"`
}

// CheckAndCreateNamespace creates the namespace and its base security policies
func CheckAndCreateNamespace(namespace string) (*NamespaceResult, error) {
	exists, err := CheckNamespaceExists(namespace)
	if err != nil {
		return nil, fmt.Errorf("namespace check error: %v", err)
	}

	if exists {
		return &NamespaceResult{
			Status:  "success",
			Message: fmt.Sprintf("Namespace %s already exists", namespace),
			Exists:  true,
		}, nil
	}

	// Create namespace
	if err := CreateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("failed to create namespace: %v", err)
	}

	// Create NetworkPolicy, RBAC, and Quota. If any fail, return an error.
	if err := CreateNetworkPolicy(namespace); err != nil {
		return nil, fmt.Errorf("failed to create network policy: %v", err)
	}
	if err := CreateRBAC(namespace); err != nil {
		return nil, fmt.Errorf("failed to create rbac: %v", err)
	}
	if err := CreateQuota(namespace); err != nil {
		return nil, fmt.Errorf("failed to create quota: %v", err)
	}

	return &NamespaceResult{
		Status:  "success",
		Message: fmt.Sprintf("Namespace %s created successfully", namespace),
		Exists:  false,
	}, nil
}

// CheckNamespaceExists checks if a given namespace already exists
func CheckNamespaceExists(namespace string) (bool, error) {
	_, err := Client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateNamespace creates a new namespace with the managed-by label
func CreateNamespace(namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"managed-by": "k8s-idp",
			},
		},
	}
	_, err := Client.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	return err
}
