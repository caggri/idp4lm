
import app_config


def get_managed_resources(namespace: str):
    try:
        v1 = app_config.core_v1
        apps_v1 = app_config.apps_v1

        pods = v1.list_namespaced_pod(
            namespace=namespace,
            label_selector="managed-by=k8s-idp"
        )
        pod_info = [{
            "name": pod.metadata.name,
            "ready": all(status.ready for status in pod.status.container_statuses) if pod.status.container_statuses else False
        } for pod in pods.items]

        services = v1.list_namespaced_service(
            namespace=namespace,
            label_selector="managed-by=k8s-idp"
        )
        service_names = [svc.metadata.name for svc in services.items]

        deployments = apps_v1.list_namespaced_deployment(
            namespace=namespace,
            label_selector="managed-by=k8s-idp"
        )
        deployment_info = [{
            "name": dep.metadata.name,
            "ready": dep.status.ready_replicas == dep.status.replicas if dep.status.ready_replicas else False
        } for dep in deployments.items]

        return {
            "pods": pod_info,
            "services": service_names,
            "deployments": deployment_info
        }

    except Exception as e:
        print(f"Error getting resources: {e}")
        return {"error": str(e)}
