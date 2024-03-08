package config

import (
    "time"
)

const (
    HeaderApiKey        = "Authorization"
    HeaderFilePath      = "File-Path"
    PrintMemorySchedule = 20 * time.Second
)

var (
    ApiKey      string
    IgnoreFiles = []string{".DS_Store"}
    BoxFolder   = "./box/"
)
