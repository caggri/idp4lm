import app_config


def create_deployment(namespace: str, deployment_type: str, replicas: int):
    try:
        apps_v1 = app_config.apps_v1

        container_configs = {
            "app": {
                "name": "nginx",
                "image": "nginx:latest",
                "port": 80
            },
            "redis": {
                "name": "redis",
                "image": "redis:latest",
                "port": 6379
            },
            "mysql": {
                "name": "mysql",
                "image": "mysql:latest",
                "port": 3306,
                "env": [
                    {
                        "name": "MYSQL_ROOT_PASSWORD",
                        "value": "password123"
                    }
                ]
            }
        }

        config = container_configs[deployment_type]

        container = app_config.kubernetes_client.V1Container(
            name=config["name"],
            image=config["image"],
            ports=[app_config.kubernetes_client.V1ContainerPort(
                container_port=config["port"])],
            resources=app_config.kubernetes_client.V1ResourceRequirements(
                requests={
                    "cpu": "250m",
                    "memory": "128Mi"
                },
                limits={
                    "cpu": "500m",
                    "memory": "768Mi"
                }
            )
        )

        if deployment_type == "mysql" and "env" in config:
            container.env = [
                app_config.kubernetes_client.V1EnvVar(
                    name=env["name"], value=env["value"])
                for env in config["env"]
            ]

        deployment = app_config.kubernetes_client.V1Deployment(
            metadata=app_config.kubernetes_client.V1ObjectMeta(
                name=f"{deployment_type}-deployment",
                namespace=namespace,
                labels={"managed-by": "k8s-idp"}
            ),
            spec=app_config.kubernetes_client.V1DeploymentSpec(
                replicas=replicas,
                selector=app_config.kubernetes_client.V1LabelSelector(
                    match_labels={
                        "app": deployment_type,
                        "managed-by": "k8s-idp"
                    }
                ),
                template=app_config.kubernetes_client.V1PodTemplateSpec(
                    metadata=app_config.kubernetes_client.V1ObjectMeta(
                        labels={
                            "app": deployment_type,
                            "managed-by": "k8s-idp"
                        }
                    ),
                    spec=app_config.kubernetes_client.V1PodSpec(
                        containers=[container]
                    )
                )
            )
        )

        apps_v1.create_namespaced_deployment(
            namespace=namespace,
            body=deployment
        )

        return {
            "status": "success",
            "message": f"{deployment_type} deployment created successfully"
        }

    except Exception as e:
        return {
            "status": "error",
            "message": str(e)
        }


def get_existing_deployments(namespace: str):
    try:
        apps_v1 = app_config.apps_v1
        deployments = apps_v1.list_namespaced_deployment(
            namespace=namespace,
            label_selector="managed-by=k8s-idp"
        )

        existing = {}
        for deployment in deployments.items:
            deployment_type = deployment.metadata.name.replace(
                "-deployment", "")
            existing[deployment_type] = deployment.spec.replicas
        return existing

    except Exception as e:
        print(f"Error getting deployments: {e}")
        return {}


def scale_deployment(namespace: str, deployment_type: str, replicas: int):
    try:
        apps_v1 = app_config.apps_v1
        deployment_name = f"{deployment_type}-deployment"

        deployment = apps_v1.read_namespaced_deployment(
            name=deployment_name,
            namespace=namespace
        )

        deployment.spec.replicas = replicas

        apps_v1.patch_namespaced_deployment(
            name=deployment_name,
            namespace=namespace,
            body=deployment
        )

        return {
            "status": "success",
            "message": f"{deployment_type} scaled to {replicas} replicas"
        }

    except Exception as e:
        return {
            "status": "error",
            "message": str(e)
        }


def get_deployment_replicas(namespace: str, deployment_type: str):
    try:
        apps_v1 = app_config.apps_v1
        deployment_name = f"{deployment_type}-deployment"

        deployment = apps_v1.read_namespaced_deployment(
            name=deployment_name,
            namespace=namespace
        )

        return deployment.spec.replicas
    except Exception as e:
        print(f"Error getting replicas: {e}")
        return 1
