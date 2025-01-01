
from fastapi.responses import HTMLResponse
from fastapi import APIRouter, WebSocket, WebSocketDisconnect

from kubernetes.stream import stream
import app_config
from fastapi import WebSocket


import asyncio

router = APIRouter()


@router.websocket("/terminal/{namespace}/{pod_name}")
async def websocket_endpoint(websocket: WebSocket, namespace: str, pod_name: str):
    print(f"WebSocket connection established: {namespace}/{pod_name}")
    await websocket.accept()

    exec_command = ["/bin/sh"]

    try:
        resp = await asyncio.to_thread(
            stream,
            app_config.core_v1.connect_get_namespaced_pod_exec,
            name=pod_name,
            namespace=namespace,
            command=exec_command,
            stderr=True,
            stdin=True,
            stdout=True,
            tty=True,
            _preload_content=False,
        )

        async def kubernetes_to_websocket():
            try:
                while True:
                    stdout_data = await asyncio.to_thread(resp.read_stdout, timeout=1)
                    if stdout_data:
                        await websocket.send_text(stdout_data)

                    stderr_data = await asyncio.to_thread(resp.read_stderr, timeout=1)
                    if stderr_data:
                        await websocket.send_text(stderr_data)
            except Exception as e:
                print(f"Error: {e}")

        async def websocket_to_kubernetes():
            try:
                while True:
                    data = await websocket.receive_text()
                    await asyncio.to_thread(resp.write_stdin, data)
            except WebSocketDisconnect:
                print("WebSocket connection closed.")
            except Exception as e:
                print(f"Error: {e}")

        await asyncio.gather(
            kubernetes_to_websocket(),
            websocket_to_kubernetes(),
        )
    except Exception as e:
        print(f"Error: {e}")
    finally:
        if 'resp' in locals():
            print("Closing exec stream...")
            resp.close()
        await websocket.close()
        print("WebSocket stream closed.")


@router.get("/terminal-page/{namespace}/{pod_name}", response_class=HTMLResponse)
async def terminal_page(namespace: str, pod_name: str):
    return f"""
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Terminal</title>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/xterm/5.5.0/xterm.min.css">
        <style>
                *{{
            height: 100%;
        }}
            #terminal {{
                width: 100%;
                height: 200vh;
                display: flex;
                flex-direction: column;
            }}
        </style>
    </head>
    <body>
        <div id="terminal"></div>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/5.5.0/xterm.js"></script>
        <script>
            document.addEventListener("DOMContentLoaded", () => {{
                const terminal = new Terminal();
                terminal.open(document.getElementById('terminal'));

                let socket = null;

                function setupSocket() {{
                    socket = new WebSocket(`ws://${{window.location.host}}/terminal/{namespace}/{pod_name}`);

                    socket.onopen = () => {{
                        terminal.write(
                            "\\r\\n*** Connecion Established on {pod_name} ***\\r\\n");
                    }};

                    socket.onmessage = (event) => {{
                        terminal.write(event.data);
                    }};

                    socket.onclose = () => {{
                        terminal.write(
                            "\\r\\n*** Connection Closed ***\\r\\n");
                        socket = null;
                    }};

                    socket.onerror = (error) => {{
                        terminal.write(
                            "\\r\\n*** WebSocket Error Occured ***\\r\\n");
                    }};
                }}

                setupSocket();

                terminal.onData((data) => {{
                    if (socket && socket.readyState === WebSocket.OPEN) {{
                        socket.send(data);
                    }} else {{
                        terminal.write(
                            "\\r\\n*** WebSocket connection is not active ***\\r\\n");
                    }}
                }});
            }});
        </script>
    </body>
    </html>
    """
