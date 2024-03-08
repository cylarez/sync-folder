import os
from config import BOX_FOLDER_PATH, IGNORE_FILES
from enum import Enum


class SyncAction(Enum):
    DOWNLOAD = 1
    UPLOAD = 2


def get_content():
    content = []
    files = get_files(BOX_FOLDER_PATH)
    for f in files:
        # Remote absolute path
        path = f[len(BOX_FOLDER_PATH)+1:]
        # Get Last Update time
        last_update = int(os.path.getmtime(f))
        content.append({'path': path, 'lastUpdate': last_update})
    return {'files': content}


def init_folder(folder):
    if not os.path.exists(folder):
        os.makedirs(folder)


def get_files(directory):
    files = []
    # Iterate over all entries in the directory
    for entry in os.listdir(directory):
        if entry in IGNORE_FILES:
            continue
        full_path = os.path.join(directory, entry)
        # Check if the entry is a file
        if os.path.isfile(full_path):
            files.append(full_path)
        # If the entry is a directory, recursively call get_files
        elif os.path.isdir(full_path):
            files.extend(get_files(full_path))
    return files


def read_content(filename):
    try:
        with open(filename, 'rb') as file:
            content = file.read()
            return content
    except FileNotFoundError:
        return None

