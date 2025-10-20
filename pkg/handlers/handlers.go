package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"log"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func HandleMain(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(res, req, "index.html")
}

func HandleUpload(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	res.Header().Set("Content-Type", "text/plain")
	
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := req.FormFile("myFile")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(res, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(res, "Error reading file", http.StatusInternalServerError)
		return
	}

	input := strings.TrimSpace(string(content))
	if input == "" {
		log.Printf("Empty file content: %v", err)
		http.Error(res, "Empty file content", http.StatusBadRequest)
		return
	}

	result, err := service.Convert(input) 
	if err != nil {
		log.Printf("Convertion error: %v", err)
		http.Error(res, fmt.Sprintf("Convertion error: %v", err), http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".txt"
	}

	outputFilename := "result_" + time.Now().UTC().Format("20060102150405") + ext

	outputFile,err := os.Create(outputFilename)
	if err != nil {
		log.Printf("Error creating result file: %v", err)
		http.Error(res, "Error saving result", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(result)
	if err != nil {
		log.Printf("Error writing result: %v", err)
		http.Error(res, "Error saving result", http.StatusInternalServerError)
		return
	}
	

	fmt.Fprint(res, result)
	// res.WriteHeader(http.StatusOK)
	// fmt.Fprint(res, result)
}

func sanitizeHTML(s string) string {
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
	return r.Replace(s)
}