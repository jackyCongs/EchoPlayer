# 🎯 EchoPlayer (IELTS Shadowing Pro)

**A lightweight, local-first web tool designed for intensive language acquisition and frame-accurate shadowing practice.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📖 Overview

As an aspiring immigrant and IELTS candidate, I found existing media players lacked the precision and persistence needed for **Intensive Listening (精听)** and **Shadowing (影子练习)**. This tool allows users to load local media folders, automatically bind subtitles, and save high-value language segments for recurring practice.

### Key Pain Points Solved:
* **Precision Control**: Frame-accurate A-B looping and ±1s fine-tuning.
* **Segment Persistence**: Save specific time-stamps (e.g., a tricky dialogue) to a local "Learning Database" (`localStorage`).
* **Zero Latency**: Local-first architecture allows instant loading of GB-sized media files.
* **Cognitive Load Management**: "Blind Mode" to disable visual cues and force auditory focus.

---

## 🖼️ Showcase

### User Interface
The dashboard features an automated playlist, synchronized bilingual subtitles, and a dedicated segment memory bank.

![UI Showcase](https://raw.githubusercontent.com/congbaochang/EchoPlayer/refs/heads/master/example.png)

### File Organization
The tool uses a regex-based engine to bind `.mp4` video files with their corresponding `.srt` (English and Bilingual) subtitles automatically.

---

## 🚀 Features

* **📦 One-Click Folder Import**: Automatically parses entire seasons of TV shows (e.g., *Modern Family*).
* **🔄 A-B Loop & Memory**: Set [A] and [B] points to loop segments. Save these "cuts" to revisit them later without re-seeking.
* **⌨️ Developer-Friendly Shortcuts**: 
    * `Space`: Play/Pause
    * `Arrow Left/Right`: -1s / +1s
    * `R`: Replay from Point [A]
* **👁️ Blind Mode**: Instantly hide the video to practice "pure" listening.
* **⏩ Variable Speed**: Adjust playback from 0.7x to 1.0x for phonetic analysis.

---

## 🛠️ Technical Implementation & Workflow

### 1. Audio Transcoding (The Codec Challenge)
Browsers often lack support for AC3/DTS audio found in many high-quality `MKV` files. I implemented a lossless video-copy workflow using **FFmpeg** to ensure web compatibility:

```powershell
# PowerShell: Lossless video stream copy + AAC audio conversion
Get-ChildItem *.mkv | ForEach-Object { ffmpeg -i "$($_.Name)" -c:v copy -c:a aac "$($_.BaseName).mp4" }
```

### 2. Subtitle Injection
The tool dynamically converts `.srt` files into `WebVTT` blobs in memory to bypass local file system restrictions and enable real-time subtitle switching.

---

## 📈 Learning Methodology

This tool is optimized for the **4-Step Shadowing Method**:
1.  **Blind Listen**: Use the "Blind Mode" to test initial comprehension.
2.  **Logic Mapping**: Analyze the script and mark the A-B loop for difficult segments.
3.  **Shadowing**: Use the `R` shortcut to repeat the segment until the rhythm matches the native speaker.
4.  **Harvesting**: Use the "Save Segment" feature to build a personal library of "Gold Sentences" for IELTS Speaking preparation.

---
