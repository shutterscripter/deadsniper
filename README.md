# Deadsniper

CLI tool that finds broken or dead links on a web page. Give it a URL; it scrapes that page only (no recursive crawl), checks every link, and reports dead links (4xx/5xx, soft 404s) and links that block the scraper (403).

## Quick start

**macOS / Linux (one-liner):**

```bash
curl -fsSL https://raw.githubusercontent.com/shutterscripter/deadsniper/main/install.sh | sh
```

Add to PATH if prompted (e.g. in `~/.zshrc` or `~/.bashrc`):

```bash
export PATH="$PATH:$HOME/.local/bin"
```

Then run:

```bash
deadsniper -u https://example.com
```

**Windows:** Download the latest `deadsniper-windows-amd64.exe` from [Releases](https://github.com/shutterscripter/deadsniper/releases), then:

```powershell
.\deadsniper-windows-amd64.exe -u https://example.com
```

## Supported platforms

| Platform        | Install |
|----------------|--------|
| macOS (Intel)  | One-liner above, or [Releases](https://github.com/shutterscripter/deadsniper/releases) → `deadsniper-darwin-amd64` |
| macOS (Apple Silicon) | One-liner above, or Releases → `deadsniper-darwin-arm64` |
| Linux (amd64/arm64)   | One-liner above, or Releases → `deadsniper-linux-amd64` / `deadsniper-linux-arm64` |
| Windows (amd64/arm64) | [Releases](https://github.com/shutterscripter/deadsniper/releases) → `deadsniper-windows-amd64.exe` / `deadsniper-windows-arm64.exe` |

## Usage

```bash
deadsniper -u <URL>
```

**Examples:**

```bash
deadsniper -u https://example.com
deadsniper --url https://mysite.com/page.html
deadsniper --version
```

## Options

| Flag | Short | Description | Default |
|------|--------|-------------|---------|
| `--url` | `-u` | URL of the page to check for dead links | (required) |
| `--verbose` | `-v` | Verbose output | false |
| `--threads` | `-t` | Number of threads | 1 |
| `--delay` | `-d` | Delay between requests (seconds) | 0.5 |
| `--timeout` | `-T` | Request timeout (seconds) | 10 |
| `--output-type` | `-o` | Output: 1=text file, 2=json file | 1 |
| `--help` | `-h` | Help | |
| `--version` | | Print version | |

## Output

- **Dead links** — URLs that returned 4xx/5xx or a “soft 404” (HTTP 200 but page says not found).
- **Blocked (403)** — URLs that returned 403 (server blocks the scraper; may work in a browser).

With `--output-type`:

- `-o 1`: write dead links to `data.txt`
- `-o 2`: write dead links to `data.json`

## How it works

1. Fetches the URL with a browser-like User-Agent.
2. Parses HTML and collects all `<a href="...">` links on that page only.
3. Requests each link and classifies: dead (4xx/5xx or soft 404), blocked (403), or OK (200/301/302/304).

## Build from source

**Single platform:**

```bash
go build -o deadsniper .
# Or: go build -o $(go env GOPATH)/bin/deadsniper .
```

**All platforms (local):**

```bash
chmod +x build.sh
./build.sh
```

Binaries go to `dist/`. Requires Go 1.23+.

## Releasing (maintainers)

Releases are built by [GitHub Actions](.github/workflows/release.yml) when you push a version tag:

```bash
git tag v0.1.4
git push origin v0.1.4
```

Binaries are published to [GitHub Releases](https://github.com/shutterscripter/deadsniper/releases); the install script uses the latest release.

## License

See [LICENSE](LICENSE) in the repository.
