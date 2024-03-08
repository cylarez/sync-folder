package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"server/internal/client"
	"server/internal/config"
	"server/internal/file"
	"server/internal/helper"
	"server/internal/middleware"
	"server/internal/route"
	"time"
)

func main() {
	// Check Help
	help := flag.Bool("help", false, "Help command")
	flag.Parse()
	if *help {
		printHelp()
		return
	}
	// Use http.FileServer to serve clients files to /download
	var fs = http.StripPrefix("/download", http.FileServer(http.Dir(config.BoxFolder)))

	// Setup routes
	http.Handle("/sync", middleware.Logger(middleware.Auth(route.Sync)))
	http.Handle("/upload", middleware.Logger(middleware.Auth(route.Upload)))
	http.Handle("/download/", middleware.Logger(middleware.Auth(fs.ServeHTTP)))

	// Load Box Files in memory
	file.Load()

	// Load API Key
	config.ApiKey = helper.MustGetEnv("API_KEY")

	// Periodic Log Status (user count, files, memory)
	go runLogStatus()

	// Starting the server
	addr := helper.GetEnv("SERVER_ADDRESS", ":8080")
	log.Println("[Server]:started on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		helper.LogErr(err)
	}
}

func runLogStatus() {
	for {
		logStatus()
		time.Sleep(config.PrintMemorySchedule)
	}
}

func logStatus() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	clientCount := client.GetClientCount()
	fileCount := file.GetFileCount()
	format := "[Status] Connected Clients: %d | Shared Files: %d | Allocated Memory: %v MB | Total Allocated Memory: %v MB"
	log.Printf(format, clientCount, fileCount, helper.BToMb(m.Alloc), helper.BToMb(m.TotalAlloc))
}

func printHelp() {
	// Open the README file
	file, err := os.Open("./README.md")
	if err != nil {
		helper.LogErr(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	openBlock := false
	for scanner.Scan() {
		line := scanner.Text()
		renderedLine := helper.ParseMarkdownLine(line, &openBlock)
		if renderedLine != "" {
			fmt.Println(renderedLine)
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}
