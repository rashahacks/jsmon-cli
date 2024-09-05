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

```
git clone https://github.com/rashahacks/jsmon-cli
cd jsmon-cli
go build -o jsmon
```

## Usage
`./jsmon [flags]`

### Flags

- `-scanUrl string`: URL or scan ID to rescan
- `-uploadUrl string`: URL to upload for scanning
- `-apikey string`: API key for authentication
- `-scanFile string`: File ID to scan
- `-uploadFile string`: File to upload (local path)
- `-automationData string`: Get all automation results
- `-scannerData`: Get scanner results
- `-size int`: Number of results to fetch (default 10000)
- `-cron string`: Set, update, or stop cronjob
- `-notifications string`: Set cronjob notification channel
- `-time int64`: Set cronjob time
- `-vulnerabilitiesType string`: Set type of cronjob (URLs, Analysis, Scanner)
- `-domains string`: Set domains for cronjob
- `-domainsNotify string`: Set notify (true/false) for each domain for cronjob
- `-urls`: View all URLs
- `-urlSize int`: Number of URLs to fetch
- `-scanDomain string`: Domain to automate scan
- `-words string`: Comma-separated list of words to include in the scan
- `-getDomains`: Get all domains for the user
- `-H string`: Custom headers in the format 'Key: Value' (can be used multiple times)
- `-usage`: View user profile
- `-files`: View all files
- `-Emails string`: View all Emails for specified domains
- `-S3Domains string`: Get all S3 Domains for specified domains
- `-ips string`: Get all IPs for specified domains
- `-DomainUrls string`: Get Domain URLs for specified domains
- `-api string`: Get the APIs for specified domains
- `-compare string`: Compare two js responses by jsmon_ids (format: JSMON_ID1,JSMON_ID2)

## Authentication

The CLI uses an API key for authentication. You can provide the API key using the `-apikey` flag or by storing it in `~/.jsmon/credentials`.

## Examples

1. Upload a URL for scanning:
```./jsmon -uploadUrl https://example.com```

2. Rescan a previously scanned URL:
```./jsmon -scanUrl SCAN_ID```

3. Upload and scan a local file:
```./jsmon -uploadFile /path/to/file.js```

4. View user profile:
```./jsmon -usage```

5. Compare two JavaScript responses:
```./jsmon -compare JSMON_ID1,JSMON_ID2```

## Contributors and Maintainers

- Inderjeet Singh
  - [GitHub](https://github.com/encodedguy)
  - [LinkedIn](https://www.linkedin.com/in/encodedguy/)
  - [Twitter](https://x.com/3nc0d3dGuY)

- Nadeem Ahmad
  - [GitHub](https://github.com/Nadeem-Ahmad-25)
  - [LinkedIn](https://www.linkedin.com/in/ndmxx0001/)
  - [Twitter](https://x.com/NadeemAhmad97)

- Nandani Kumari
  - [GitHub](https://github.com/nandini-56)
  - [LinkedIn](https://www.linkedin.com/in/nandini-kumari-693257296/)
  - [Twitter](https://x.com/Nandini17060041)
 
- Akash Litoriya
  - [GitHub](https://github.com/akashlitoriya)
  - [LinkedIn](https://www.linkedin.com/in/akashlitoriya/)
  - [Twitter](https://x.com/akashhhh_l)

## License

This project is licensed under the Apache License, Version 2.0 (the "License").

You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
