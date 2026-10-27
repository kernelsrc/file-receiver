# Go Mobile File Receiver

## Build / Run
```bash
go build -o server .
./server
```

## Global Install (`go install`)
```bash
go install github.com/user/file-receiver@latest
# or manually:
mkdir -p ~/.local/bin
go build -o ~/.local/bin/file-receiver .
export PATH="$HOME/.local/bin:$PATH"
```

Then just type `file-receiver` from any terminal.

## Requirements
- Go 1.20+
- Linux (Debian)
- `qrencode` optional (`sudo apt install qrencode`)

## Usage
Launch and open `http://<IP>:8080` from mobile browser on same Wi-Fi.
Files save to `~/Downloads/mobile_uploads`.
