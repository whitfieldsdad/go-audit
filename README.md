# go-audit

## Features

- Detect when processes start<sub>1</sub>/stop

1<sub>1</sub>. As a non-elevated user, we simply poll the process list every 10 milliseconds. This is surprisingly reliable and efficient on macOS.

1<sub>2</sub>. On Windows, when running as an elevated user, we detect when processes start/stop by tracing [Microsoft-Windows-Kernel-Process](https://github.com/repnz/etw-providers-docs/blob/master/Manifests-Win7-7600/Microsoft-Windows-Kernel-Process.xml) ([{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}](https://github.com/search?q=%7B22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716%7D+language%3AMarkdown&type=code&l=Markdown)) with Event Tracing for Windows (ETW).
