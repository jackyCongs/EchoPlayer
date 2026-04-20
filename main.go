package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// 你指定的最新配置
type DirConfig struct {
	Path  string
	Alias string
}

var videoConfigs = []DirConfig{
	{Path: `F:\baidunetdisk\young sheldon\Young.Sheldon.S01`, Alias: "Young Sheldon S01"},
	{Path: `F:\baidunetdisk\young sheldon\Young.Sheldon.S02`, Alias: "Young Sheldon S02"},
	{Path: `F:\baidunetdisk\young sheldon\Young.Sheldon.S03`, Alias: "Young Sheldon S03"},
	{Path: `F:\baidunetdisk\S01`, Alias: "Modern Family S01"},
}

type Episode struct {
	Key   string `json:"key"`
	Video string `json:"video"`
	EnSrt string `json:"enSrt"`
	CnSrt string `json:"cnSrt"`
}

type Category struct {
	Alias    string    `json:"alias"`
	Episodes []Episode `json:"episodes"`
}

var keyRegex = regexp.MustCompile(`S\d+E\d+|\d+X\d+`)

const backupFile = "backup.json"

func main() {
	for _, config := range videoConfigs {
		folderName := filepath.Base(config.Path)
		http.Handle("/"+folderName+"/", http.StripPrefix("/"+folderName+"/", http.FileServer(http.Dir(config.Path))))
	}

	http.HandleFunc("/api/list", listVideosHandler)
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	})

	// 数据同步接口
	http.HandleFunc("/api/segments/get", getSegmentsHandler)
	http.HandleFunc("/api/segments/save", saveSegmentsHandler)

	fmt.Println("EchoPlayer Resource Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
}

func getSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(backupFile)
	if err != nil {
		if os.IsNotExist(err) {
			w.Write([]byte("{}"))
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func saveSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", 405)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	// 直接覆盖写入，确保数据最新
	err := ioutil.WriteFile(backupFile, body, 0644)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func listVideosHandler(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	for _, config := range videoConfigs {
		folderName := filepath.Base(config.Path)
		files, _ := os.ReadDir(config.Path)
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

			// 使用你指定的 Alias 拼接唯一 Key
			fullKey := config.Alias + " " + rawKey

			if _, ok := tempMap[fullKey]; !ok {
				tempMap[fullKey] = &Episode{Key: fullKey}
				keys = append(keys, fullKey)
			}

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

func extractKey(name string) string {
	match := keyRegex.FindString(strings.ToUpper(name))
	return match
}

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
