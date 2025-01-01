import os
from kubernetes import client, config
from kubernetes.stream import stream

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
STATIC_DIR = os.path.join(BASE_DIR, "static")
TEMPLATES_DIR = os.path.join(BASE_DIR, "templates")


config.load_kube_config()
core_v1 = client.CoreV1Api()
apps_v1 = client.AppsV1Api()
network_v1 = client.NetworkingV1Api()
kubernetes_client = client
kubernetes_stream = stream
