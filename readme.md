
# IDP4LM

**IDP4LM** (Identity Provider for Lightweight Management) is a web application aimed at providing an easy-to-use solution for creating Python development sandboxes using Kubernetes resources. This ongoing project was initiated during **Winter Fest**.

## Features

- **Isolated Environments**: Sandboxes, represented as Kubernetes namespaces, provide isolated environments. Creating a sandbox automatically generates RBAC, NetworkPolicy, and ResourceQuota objects.
- **Easy Conatainer Access**: The platform provides Websocket Ingtegration, with WS users can drop-in pods' shell using browser.
- **Resource Management**: Every resource created by a user can be viewed, updated (some extent), and deleted using the UI.

## Prerequisites

- Python 3.8+
- Properly configured `kubeconfig` 
- Network Policy supported CNI (eg. Amazon VPC CNI, Calico, Weave Net)
- ResourceQuota plugin is enabled by default in the latest version of Kubernetes. If you are using an older version or changed the Admission Controller settings, verify it.

## Installation

1. **Clone the Repository**:
    ```sh
    git clone https://github.com/caggri/idp4lm.git
    cd idp4lm/app
    ```

2. **Install Dependencies**:
    ```sh
    pip install -r requirements.txt
    ```

3. **Configure Kubernetes**:
    Ensure your Kubernetes config file (`~/.kube/config`) is correctly set up with appropriate cluster access.

## Running the Application

1. **Start the Server**:
    ```sh
    uvicorn app.main:app --reload
    ```

2. **Access the Application**:
    Open your browser and navigate to:
    ```
    http://127.0.0.1:8000
    ```

## Development Notes

- This project is actively under development. Contributions are welcome via pull requests.


## Planned Features

- Implementing IaC to create an EKS Cluster
- OIDC integration with Keycloak
- Implement a public endpoint for users via Route53 and Ingress

###### *Lifemote Networks Winter Fest 2024 is a 4 day event to learn new things, build prototypes, and work with other teams on new ideas.