# JSMON CLI

JSMON CLI is a command-line interface for interacting with the jsmon.sh web application. It provides a convenient way to access various features of JSMON directly from your terminal.

## Features

- Upload URLs for scanning
- Rescan previously scanned URLs
- Upload and scan files
- View scan results
- Manage domains
- Set up and manage cron jobs for automated scanning
- Compare JavaScript responses
- View user profile and usage information

## Installation

Ensure you have Golang installed on your system. If not, download and install it from golang.org.
Then, clone the repository and build the binary:
```
git clone https://github.com/rashahacks/jsmon-cli
cd jsmon-cli
go mod download
go build -o jsmon-cli
```
Alternatively, you can install it directly using:
```
go install github.com/rashahacks/jsmon-cli@latest
```

## Usage

Run the CLI tool using:
`jsmon-cli [flags]`

### Flags

- `-u string`: URL to upload for scanning
- `-d string`: Domain to automate scan
- `-f string`: File to upload (local path)
- `-key string`: API key for authentication
- `-jsi string`: Get all automation results
- `-urls`: View all URLs
- `-us int`: Number of URLs to fetch
- `-w string`: Comma-separated list of words to include in the scan
- `-domains`: Get all domains for the user
- `-H string`: Add custom headers in the format 'Key: Value' (can be used multiple times)
- `-profile`: View user profile
- `-files`: View all files
- `-query string`: string = <field=apiPaths domain=example.com page=1 sub=true>

## Authentication

The CLI uses an API key for authentication. You can provide the API key using the `-key` flag or by storing it in `~/.jsmon/credentials`.

## Example Commands

1. Upload a URL for scanning:
```jsmon-cli -u https://example.com/main.js -wksp <WORKSPACE_ID>```

2. Upload a file containing JS URLs:
```jsmon-cli -f jsurls.txt -wksp <WORKSPACE_ID>```

3. Scan a domain, subdomain or URL:
```jsmon-cli -d <sub.example.com> -wksp <WORKSPACE_ID>```

4. View user profile:
```jsmon-cli -profile```

5. Query Data from your account:
```
jsmon-cli -query field=apiPaths -wksp <WORKSPACE_ID>
jsmon-cli -query field=extractedUrls -wksp <WORKSPACE_ID>
jsmon-cli -query field=extractedDomains -wksp <WORKSPACE_ID>
jsmon-cli -query field=emails -wksp <WORKSPACE_ID>
jsmon-cli -query field=apiPaths domain=example.com page=2 sub=true> -wksp <WORKSPACE_ID>
```

## Query Guide

Learn more about -query flags here via query guide: <a href="https://knowledge.jsmon.sh/query-data/query-guide">https://knowledge.jsmon.sh/query-data/query-guide</a>

## License

This project is licensed under the Apache License, Version 2.0 (the "License").

You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
