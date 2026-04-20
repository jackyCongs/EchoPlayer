package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// DirConfig represents a directory mapping for video resources.
type DirConfig struct {
	Path  string
	Alias string
}

// Pre-configured video directories and their aliases.
var videoConfigs = []DirConfig{
	{Path: `F:\baidunetdisk\young sheldon\Young.Sheldon.S01`, Alias: "Young Sheldon S01"},
	{Path: `F:\baidunetdisk\young sheldon\Young.Sheldon.S02`, Alias: "Young Sheldon S02"},
	{Path: `F:\baidunetdisk\young sheldon\Young.Sheldon.S03`, Alias: "Young Sheldon S03"},
	{Path: `F:\baidunetdisk\S01`, Alias: "Modern Family S01"},
}

// Episode represents a single video episode with its associated subtitle files.
type Episode struct {
	Key   string `json:"key"`
	Video string `json:"video"`
	EnSrt string `json:"enSrt"`
	CnSrt string `json:"cnSrt"`
}

// Category groups episodes by their directory alias.
type Category struct {
	Alias    string    `json:"alias"`
	Episodes []Episode `json:"episodes"`
}

// Regex to extract standard episode identifiers (e.g., S01E01, 1x01).
var keyRegex = regexp.MustCompile(`S\d+E\d+|\d+X\d+`)

const backupFile = "backup.json"

func main() {
	// Map physical video directories to HTTP endpoints.
	for _, config := range videoConfigs {
		folderName := filepath.Base(config.Path)
		http.Handle("/"+folderName+"/", http.StripPrefix("/"+folderName+"/", http.FileServer(http.Dir(config.Path))))
	}

	// Serve the static frontend files (e.g., index.html) from the current working directory.
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Define API routes.
	http.HandleFunc("/api/list", listVideosHandler)
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/api/segments/get", getSegmentsHandler)
	http.HandleFunc("/api/segments/save", saveSegmentsHandler)

	fmt.Println("EchoPlayer Local Server started at http://localhost:8000")

	// Wrap the default multiplexer with CORS middleware.
	log.Fatal(http.ListenAndServe(":8000", corsMiddleware(http.DefaultServeMux)))
}

// getSegmentsHandler reads and returns the saved shadowing segments.
func getSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(backupFile)
	if err != nil {
		if os.IsNotExist(err) {
			w.Write([]byte("{}"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// saveSegmentsHandler overwrites the backup file with new segment data.
func saveSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	err := os.WriteFile(backupFile, body, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// listVideosHandler scans the configured directories and returns a structured playlist.
func listVideosHandler(w http.ResponseWriter, r *http.Request) {
	var categories []Category

	for _, config := range videoConfigs {
		folderName := filepath.Base(config.Path)
		files, err := os.ReadDir(config.Path)
		if err != nil {
			log.Printf("Warning: Cannot read directory %s: %v", config.Path, err)
			continue
		}

		var fileNames []string
		for _, f := range files {
			if !f.IsDir() {
				fileNames = append(fileNames, f.Name())
			}
		}
		sort.Strings(fileNames)

		tempMap := make(map[string]*Episode)
		var keys []string

		for _, name := range fileNames {
			lowName := strings.ToLower(name)
			rawKey := extractKey(name)
			if rawKey == "" {
				rawKey = strings.TrimSuffix(name, filepath.Ext(name))
			}

			fullKey := config.Alias + " " + rawKey

			if _, ok := tempMap[fullKey]; !ok {
				tempMap[fullKey] = &Episode{Key: fullKey}
				keys = append(keys, fullKey)
			}

			// Construct dynamic URLs for videos and subtitles.
			fullURL := fmt.Sprintf("http://%s/%s/%s", r.Host, folderName, name)
			if strings.HasSuffix(lowName, ".mp4") || strings.HasSuffix(lowName, ".mkv") {
				tempMap[fullKey].Video = fullURL
			} else if strings.HasSuffix(lowName, "-en.srt") {
				tempMap[fullKey].EnSrt = fullURL
			} else if strings.HasSuffix(lowName, ".srt") {
				tempMap[fullKey].CnSrt = fullURL
			}
		}

		cat := Category{Alias: config.Alias, Episodes: []Episode{}}
		for _, k := range keys {
			if tempMap[k].Video != "" {
				cat.Episodes = append(cat.Episodes, *tempMap[k])
			}
		}
		categories = append(categories, cat)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// extractKey extracts standard episode identifiers using regex.
func extractKey(name string) string {
	return keyRegex.FindString(strings.ToUpper(name))
}

// corsMiddleware applies necessary CORS headers to all responses.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
