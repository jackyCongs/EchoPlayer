# 🎯 EchoPlayer (IELTS Shadowing Pro)

**A professional-grade, local-first web application engineered for IELTS candidates to perform intensive shadowing and phonetic analysis across all devices in a local network.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)

## 📖 Project Vision

EchoPlayer was designed to solve the "Persistence Gap" in language learning. Unlike standard players, it treats shadowing as a data-driven process—allowing you to mark, save, and sync tricky dialogue segments across your PC and mobile devices.

### Why this architecture?
* **Centralized Logic**: Your Go server handles the heavy lifting (file serving, data persistence), keeping the client (browser) lightweight and fast.
* **LAN Synchronization**: Save a segment on your PC, and it's immediately available on your iPad or phone for reviews.
* **Codec Freedom**: The built-in Go converter ensures high-quality MKV files play smoothly on mobile Safari/Chrome with full audio support.

---

## 🛠️ Setup & Workflow (The 4-Step Start)

### 1. Configuration (Path Setup)
Open `main.go` and map your local media directories to the server. This allows the Go backend to index your files:
```go
var videoConfigs = []DirConfig{
    {Path: `F:\BaiduNetdisk\Modern.Family.S01`, Alias: "Modern Family S01"},
    {Path: `F:\BaiduNetdisk\Young.Sheldon.S01`, Alias: "Young Sheldon S01"},
}
```

### 2. Media Normalization (The Converter)
Mobile browsers often fail to play AC3/5.1 audio found in MKV files. Run our specialized Go tool to ensure 100% compatibility:
```bash
go run converter.go
```
* **Performance**: Uses lossless `-c:v copy` for video (near-instant speed).
* **Fix**: Transcodes audio to **Stereo AAC** (`-ac 2`), fixing "no sound" issues on iPhones/tablets.

### 3. Start the Server
Compile and launch the backend. This will host the web interface and the REST APIs:
```bash
go run main.go
```

### 4. Access the Player (Connection Guide)
The server runs on port **8080**. You can access it from any device on your Wi-Fi:

* **On your Host PC**: Open [http://localhost:8080](http://localhost:8080)
* **On Mobile/Tablet (LAN)**:
  1. Find your PC's IP (run `ipconfig` in CMD, e.g., `192.168.31.100`).
  2. Open your device browser and go to: `http://192.168.31.100:8080`
     *Note: The app will automatically remember this IP for future sessions.*

---

## 🚀 Key Features

* **🔄 Segment Memory**: Set A-B points and click "Save Segment". The timestamps are persisted to `backup.json` on the server.
* **👁️ Blind Mode**: Instantly dim the video to 0% brightness to force total auditory focus.
* **⌨️ Hotkeys**:
  * `Space`: Play/Pause
  * `Arrow Left/Right`: ±1s Seek
  * `R`: Replay from Point [A]
* **📱 Responsive Design**: A tailored CSS layout that works perfectly on vertical phone screens and wide desktop monitors.

---

## 📈 Learning Methodology

This tool is optimized for the **4-Step Shadowing Method**:
1.  **Blind Listen**: Use the "Blind Mode" to test initial comprehension.
2.  **Logic Mapping**: Analyze the script and mark the A-B loop for difficult segments.
3.  **Shadowing**: Use the `R` shortcut to repeat the segment until your rhythm matches the native speaker.
4.  **Harvesting**: Save segments to build a personal library of "Gold Sentences" for IELTS Speaking preparation.

---

## 🖼️ Showcase

### User Interface
The dashboard features an automated playlist, synchronized bilingual subtitles, and a dedicated segment memory bank.

![UI Showcase](https://raw.githubusercontent.com/congbaochang/EchoPlayer/refs/heads/master/example.png)

## 👨‍💻 Author
**Jacky Cong**
*Backend Developer | Golang Expert | IELTS Candidate*