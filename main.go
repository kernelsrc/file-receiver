package main

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const uploadDirName = "~/Downloads/mobile_uploads"

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "127.0.0.1"
}

func expandHome(path string) string {
	home := os.Getenv("HOME")
	if home == "" {
		home = "/tmp"
	}
	return strings.Replace(path, "~", home, 1)
}

func main() {
	home := expandHome(uploadDirName)
	os.MkdirAll(home, 0755)

	port := "8080"
	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:%s", ip, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Connection", "close")
		t := template.Must(template.New("t").Parse(htmlPage))
		t.Execute(w, url)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("  [server] Upload error getting file: %v\n", err)
			http.Error(w, "No file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := header.Filename
		if filename == "" {
			filename = "uploaded_file.bin"
		}
		filename = filepath.Base(filename)

		outPath := filepath.Join(home, filename)
		outFile, err := os.Create(outPath)
		if err != nil {
			fmt.Printf("  [server] Error creating %s: %v\n", outPath, err)
			http.Error(w, "File creation error", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		bytesWritten, err := io.Copy(outFile, file)
		if err != nil {
			fmt.Printf("  [server] Write error: %v\n", err)
		} else {
			fmt.Printf("  [server] Upload OK: %s (%d bytes)\n", filename, bytesWritten)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("===================================================")
	fmt.Println("  Go File Receiver Server")
	fmt.Println("===================================================")
	fmt.Printf("Upload dir:  %s\n", home)
	fmt.Printf("Local IP:    %s\n", ip)
	fmt.Printf("Port:        %s\n", port)
	fmt.Println("QR Code:     qrencode -t ANSIUTF8 '" + url + "'")
	fmt.Println("============================================")
	fmt.Printf("  [server] Server started on %s:%s. Waiting...\n", ip, port)

	http.ListenAndServe(ip+":"+port, mux)
}

const htmlPage = `<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>File Receiver</title>
<style>
body{font-family:sans-serif;margin:0;padding:12px;background:#f0f0f0;text-align:center}
h1{font-size:1.2em;color:#333}
.zone{background:#fff;border:3px dashed #4A90D9;border-radius:10px;padding:40px;margin:12px 0;font-size:.95em;color:#555;cursor:pointer;transition:background .2s}
.zone:hover,.zone:active{background:#e6f2ff}
.status{font-weight:bold;color:#4A90D9;margin-top:10px}
.progress-wrap{width:100%;max-width:400px;margin:10px auto;background:#ddd;border-radius:8px;overflow:hidden}
.progress-bar{height:24px;background:#4A90D9;width:0%;color:#fff;font-size:.85em;line-height:24px;text-align:center}
</style>
</head>
<body>
<h1>Mobile File Upload</h1>
<div class="zone" onclick="document.getElementById('fileInput').click()">Tap here to choose a file from your phone</div>
<input type="file" id="fileInput" name="file" style="display:none" onchange="uploadFile(this)">
<div class="status" id="status">Waiting for file...</div>
<div class="progress-wrap"><div class="progress-bar" id="progressBar">0%</div></div>
<script>
function uploadFile(input){
  var f=input.files[0];
  if(!f){document.getElementById('status').textContent='No file selected';return;}
  document.getElementById('status').textContent='Uploading: '+f.name+'...';
  document.getElementById('progressBar').style.width='0%';
  document.getElementById('progressBar').textContent='0%';
  var xhr=new XMLHttpRequest();
  var fd=new FormData();
  fd.append('file',f);
  xhr.upload.onprogress=function(e){
    if(e.lengthComputable){
      var pct=Math.round((e.loaded/e.total)*100);
      document.getElementById('progressBar').style.width=pct+'%';
      document.getElementById('progressBar').textContent=pct+'%';
    }
  };
  xhr.onload=function(){
    if(xhr.status===200){
      document.getElementById('status').textContent='Success! Uploaded '+f.name;
    } else {
      document.getElementById('status').textContent='Upload failed (HTTP '+xhr.status+')';
    }
  };
  xhr.onerror=function(){
    document.getElementById('status').textContent='Upload network error';
  };
  xhr.open('POST','/upload');
  xhr.send(fd);
}
</script>
</body>
</html>`
