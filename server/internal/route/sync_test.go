package route

import (
    "bytes"
    "fmt"
    "net/http"
    "net/http/httptest"
    "os"
    "path"
    "path/filepath"
    "server/internal/client"
    "server/internal/config"
    "server/internal/file"
    "strings"
    "testing"
    "time"
)

func TestSync(t *testing.T) {
    tmp := os.TempDir()
    tmpBoxDir := path.Join(tmp, "box")
    err := os.MkdirAll(tmpBoxDir, 0750)
    // Cleanup files from previous tests
    err = removeAllFilesInFolder(tmpBoxDir)
    config.BoxFolder = tmpBoxDir + "/"
    newFile := "test-file-1.txt"
    // Add a file in the Box
    _, err = os.Create(config.BoxFolder + newFile)
    file.Load()
    // Make a Body
    content := `{"Files":[{"Path":"test-file-1.txt","Hash":"d41d8cd98f00b204e9800998ecf8427e"},{"Path":"test-file-2.txt","Hash":""}]}`
    body := []byte(content)
    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/sync", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()

    // Register client
    c := client.Register(rr, req)

    go Sync(rr, req)
    time.Sleep(time.Second)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
    }
    // test-file-1.txt should not be in the list since server already has it
    expected := `{"Path":"test-file-2.txt","Action":2}`
    if strings.TrimSpace(rr.Body.String()) != expected {
        t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
    }
    req.Context().Done()
    c.Destroy()
}

func removeAllFilesInFolder(folderPath string) error {
    // Open the folder
    dir, err := os.Open(folderPath)
    if err != nil {
        return err
    }
    defer dir.Close()
    // Read all files in the folder
    files, err := dir.Readdir(-1)
    if err != nil {
        return err
    }
    // Iterate over the files and remove them
    for _, file := range files {
        filePath := filepath.Join(folderPath, file.Name())
        err = os.Remove(filePath)
        if err != nil {
            return err
        }
        fmt.Println("Removed:", filePath)
    }
    return nil
}
