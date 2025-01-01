import app_config


def check_and_create_service(namespace: str, service_name: str):

    v1 = app_config.core_v1

    service_ports = {
        'app': 80,
        'mysql': 3306,
        'redis': 6379
    }

    port = service_ports.get(service_name, 80)

    service_manifest = {
        "apiVersion": "v1",
        "kind": "Service",
        "metadata": {
                "name": service_name + "-service",
                "namespace": namespace,
                "labels": {
                    "app": service_name,
                    "managed-by": "k8s-idp"
                }
        },
        "spec": {
            "selector": {
                "app": service_name
            },
            "ports": [
                {
                    "protocol": "TCP",
                    "port": port,
                    "targetPort": port
                }
            ],
            "type": "ClusterIP"
        }
    }

    try:
        v1.create_namespaced_service(
            namespace=namespace,
            body=service_manifest
        )

        return (f"Service '{service_name}' created successfully'.")

    except:
        v1.read_namespaced_service(name=service_name, namespace=namespace)
        return (f"error. Service '{service_name}' already exists in namespace '{namespace}'.")


def get_existing_services(namespace: str):
    try:
        v1 = app_config.core_v1
        services = v1.list_namespaced_service(
            namespace=namespace,
            label_selector="managed-by=k8s-idp"
        )

        existing = []
        for service in services.items:
            existing.append(service.metadata.name)
        return existing

    except Exception as e:

        print(f"Error getting services: {e}")
        return {}
