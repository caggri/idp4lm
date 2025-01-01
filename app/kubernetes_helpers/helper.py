from kubernetes import client
import app_config


def crete_rbac(namespace: str):
    role = client.V1Role(
        api_version="rbac.authorization.k8s.io/v1",
        kind="Role",
        metadata=client.V1ObjectMeta(
            name="isolated-role",
            namespace=namespace
        ),
        rules=[

            client.V1PolicyRule(
                api_groups=[""],
                resources=["*"],
                verbs=["*"]
            ),

            client.V1PolicyRule(
                api_groups=["apps", "batch", "extensions"],
                resources=["*"],
                verbs=["*"]
            ),

            client.V1PolicyRule(
                api_groups=["rbac.authorization.k8s.io"],
                resources=["roles", "rolebindings"],
                verbs=[""]
            )
        ]
    )

    rbac_api = client.RbacAuthorizationV1Api()
    try:
        rbac_api.create_namespaced_role(namespace, role)
        print(
            f"Role 'full-access-except-rbac' created in namespace '{namespace}'.")
    except client.exceptions.ApiException as e:
        if e.status == 409:
            print(
                f"Role 'full-access-except-rbac' already exists in namespace '{namespace}'.")
        else:
            print(f"Error creating Role: {e}")

    role_binding = client.V1RoleBinding(
        api_version="rbac.authorization.k8s.io/v1",
        kind="RoleBinding",
        metadata=client.V1ObjectMeta(
            name="default-rolebinding",
            namespace=namespace
        ),
        subjects=[
            client.RbacV1Subject(
                kind="ServiceAccount",
                name="default",
                namespace=namespace
            )
        ],
        role_ref=client.V1RoleRef(
            kind="Role",
            name="isolated-role",
            api_group="rbac.authorization.k8s.io"
        )
    )

    try:
        rbac_api.create_namespaced_role_binding(namespace, role_binding)
        print(
            f"RoleBinding 'default-rolebinding' created in namespace '{namespace}'.")
    except client.exceptions.ApiException as e:
        if e.status == 409:
            print(
                f"RoleBinding 'default-rolebinding' already exists in namespace '{namespace}'.")
        else:
            print(f"Error creating RoleBinding: {e}")


def create_quota(namespace: str):
    v1 = client.CoreV1Api()

    resource_quota = client.V1ResourceQuota(
        metadata=client.V1ObjectMeta(
            name="resource-quota",
            namespace=namespace
        ),
        spec=client.V1ResourceQuotaSpec(
            hard={
                "requests.cpu": "2",
                "requests.memory": "4Gi",
                "limits.cpu": "6",
                "limits.memory": "8Gi",
                "persistentvolumeclaims": "5",
                "requests.storage": "100Gi",
                "pods": "10",
                "services": "5",
                "configmaps": "10",
                "secrets": "10",
            }
        )
    )

    try:
        v1.create_namespaced_resource_quota(
            namespace=namespace, body=resource_quota)
        print("ResourceQuota created!")
    except client.exceptions.ApiException as e:
        print(f"Error on creating ResourceQuota: {e}")


def create_network_policy(namespace: str):
    network_instance = app_config.network_v1

    network_policy = client.V1NetworkPolicy(
        api_version="networking.k8s.io/v1",
        kind="NetworkPolicy",
        metadata=client.V1ObjectMeta(
            name="allow-internal-comm",
            namespace=namespace
        ),
        spec=client.V1NetworkPolicySpec(
            pod_selector=client.V1LabelSelector(),
            policy_types=["Ingress"],
            ingress=[
                client.V1NetworkPolicyIngressRule(
                    _from=[
                        client.V1NetworkPolicyPeer(
                            pod_selector=client.V1LabelSelector()
                        )
                    ]
                )
            ]
        )
    )

    try:
        policy_result = network_instance.create_namespaced_network_policy(
            namespace=namespace, body=network_policy)
        return policy_result
    except client.exceptions.ApiException as e:
        return (f"Error on creating Network Policy: {e}")
