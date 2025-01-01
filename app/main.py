from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from routes import route, ws
import app_config

from kubernetes import config

app = FastAPI()

app.mount("/static", StaticFiles(directory=app_config.STATIC_DIR), name="static")

try:
    config.load_kube_config()
except Exception as e:
    print(f"Warning: Could not load kube config: {e}")

app.include_router(route.router, tags=["all_routes"])
app.include_router(ws.router, tags=["websockets"])

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host="127.0.0.1", port=8000, log_level="ERROR")
