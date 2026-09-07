# AOF CLI

AOF CLI is the official repository bootstrap utility for adopting the [AI Orchestration Framework (AOF)](https://github.com/aof-framework/aof).

It configures project context and developer-facing adoption artifacts. It is **not** an AOF runtime and is not required by applications after initialization.

---

## Installation

### Option 1: Go Install (Recommended)

If you have Go installed on your machine:

```bash
go install github.com/aof-framework/aof-cli/cmd/aof@latest
```

> Ensure `$GOPATH/bin` (or `%USERPROFILE%\go\bin` on Windows) is in your system's `PATH`.

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/aof-framework/aof-cli.git
cd aof-cli

# Build the binary
go build -o aof ./cmd/aof
```

Move the compiled `aof` binary to a directory in your `PATH`:
- **Linux/macOS:** `sudo mv aof /usr/local/bin/`
- **Windows:** Move `aof.exe` to a directory included in your system `PATH`.

### Option 3: Run directly without installation

```bash
cd /path/to/your-project
go run github.com/aof-framework/aof-cli/cmd/aof@latest init
```

---

## Commands

| Command | Description |
| :--- | :--- |
| `aof init [flags]` | Initialize AOF adoption in the current repository |
| `aof version` | Display AOF CLI version and AOF specification baseline |
| `aof help` | Show help information |

---

## Usage Examples

### 1. Interactive Initialization (Wizard)

Run `aof init` inside your target project directory:

```bash
cd itsm-backend
aof init
```

The CLI will automatically detect the project language/manifest and guide you through selecting the AOF profile (`core`, `governed`, `assured`) and AI use cases.

### 2. Automated / Non-Interactive (CI/CD)

```bash
aof init \
  --name itsm-backend \
  --language go \
  --type backend \
  --domain itsm \
  --profile core \
  --ai \
  --ai-analysis \
  --ai-recommendation \
  --non-interactive
```

#### Flags for `aof init`:

- `--name <string>`: Project name (default: current directory name)
- `--language <string>`: Primary language (auto-detected if omitted)
- `--type <string>`: Project type (e.g., `backend`, `frontend`, `fullstack`, `service`, `library`)
- `--domain <string>`: Problem domain (e.g., `itsm`, `fintech`, `ecommerce`)
- `--profile <string>`: AOF profile (`core`, `governed`, `assured`, default: `core`)
- `--features <string>`: Comma-separated list of project features
- `--ai`: Declare AI usage in the project
- `--ai-analysis`: Declare AI analysis capability
- `--ai-classification`: Declare AI classification capability
- `--ai-recommendation`: Declare AI recommendation capability
- `--ai-consequential`: Declare AI consequential execution capability
- `--non-interactive`: Disable interactive prompts

---

## Generated Artifacts

Running `aof init` creates:
- `.aof/`: Machine-readable project definition (`project.yaml`), adoption config (`config.yaml`), and provenance manifest (`manifest.json`).
- `.agents/`: Pinned agent skills (`.agents/skills/`) and agent role definitions (`.agents/agents/`).
- `AGENTS.md` & `AOF.md`: Repository instructions and governance posture for humans and AI coding assistants.
- `docs/aof/`: Embedded static AOF semantic and execution boundary documentation.

---

## Canonical AOF

- Repository: https://github.com/aof-framework/aof
- Specification: https://github.com/aof-framework/aof/blob/main/specification/AOF-v1.0-Framework-Specification.md

`aof init` is offline-capable and does not clone the mutable upstream `main` branch. Static bootstrap guidance is embedded in the binary and records canonical upstream provenance in generated project metadata.
