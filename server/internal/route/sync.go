package route

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "server/internal/client"
    "server/internal/file"
    "server/internal/helper"
    "time"
)

type syncContent struct {
    Files []file.File
}

func Sync(w http.ResponseWriter, r *http.Request) {
    // Set headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Transfer-Encoding", "chunked")
    // Register the client
    c := client.Register(w, r)
    // Check Content
    body, err := io.ReadAll(r.Body)
    if err != nil {
        helper.LogErr(err)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    newContent := syncContent{}
    err = json.Unmarshal(body, &newContent)
    if err != nil {
        helper.LogErr(err)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    var (
        streamUpdates []file.File
        toUpload      int
        toDownload    int
    )
    file.Lock.RLock()
    alreadySynced := make(map[string]bool)
    for _, f := range newContent.Files {
        // Check file exist and up to date (same modified date)
        if _, hasFile := file.SharedFiles[f.Path]; hasFile && f.LastUpdate <= f.LocalLastUpdate() {
            alreadySynced[f.Path] = true
        } else {
            log.Printf("[Sync]:pushing file upload %s to clientId:%d", f.Path, c.Id)
            toUpload++
            f.Action = file.Upload
            streamUpdates = append(streamUpdates, f)
            f.ClientId = c.Id
            file.NewPendingUpload(f)
        }
    }
    // Send files to be downloaded
    for key, f := range file.SharedFiles {
        if !alreadySynced[key] {
            log.Printf("[Sync]:pushing file download %s to clientId:%d", f.Path, c.Id)
            toDownload++
            f.Action = file.Download
            streamUpdates = append(streamUpdates, f)
        }
    }
    file.Lock.RUnlock()
    log.Printf("[Sync] ClientId:%d has %d to upload / %d to download", c.Id, toUpload, toDownload)
    go func() {
        // Wait to make sure client chan is ready
        time.Sleep(100 * time.Millisecond)
        for _, f := range streamUpdates {
            c.FileChan <- f
        }
    }()
    // Keeps connection open waiting for updates
    c.Stream()
}
