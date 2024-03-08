package file

import (
    "encoding/json"
    "log"
    "os"
    "path"
    "server/internal/config"
    "server/internal/helper"
    "sync"
)

var (
    Lock           sync.RWMutex
    SharedFiles    = make(map[string]File)
    pendingUploads = make(map[string]File)
)

type File struct {
    Path       string
    ClientId   int32      `json:",omitempty"`
    LastUpdate int64      `json:",omitempty"`
    Action     SyncAction `json:",omitempty"`
}

func (f File) ToJSON() (str string, err error) {
    data, err := json.Marshal(f)
    if err != nil {
        return
    }
    str = string(data)
    return
}

func (f File) LocalLastUpdate() (lastUpdate int64) {
    fileInfo, err := os.Stat(path.Join(config.BoxFolder, f.Path))
    if err != nil {
        helper.LogErr(err)
        return
    }
    // Get the last modified time
    return fileInfo.ModTime().Unix()
}

type SyncAction int

const (
    Download SyncAction = iota + 1
    Upload
)

func NewSharedFile(filePath string) File {
    f := File{Path: filePath}
    // Lock the Map to avoid concurrent access
    Lock.Lock()
    defer Lock.Unlock()
    SharedFiles[filePath] = f
    // Add original clientId to avoid pushing to file to him
    if i, hasItem := pendingUploads[filePath]; hasItem {
        f.ClientId = i.ClientId
    }
    f.Action = Download
    return f
}

func NewPendingUpload(f File) {
    pendingUploads[f.Path] = f
}

func GetFileCount() int {
    Lock.RLock()
    defer Lock.RUnlock()
    return len(SharedFiles)
}

func Load() () {
    log.Println("[Shared box]:loading local files...")
    // Create Box Folder if not exists
    helper.CreateLocalFolder(config.BoxFolder)
    // Read the contents of the folder
    getAll(config.BoxFolder, SharedFiles)
    log.Printf("[Shared Box]:loaded with %d files", len(SharedFiles))
    return
}

func getAll(folderPath string, b map[string]File) {
    fileInfos, err := os.ReadDir(folderPath)
    toCrop := len(config.BoxFolder)
    if folderPath[:2] == "./" {
        // remove ./ from string when cleaning path as it's not part of file relative path
        toCrop -= 2
    }
    if err != nil {
        log.Fatal("Cannot load Box Folder", folderPath, err)
    }
    for _, file := range fileInfos {
        name := file.Name()
        fullPath := path.Join(folderPath, name)
        if file.IsDir() {
            getAll(fullPath, b)
            continue
        }
        var ignoreFile bool
        for _, i := range config.IgnoreFiles {
            if i == name {
                ignoreFile = true
                break
            }
        }
        if ignoreFile {
            continue
        }
        if err != nil {
            helper.LogErr(err)
            continue
        }
        finalPath := fullPath[toCrop:]
        b[finalPath] = File{
            Path: finalPath,
        }
    }
}
