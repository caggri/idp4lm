package handlers

import (
	"net/http"
)

// SetupRoutes registers all the HTTP handlers with the standard DefaultServeMux
func SetupRoutes() {
	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// HTML Pages
	http.HandleFunc("/", LandingPage)
	http.HandleFunc("/dashboard", DashboardPage)
	http.HandleFunc("/configure", ConfigurePage)
	http.HandleFunc("/terminal-page/", TerminalPage)
	http.HandleFunc("/terminal/", WebSocketEndpoint)

	// REST APIs
	http.HandleFunc("/api/namespaces", ListNamespacesAPI)
	http.HandleFunc("/api/namespace", CreateNamespaceAPI)
	http.HandleFunc("/api/deployments", CreateDeploymentsAPI)
	http.HandleFunc("/api/services", CreateServicesAPI)
	http.HandleFunc("/api/resources/", DeleteResourceAPI) // Prefix match for dynamic path
}
