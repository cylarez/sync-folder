package route

import (
    "bufio"
    "bytes"
    "net/http"
    "net/http/httptest"
    "os"
    "path"
    "server/internal/client"
    "server/internal/config"
    "server/internal/file"
    "testing"
)

func TestUpload(t *testing.T) {
    // Create a file
    filePath := "test-upload.txt"
    f, err := os.CreateTemp("", filePath)
    writer := bufio.NewWriter(f)
    // Write the content to the file
    content := "Hello World!"
    _, err = writer.WriteString(content)
    err = writer.Flush()
    b, err := os.ReadFile(f.Name())

    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer(b))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    req.Header.Set(config.HeaderFilePath, filePath)

    // Register client
    c := client.Register(rr, req)
    fi := file.File{
        Path:     filePath,
        ClientId: c.Id,
        Action:   0,
    }
    file.NewPendingUpload(fi)

    Upload(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("wrong status code: got %v want %v", status, http.StatusForbidden)
    }
    uploaded, err := os.ReadFile(path.Join(config.BoxFolder, filePath))

    if string(uploaded) != content {
        t.Errorf("unexpected file content: got %v want %v", string(uploaded), content)
    }
    c.Destroy()
}
