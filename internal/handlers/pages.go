package handlers

import (
	"html/template"
	"net/http"

	"github.com/caggri/idp4lm/internal/k8s"
)

// PageData represents the template variables needed for the HTML pages
type PageData struct {
	Namespace string
	Resources *k8s.ManagedResources

	// For configure.html checkbox states
	HasApp        bool
	HasAppSvc     bool
	HasRedis      bool
	HasRedisSvc   bool
	HasMysql      bool
	HasMysqlSvc   bool
	MysqlReplicas int32
}

func renderPage(w http.ResponseWriter, page string, data interface{}) {
	t, err := template.ParseFiles("web/templates/base.html", "web/templates/"+page)
	if err != nil {
		http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Template Render Error: "+err.Error(), http.StatusInternalServerError)
	}
}

// LandingPage renders the initial namespace selection/creation page
func LandingPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderPage(w, "landing.html", PageData{})
}

// DashboardPage renders the dashboard showing resources for a namespace
func DashboardPage(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("ns")
	if ns == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	resources, err := k8s.GetManagedResources(ns)
	if err != nil {
		// If error occurs, render the error page
		http.ServeFile(w, r, "web/templates/error.html")
		return
	}

	renderPage(w, "dashboard.html", PageData{
		Namespace: ns,
		Resources: resources,
	})
}

// ConfigurePage renders the deployment and service configuration form
func ConfigurePage(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("ns")
	if ns == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	resources, err := k8s.GetManagedResources(ns)
	if err != nil {
		http.ServeFile(w, r, "web/templates/error.html")
		return
	}

	data := PageData{
		Namespace: ns,
		Resources: resources,
	}

	// Calculate state for configure.html
	for _, dep := range resources.Deployments {
		switch dep.Name {
		case "app-deployment":
			data.HasApp = true
		case "redis-deployment":
			data.HasRedis = true
		case "mysql-deployment":
			data.HasMysql = true
			data.MysqlReplicas = dep.Replicas
		}
	}

	for _, svc := range resources.Services {
		switch svc {
		case "app-service":
			data.HasAppSvc = true
		case "redis-service":
			data.HasRedisSvc = true
		case "mysql-service":
			data.HasMysqlSvc = true
		}
	}

	renderPage(w, "configure.html", data)
}
