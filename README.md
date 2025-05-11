# infra-health-cli

A handytool, cross-platform CLI tool for performing on-demand health checks across common network layers. Designed for SREs, DevOps engineers, or anyone needing a fast, scriptable way to verify service availability or system status.

---

## 🚀 Features

- ✅ Supports **HTTP, HTTPS, TCP, and ICMP (ping)** checks
- 📊 Optional **JSON output** for logs, pipelines, or integration
- 🧪 Designed for use in **on-call triage**, **CI/CD gates**, or **bootstrap debugging**
- 🖥 Works on **Linux, macOS, and Windows**
- ⚙️ Interactive and flag-based usage modes

---

## 📦 Installation

```bash
go build -o infra-health-cli
```
## 🔧 Usage

### Interactive mode (default)
```
./infra-health-cli
```
### Flag-based (for automation or scripting)
```
./infra-health-cli --choice=1 --endpoint=example.com --json
```

| Flag         | Description                                             |
| ------------ | ------------------------------------------------------- |
| `--choice`   | Monitoring type: `1=HTTP`, `2=HTTPS`, `3=TCP`, `4=ICMP` |
| `--endpoint` | Target address or hostname                              |
| `--port`     | Port number (for TCP checks only)                       |
| `--json`     | Output result as JSON                                   |

##   Example Output
### Console Output (default)
```
HTTPS check successful for example.com (Status: 200)
```
### JSON Output (--json)
```
{
  "type": "HTTPS",
  "endpoint": "example.com",
  "status": "OK",
  "status_code": 200,
  "latency": "35.6ms",
  "timestamp": "2025-05-05T20:33:45Z"
}

```
