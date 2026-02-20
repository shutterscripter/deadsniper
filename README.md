# Deadsniper

A CLI tool that finds broken or dead links on a web page. It scrapes the given URL, collects all links on that page only (no recursive crawling), and checks each link. It reports dead links (4xx/5xx, soft 404s) and links that block the scraper (403).

## Requirements

- Go 1.25+

## Install

**One-liner (macOS / Linux, after a release exists):**

```bash
curl -fsSL https://raw.githubusercontent.com/shutterscripter/deadsniper/main/install.sh | sh
```

Then ensure `~/.local/bin` is in your `PATH` (the script will remind you). Run:

```bash
deadsniper -u https://example.com
```

**From source (single platform):**

```bash
go build -o deadsniper .
# Or install to $GOPATH/bin:
go build -o $(go env GOPATH)/bin/deadsniper .
```

**Build all platforms locally (no GoReleaser):**

```bash
chmod +x build.sh
./build.sh
```

Binaries are written to `dist/` (see table below). On Windows, download the `.exe` from [GitHub Releases](https://github.com/shutterscripter/deadsniper/releases) and run it.

## Usage

```bash
deadsniper -u <URL>
```

Example:

```bash
deadsniper -u https://example.com
```

The tool fetches the page at the given URL, finds all `<a href="...">` links on that page, and checks each linked URL. Results are printed to stdout and optionally written to a file.

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

- **Dead links:** URLs that returned 4xx/5xx or a “soft 404” (HTTP 200 but page content indicates not found).
- **Blocked by bot (403):** URLs that returned 403 Forbidden (server blocks the scraper; link may work in a browser).

When `--output-type` is set, dead links are also written to a file:

- `-o 1`: `data.txt` (one URL per line)
- `-o 2`: `data.json` (one JSON string per line)

## How it works

1. Fetches the given URL with a browser-like User-Agent and headers.
2. Parses the HTML and collects links only from that page (does not follow links to other pages).
3. Requests each collected link and classifies the response:
   - **Dead:** status 4xx/5xx (except 403), or HTTP 200 with body that looks like a 404 page (e.g. “404”, “not found”).
   - **Blocked:** status 403 (not counted as dead).
4. 304 Not Modified and 200/301/302 are treated as success; 403 is reported separately as “blocked by server”.

## Publishing releases (Mac, Linux, Windows)

**Option A: GoReleaser (recommended)**

1. Install GoReleaser once: `go install github.com/goreleaser/goreleaser@latest`
2. Create a tag and push: `git tag v0.1.0 && git push origin v0.1.0`
3. Run: `goreleaser release`
4. Binaries appear on [GitHub Releases](https://github.com/shutterscripter/deadsniper/releases). The `install.sh` script downloads from `latest`.

**Option B: Manual**

1. Run `./build.sh` to populate `dist/`.
2. Create a new release on GitHub, tag (e.g. `v0.1.0`), and upload the files from `dist/` as assets.

Asset names expected by `install.sh`: `deadsniper-darwin-amd64`, `deadsniper-darwin-arm64`, `deadsniper-linux-amd64`, `deadsniper-linux-arm64`. Windows: `deadsniper-windows-amd64.exe` (users download manually or via scoop).

## License

See the LICENSE file in the repository.
