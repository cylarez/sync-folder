package helper

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatal("Missing Environment Variable ", key)
	}
	return value
}

func LogErr(err error) {
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[Error] at %s:%d %v", filename, line, err)
}

func CreateLocalFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			log.Panic(err)
		}
	}
}

func StrToInt(str string) int32 {
	if i, err := strconv.Atoi(str); err == nil {
		return int32(i)
	}
	return 0
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// HashFileContent return md5 hash with very small memory allocating
func HashFileContent(path string) (md5Hash string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	var n int
	defer file.Close()
	hash := md5.New()
	// Use a buffer to read the file in chunks
	buffer := make([]byte, 4096)
	for {
		n, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			break
		}
		hash.Write(buffer[:n])
	}
	// Sum the hash and return as a hexadecimal string
	hashSum := hash.Sum(nil)
	md5Hash = hex.EncodeToString(hashSum)
	err = nil
	return
}

func ParseMarkdownLine(line string, openedBlock *bool) string {
	if strings.HasPrefix(line, "```") { // Code block in colour
		*openedBlock = !*openedBlock
		if *openedBlock {
			return "\x1b[34m"
		} else {
			return "\x1b[0m"
		}
	} else if strings.HasPrefix(line, "# ") {
		return "\033[1m" + line[2:] + "\033[0m" // Bold
	} else if strings.HasPrefix(line, "## ") {
		return "\n\033[4m" + line[3:] + "\033[0m" // Underline
	} else if strings.HasPrefix(line, "### ") {
		return "\033[1m" + line[4:] + "\033[0m" // Bold
	} else {
		return line
	}
}
