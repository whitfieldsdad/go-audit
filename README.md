# go-audit

## Features

- Detect when processes start<sub>1</sub>/stop
- JSONL output

1<sub>1</sub>. As a non-elevated user, we simply poll the process list every 10 milliseconds. This is surprisingly reliable and efficient on macOS.

1<sub>2</sub>. On Windows, when running as an elevated user, we detect when processes start/stop by tracing [Microsoft-Windows-Kernel-Process](https://github.com/repnz/etw-providers-docs/blob/master/Manifests-Win7-7600/Microsoft-Windows-Kernel-Process.xml) ([{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}](https://github.com/search?q=%7B22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716%7D+language%3AMarkdown&type=code&l=Markdown)) with Event Tracing for Windows (ETW).

## Usage

### Command line

To watch for new processes:

```bash
go run main.go run
```

```json
...
{
  "header": {
    "id": "70fcf683-a749-4dd0-880f-70eaa0d50ebd",
    "time": "2024-02-06T11:16:38.696857-05:00",
    "object_type": "process",
    "event_type": "started"
  },
  "data": {
    "pid": 94067,
    "ppid": 3333,
    "name": "ps",
    "create_time": "2024-02-06T11:16:38.688-05:00",
    "executable": {
      "path": "/bin/ps",
      "filename": "ps",
      "hashes": {
        "md5": "c69d135ec952c1e7e71a6661d7f2c668",
        "sha1": "35e5b335a5858e58728c32de1f5812af87e8a1f4",
        "sha256": "0146891ae982b8ac830beea880da94eb00e7c456820ca54c0f7523a6fbedb096",
        "xxh3": 5261099236022621584
      }
    }
  }
}
...
```

To only select processes that are a descendant of a particular process:

```bash
go run main.go run --ancestor-pid 93707 
```

```json
...
{
  "header": {
    "id": "ba6e3e03-0381-408a-a58e-1925a8e791f2",
    "time": "2024-02-06T11:18:32.319187-05:00",
    "object_type": "process",
    "event_type": "started"
  },
  "data": {
    "pid": 94335,
    "ppid": 93707,
    "name": "whoami",
    "argv": [
      "whoami"
    ],
    "argc": 1,
    "command_line": "whoami",
    "create_time": "2024-02-06T11:18:32.315-05:00",
    "executable": {
      "path": "/usr/bin/whoami",
      "filename": "whoami",
      "hashes": {
        "md5": "9998866dc9ea32e4b4cff7ce737272ab",
        "sha1": "3c62dad9b22c6bf437c4eb4dd73a13175d326575",
        "sha256": "c4167b65515e95be93ecb3cdc555096bb088bccaeb7ee22cc0f817d040761b25",
        "xxh3": 6707975974401128194
      }
    }
  }
}
...
```
