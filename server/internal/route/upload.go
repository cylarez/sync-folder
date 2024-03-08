package route

import (
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "server/internal/client"
    "server/internal/config"
    "server/internal/file"
    "server/internal/helper"
)

func Upload(w http.ResponseWriter, r *http.Request) {
    // Retrieve the filename from the custom header
    filePath := r.Header.Get(config.HeaderFilePath)
    fullPath := config.BoxFolder + filePath

    // Create folder
    err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
    if err != nil {
        helper.LogErr(err)
        return
    }
    // Open a new file for writing with the provided filename
    fi, err := os.Create(fullPath)
    if err != nil {
        helper.LogErr(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer fi.Close()

    // Copy the whole request body to the file
    _, err = io.Copy(fi, r.Body)
    if err != nil {
        helper.LogErr(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Printf("[Upload]:completed file: %s", filePath)

    // Store File to Shared Box
    f := file.NewSharedFile(filePath)

    // Sync new content to other clients
    client.BroadcastNewFile(f, f.ClientId)
}
