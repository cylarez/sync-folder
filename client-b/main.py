import http.client
import json
import logging
import os
import time
import urllib.parse
import argparse
import signal
import sys
import file
import helper
from config import BOX_FOLDER_PATH

# Program args and config
parser = argparse.ArgumentParser(description='Client Demo Python')
parser.add_argument('-help', action='help', help='Show this help message and exit')
parser.add_argument('-host', help='Server Hostname | default to 127.0.1')
parser.add_argument('-port', help='Server Port | default to 8080')
parser.add_argument('-api-key', help='API Key | default to 12345')
args = parser.parse_args()

SERVER_HOST = helper.get_arg(args.host, "127.0.0.1")
SERVER_PORT = helper.get_arg(args.port, "8080")
SERVER_API_KEY = helper.get_arg(args.api_key, "12345")

# Create a connection to the server
headers = {
    'Connection': 'keep-alive',
    'Authorization': SERVER_API_KEY
}

# Configure Logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger()


def sync():
    # Prepare local content
    content = file.get_content()
    data = json.dumps(content)
    # Run Sync Request
    conn = http.client.HTTPConnection(SERVER_HOST, int(SERVER_PORT))
    logger.info(f"[Sync]:starting...")
    try:
        conn.request("POST", "/sync", data, headers)
    except ConnectionError:
        reconnect()
    # Get the response
    response = conn.getresponse()
    # Check if the connection was successful
    if response.status != 200:
        error = response.read().decode("utf-8")
        logger.error(f'[Sync]:failed to connect: {response.status} {response.reason} - {error}')
        return
    logger.info(f"[Sync]:sent {len(content['files'])} files")
    logger.info("[Sync]:waiting for updates...")
    while True:
        line = response.readline().strip()
        if not line:
            logger.error('Connection closed by server')
            break  # Exit the loop when connection is closed
        try:
            f = json.loads(line)
        except ValueError:
            logger.error(f'[Sync]:error cannot handle Stream response: {line}')
            return
        if f['Action'] == file.SyncAction.DOWNLOAD.value:
            download(f)
        elif f['Action'] == file.SyncAction.UPLOAD.value:
            upload(f)
        else:
            logger.error(f'[Sync]:error unsupported action for file: {f}')
    # Retry to connect
    conn.close()
    reconnect()


def reconnect():
    logger.error("[Sync]:error reconnecting in 5sec...")
    time.sleep(5)
    sync()


def download(f):
    conn = http.client.HTTPConnection(SERVER_HOST, int(SERVER_PORT))
    file_path = f['Path']
    logger.info(f"[Download]:starting {file_path}...")

    # Ensure local folder exists
    full_path = os.path.join(BOX_FOLDER_PATH, file_path)
    folder = os.path.dirname(full_path)
    file.init_folder(folder)

    # Fetch the file
    encoded_path = urllib.parse.quote(file_path)
    conn.request("GET", f"/download/{encoded_path}", headers=headers)
    response = conn.getresponse()

    if response.status == http.HTTPStatus.OK:
        with open(full_path, 'wb') as f:
            # Read the response data in chunks and write it to the local file
            while True:
                chunk = response.read(8192)  # Read 8KB at a time
                if not chunk:
                    break
                f.write(chunk)
        logger.info(f"[Download]:completed file:{file_path}")
    else:
        logger.info(f"[Download]:error status:{response.status}")
    conn.close()


def upload(f):
    logger.info(f"[Upload]:sending file: {f['Path']}")
    # Load file content
    full_path = os.path.join(BOX_FOLDER_PATH, f['Path'])
    f['Data'] = file.read_content(full_path)
    # Send Upload request
    conn = http.client.HTTPConnection(SERVER_HOST, int(SERVER_PORT))
    conn.request("POST", "/upload", f['Data'], {**headers, **{'File-Path': f['Path']}})
    response = conn.getresponse()
    response_content = str(response.read().decode('utf-8'))

    conn.close()
    if response.status == http.HTTPStatus.OK:
        logger.info(f"[Upload]:completed {f['Path']}")
    else:
        logger.error(f"[Upload]:error status:{response.status} {response_content}")


def signal_handler(sig, frame):
    logger.info("Terminating...")
    sys.exit(0)


def main():
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    file.init_folder(BOX_FOLDER_PATH)
    sync()


main()
