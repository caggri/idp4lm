# IDP4LM

**IDP4LM** (Internal Developer Platform LM) is a web application aimed at providing an easy-to-use solution for creating development sandboxes using Kubernetes resources. This project has been fully rewritten in Go for maximum performance and reliability.

## Features

- **Isolated Environments**: Sandboxes, represented as Kubernetes namespaces, provide isolated environments. Creating a sandbox automatically generates RBAC, NetworkPolicy, and ResourceQuota objects.
- **Easy Container Access**: The platform provides native WebSocket integration. Users can drop into a pod's shell directly from their browser (xterm.js).
- **Resource Management**: Every resource created by a user can be viewed, updated, and deleted using the UI.
- **High Performance**: Built with Go standard library `net/http` and `client-go`, resulting in a fast, lightweight, and dependency-minimal backend.

## Prerequisites

- **Go 1.21+**
- Properly configured `kubeconfig` (`~/.kube/config`)
- Network Policy supported CNI (eg. Amazon VPC CNI, Calico, Weave Net)
- ResourceQuota plugin is enabled by default in the latest version of Kubernetes. If you are using an older version or changed the Admission Controller settings, please verify it.

## Installation

1. **Clone the Repository**:
    ```sh
    git clone https://github.com/caggri/idp4lm.git
    cd idp4lm
    ```

2. **Download Dependencies**:
    ```sh
    go mod tidy
    ```

3. **Build the Project**:
    ```sh
    go build -o bin/idp4lm ./cmd/idp4lm
    ```

## Running the Application

1. **Start the Server**:
    ```sh
    ./bin/idp4lm
    ```

2. **Access the Application**:
    Open your browser and navigate to:
    ```
    http://localhost:8080
    ```

## Development Notes

- This project is actively under development. Contributions are welcome via pull requests.
- The web interface uses standard Go `html/template` with Vanilla CSS and JS located in the `web/` directory.

## Planned Features

- Implementing IaC to create an EKS Cluster
- OIDC integration with Keycloak
- Implement a public endpoint for users via Route53 and Ingress

###### *Lifemote Networks Winter Fest 2024 is a 4-day event to learn new things, build prototypes, and work with other teams on new ideas.