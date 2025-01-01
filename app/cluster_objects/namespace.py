from kubernetes_helpers import helper
import app_config


def check_and_create_namespace(namespace: str):
    try:

        result = create_namespace(namespace)

        helper.create_network_policy(namespace=namespace)
        helper.crete_rbac(namespace=namespace)
        helper.create_quota(namespace=namespace)

        return result
    except Exception as e:
        print(f"Namespace utils error: {str(e)}")
        return {"error": str(e)}


def create_namespace(namespace: str):
    try:
        print(f"Checking if namespace exists: {namespace}")
        exists = check_namespace_exists(namespace)
        print(f"Namespace exists: {exists}")

        if exists:
            print("Returning existing namespace")
            return {
                "status": "success",
                "message": f"Namespace {namespace} already exists",
                "exists": True
            }

        print("Creating new namespace")
        namespace_manifest = app_config.kubernetes_client.V1Namespace(
            metadata=app_config.kubernetes_client.V1ObjectMeta(
                name=namespace,
                labels={"managed-by": "k8s-idp"}
            )
        )

        app_config.core_v1.create_namespace(namespace_manifest)
        print("Namespace created successfully")
        return {
            "status": "success",
            "message": f"Namespace {namespace} created successfully",
            "exists": False
        }
    except Exception as e:
        print(f"Error in create_namespace: {str(e)}")
        return {"status": "error", "message": str(e)}


def check_namespace_exists(namespace: str):
    try:
        app_config.core_v1.read_namespace(namespace)
        return True
    except app_config.kubernetes_client.rest.ApiException as e:
        if e.status == 404:
            return False
        raise e
