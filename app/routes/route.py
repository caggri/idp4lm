from fastapi import APIRouter
from fastapi.responses import HTMLResponse, JSONResponse
from fastapi.templating import Jinja2Templates
import app_config
from fastapi import Request

from cluster_objects.resources import get_managed_resources
from cluster_objects.deployments import get_existing_deployments, create_deployment, scale_deployment
from cluster_objects.namespace import check_and_create_namespace
from cluster_objects.services import check_and_create_service, get_existing_services

router = APIRouter()

templates = Jinja2Templates(directory=app_config.TEMPLATES_DIR)


@router.get("/", response_class=HTMLResponse)
async def landing(request: Request):
    return templates.TemplateResponse("landing.html", {"request": request})


@router.get("/api/namespaces")
async def list_namespaces():
    try:
        v1 = app_config.core_v1
        namespaces = v1.list_namespace()
        managed_namespaces = []

        for ns in namespaces.items:
            if ns.metadata.labels and ns.metadata.labels.get("managed-by") == "k8s-idp":
                managed_namespaces.append({
                    "name": ns.metadata.name,
                    "created": ns.metadata.creation_timestamp.strftime("%Y-%m-%d %H:%M:%S")
                })

        return {"namespaces": managed_namespaces}

    except Exception as e:
        print(f"Error listing namespaces: {e}")
        return JSONResponse(
            status_code=500,
            content={"error": str(e)}
        )


@router.post("/api/namespace")
async def create_namespace(request: Request):
    try:
        print("Namespace creation request received")
        data = await request.json()
        namespace = data.get("namespace")
        print(f"Requested namespace: {namespace}")

        print(f"Checking namespace: {namespace}")
        result = check_and_create_namespace(namespace)
        print(f"Check result: {result}")

        if "error" in result:
            print(f"Error in result: {result['error']}")
            return JSONResponse(
                status_code=400,
                content={"error": result["error"]}
            )

        print(f"Success: {result}")
        return JSONResponse(content=result)

    except Exception as e:
        print(f"Unexpected error: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={"error": str(e)}
        )


@router.post("/api/services")
async def create_services(request: Request):
    try:
        print("Service creation request received")
        data = await request.json()
        namespace = data.get("namespace")
        servicess = data.get("services")
        for x in servicess:

            result = check_and_create_service(
                namespace=namespace, service_name=x.get("type"))

            if "error" in result:

                print(f"Error in result: {result['error']}")

        return result

    except Exception as e:

        print(f"Unexpected error: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={"error": str(e)}
        )


@router.post("/api/deployments")
async def create_deployments(request: Request):
    try:
        data = await request.json()
        namespace = data.get("namespace")
        deployments = data.get("deployments")

        if not namespace or not deployments:
            return JSONResponse(
                status_code=400,
                content={"error": "Namespace and deployments are required"}
            )

        for deployment in deployments:
            result = create_deployment(
                namespace=namespace,
                deployment_type=deployment["type"],
                replicas=deployment["replicas"]
            )
            if "error" in result:
                return JSONResponse(
                    status_code=400,
                    content={"error": f"Failed to create {
                        deployment['type']} deployment: {result['error']}"}
                )

        return JSONResponse(content={"message": "Deployments created successfully"})

    except Exception as e:
        print(f"Error creating deployments: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={"error": str(e)}
        )


@router.get("/configure")
async def configure(request: Request, ns: str):
    try:
        existing_workloads = get_existing_deployments(ns)
        existing_services = get_existing_services(ns)

        return templates.TemplateResponse(
            "configure.html",
            {
                "request": request,
                "namespace": ns,
                "existing_deployments": existing_workloads,
                "existing_services": existing_services
            }
        )
    except Exception as e:
        print(f"Configure error: {e}")
        return templates.TemplateResponse(
            "error.html",
            {"request": request, "error": str(e)}
        )


@router.delete("/api/resources/{type}/{name}")
async def delete_resource(type: str, name: str, ns: str):
    try:
        v1 = app_config.core_v1
        apps_v1 = app_config.apps_v1

        if type == "pod":
            v1.delete_namespaced_pod(name=name, namespace=ns)
        elif type == "service":
            v1.delete_namespaced_service(name=name, namespace=ns)
        elif type == "deployment":
            apps_v1.delete_namespaced_deployment(name=name, namespace=ns)
        elif type == "namespace":
            v1.delete_namespace(name=name)
        else:
            return JSONResponse(
                status_code=400,
                content={"error": f"Unknown resource type: {type}"}
            )

        return JSONResponse(content={"status": "success"})

    except Exception as e:
        print(f"Error deleting resource: {e}")
        return JSONResponse(
            status_code=500,
            content={"error": str(e)}
        )


@router.get("/dashboard", response_class=HTMLResponse)
async def dashboard(request: Request, ns: str):
    try:
        resources = get_managed_resources(ns)

        if "error" in resources:
            return templates.TemplateResponse(
                "error.html",
                {"request": request, "error": resources["error"]}
            )

        return templates.TemplateResponse(
            "dashboard.html",
            {
                "request": request,
                "namespace": ns,
                "resources": resources
            }
        )
    except Exception as e:
        return templates.TemplateResponse(
            "error.html",
            {"request": request, "error": str(e)}
        )
