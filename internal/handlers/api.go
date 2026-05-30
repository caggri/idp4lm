package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/caggri/idp4lm/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper to read JSON request
func readJSON(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func ListNamespacesAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	namespaces, err := k8s.Client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{
		LabelSelector: "managed-by=k8s-idp",
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	type NamespaceInfo struct {
		Name    string `json:"name"`
		Created string `json:"created"`
	}

	var names []NamespaceInfo
	for _, ns := range namespaces.Items {
		names = append(names, NamespaceInfo{
			Name:    ns.Name,
			Created: ns.CreationTimestamp.String(),
		})
	}
	// Return empty array instead of null
	if names == nil {
		names = []NamespaceInfo{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"namespaces": names})
}

type CreateNamespaceReq struct {
	Namespace string `json:"namespace"`
}

func CreateNamespaceAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateNamespaceReq
	if err := readJSON(r, &req); err != nil || req.Namespace == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	res, err := k8s.CheckAndCreateNamespace(req.Namespace)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

type DeploymentReq struct {
	Type     string `json:"type"`
	Replicas int32  `json:"replicas"`
}

type CreateDeploymentsReq struct {
	Namespace   string          `json:"namespace"`
	Deployments []DeploymentReq `json:"deployments"`
}

func CreateDeploymentsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateDeploymentsReq
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	for _, dep := range req.Deployments {
		if err := k8s.CreateDeployment(req.Namespace, dep.Type, dep.Replicas); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Failed to create %s deployment: %v", dep.Type, err),
			})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Deployments created/updated successfully"})
}

type ServiceReq struct {
	Type string `json:"type"`
}

type CreateServicesReq struct {
	Namespace string       `json:"namespace"`
	Services  []ServiceReq `json:"services"`
}

func CreateServicesAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateServicesReq
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	// BUG FIX: Loop through all services instead of returning on the first one
	var errors []string
	for _, svc := range req.Services {
		if err := k8s.CreateService(req.Namespace, svc.Type); err != nil {
			errors = append(errors, fmt.Sprintf("failed to create %s service: %v", svc.Type, err))
		}
	}

	if len(errors) > 0 {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"error":   "Failed to create some services",
			"details": errors,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Services created successfully"})
}

func DeleteResourceAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/resources/"), "/")
	if len(pathParts) != 2 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid resource path"})
		return
	}
	resType := pathParts[0]
	resName := pathParts[1]
	namespace := r.URL.Query().Get("ns")

	if namespace == "" && resType != "namespace" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Namespace is required"})
		return
	}

	ctx := context.Background()
	var err error
	switch resType {
	case "pod":
		err = k8s.Client.CoreV1().Pods(namespace).Delete(ctx, resName, metav1.DeleteOptions{})
	case "service":
		err = k8s.Client.CoreV1().Services(namespace).Delete(ctx, resName, metav1.DeleteOptions{})
	case "deployment":
		err = k8s.Client.AppsV1().Deployments(namespace).Delete(ctx, resName, metav1.DeleteOptions{})
	case "namespace":
		err = k8s.Client.CoreV1().Namespaces().Delete(ctx, resName, metav1.DeleteOptions{})
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unsupported resource type"})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success", "message": fmt.Sprintf("%s deleted", resName)})
}
