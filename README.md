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
go build -o jsmon
```
Alternatively, you can install it directly using:
```
go install github.com/rashahacks/jsmon-cli@latest
```

## Usage
Run the CLI tool using:
`./jsmon [flags]`

### Flags

- `-u string`: URL to upload for scanning
- `-d string`: Domain to automate scan
- `-f string`: File to upload (local path)
- `-key string`: API key for authentication
- `-fid string`: File ID to scan
- `-jsi string`: Get all automation results
- `-secrets`: Get scanner results
- `-urls`: View all URLs
- `-us int`: Number of URLs to fetch
- `-w string`: Comma-separated list of words to include in the scan
- `-domains`: Get all domains for the user
- `-H string`: Add custom headers in the format 'Key: Value' (can be used multiple times)
- `-profile`: View user profile
- `-files`: View all files
- `-emails string`: View all Emails for specified domains
- `-buckets string`: Get all S3 Domains for specified domains
- `-ips string`: Get all IPs for specified domains
- `-domainUrls string`: Get Domain URLs for specified domains
- `-apis string`: Get API paths for specified domains
- `-query string`: string = <field=apiPaths domain=example.com page=1 sub=true>
- `-compare string`: Compare two js responses by jsmon_ids (format: JSMON_ID1,JSMON_ID2)

## Authentication

The CLI uses an API key for authentication. You can provide the API key using the `-apikey` flag or by storing it in `~/.jsmon/credentials`.

## Example Commands

1. Upload a URL for scanning:
```./jsmon -u https://example.com/main.js -wksp <WORKSPACE_ID>```

2. Upload a file containing JS URLs:
```./jsmon -f jsurls.txt -wksp <WORKSPACE_ID>```

3. Scan a domain, subdomain or URL:
```./jsmon -d <sub.example.com> -wksp <WORKSPACE_ID>```

4. View user profile:
```./jsmon -profile```

5. Query Data from your account:
```
./jsmon -query <field=apiPaths> -wksp <WORKSPACE_ID>
./jsmon -query <field=extractedUrls> -wksp <WORKSPACE_ID>
./jsmon -query <field=extractedDomains> -wksp <WORKSPACE_ID>
./jsmon -query <field=emails> -wksp <WORKSPACE_ID>
./jsmon -query <field=apiPaths domain=example.com page=2 sub=true> -wksp <WORKSPACE_ID>
```

## Contributors and Maintainers

- Inderjeet Singh
  - [GitHub](https://github.com/encodedguy)
  - [LinkedIn](https://www.linkedin.com/in/encodedguy/)
  - [Twitter](https://x.com/3nc0d3dGuY)

- Nadeem Ahmad
  - [GitHub](https://github.com/Nadeem-Ahmad-25)
  - [LinkedIn](https://www.linkedin.com/in/ndmxx0001/)
  - [Twitter](https://x.com/NadeemAhmad97)

- Nandini Kumari
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
