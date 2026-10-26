# Go Mobile File Receiver

## Build / Run
```bash
go build -o server .  # or: go run server.go
./server
```

## Requirements
- Go 1.20+
- Linux (Debian)
- `qrencode` optional (`sudo apt install qrencode`)

## Usage
Launch and open `http://<IP>:8080` from mobile browser on same Wi-Fi.
Files save to `~/Downloads/mobile_uploads`.
