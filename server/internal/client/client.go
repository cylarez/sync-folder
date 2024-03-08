package client

import (
    "fmt"
    "log"
    "net/http"
    "server/internal/file"
    "server/internal/helper"
    "sync"
)

type Client struct {
    Id       int32
    FileChan chan file.File
    writer   http.ResponseWriter
    req      *http.Request
}

var (
    clients      = make(map[int32]*Client)
    lock         sync.RWMutex
    clientIdIncr int32
)

func Register(w http.ResponseWriter, r *http.Request) *Client {
    lock.Lock()
    defer lock.Unlock()
    clientIdIncr++
    cli := &Client{clientIdIncr, make(chan file.File), w, r}
    clients[clientIdIncr] = cli
    return cli
}

func Get(id int32) *Client {
    lock.Lock()
    defer lock.Unlock()
    return clients[id]
}

func GetClientCount() int {
    lock.RLock()
    defer lock.RUnlock()
    return len(clients)
}

func BroadcastNewFile(f file.File, excludeClientId int32) {
    log.Printf("Sending message to %d client(s)", len(clients)-1)
    lock.Lock()
    defer lock.Unlock()
    for id, client := range clients {
        if id == excludeClientId || client.FileChan == nil {
            continue
        }
        client.FileChan <- f
    }
}

// Stream keeps connection open to push updates to client channel
func (c Client) Stream() {
    // Disconnection
    defer func() {
        log.Println("Closing Connection with clientId:", c.Id)
        c.Destroy()
    }()
    // Prepare flusher to push updates
    flusher, ok := c.writer.(http.Flusher)
    if !ok {
        http.Error(c.writer, "Streaming not supported", http.StatusInternalServerError)
        return
    }
    flusher.Flush()
    for {
        select {
        case f := <-c.FileChan:
            var json string
            json, err := f.ToJSON()
            if err != nil {
                helper.LogErr(err)
                continue
            }
            _, err = fmt.Fprintf(c.writer, "%s\n", json)
            if err != nil {
                helper.LogErr(err)
                continue
            }
            flusher.Flush()
        case <-c.req.Context().Done():
            // Client disconnected
            return
        }
    }
}

func (c Client) Destroy() {
    lock.Lock()
    defer lock.Unlock()
    close(c.FileChan)
    delete(clients, c.Id)
}
