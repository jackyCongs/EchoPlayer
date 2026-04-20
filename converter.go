package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Define target directories to scan
var targetFolders = []string{
	//`F:\baidunetdisk\S01`,
}

func main() {
	fmt.Println("🚀 EchoPlayer Video Converter Started...")
	for _, folder := range targetFolders {
		err := processFolder(folder)
		if err != nil {
			log.Printf("❌ Error processing folder %s: %v", folder, err)
		}
	}
	fmt.Println("🎉 All conversions completed!")
}

func processFolder(folderPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", folderPath)
	}
	fmt.Printf("\n📂 Scanning folder: %s\n", folderPath)
	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Process files only, skip directories
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		// Process .mkv and .mp4 files (excluding temporary files during conversion)
		if (ext == ".mkv" || ext == ".mp4") && !strings.HasSuffix(info.Name(), "_mobile.mp4") {
			return convertFile(path)
		}
		return nil
	})
}

func convertFile(inputPath string) error {
	dir := filepath.Dir(inputPath)
	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)

	// Temporary and final output paths
	tempOutputPath := filepath.Join(dir, baseName+"_mobile.mp4")
	finalOutputPath := filepath.Join(dir, baseName+".mp4")
	fmt.Printf("🎬 Converting: %s\n", filepath.Base(inputPath))

	// Build FFmpeg command
	cmd := exec.Command("ffmpeg", "-v", "warning", "-stats", "-i", inputPath,
		"-c:v", "copy", // Copy video stream directly for speed and lossless quality
		"-c:a", "aac", // Transcode audio to AAC
		"-ac", "2", // 🏆 Force downmix to stereo (fixes 5.1 channel errors)
		"-b:a", "192k",
		"-y", // Overwrite existing temporary files
		tempOutputPath)

	// Redirect stderr to console to view progress
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg execution failed: %v", err)
	}

	// Upon success: remove original file, then rename the temporary file
	if err := os.Remove(inputPath); err != nil {
		return fmt.Errorf("failed to remove original file: %v", err)
	}

	if err := os.Rename(tempOutputPath, finalOutputPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %v", err)
	}

	fmt.Printf("✅ Success: %s\n", filepath.Base(finalOutputPath))
	return nil
}
