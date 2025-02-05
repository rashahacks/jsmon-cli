package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type stringSliceFlag []string

type Workspace struct {
	WkspId string `json:"wkspId"`
	Name   string `json:"name"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}

func getWorkspaces() ([]Workspace, error) {
	endpoint := fmt.Sprintf("%s/workspaces", apiBaseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("X-Jsmon-Key", strings.TrimSpace(getAPIKey()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("unexpected response: %s", string(body))
		}
		return nil, fmt.Errorf("API key error: %s", errorResp.Message)
	}

	var workspaces []Workspace
	err = json.Unmarshal(body, &workspaces) // Unmarshal directly into array
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return workspaces, nil
}

func displayWorkspaces() error {
	workspaces, err := getWorkspaces()
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		return fmt.Errorf("no workspaces found")
	}

	fmt.Println("Available Workspaces:")
	for _, ws := range workspaces {
		fmt.Printf("%s (ID: %s)\n", ws.Name, ws.WkspId)
	}
	fmt.Println("\nUse -wksp <workspace_id> to specify a workspace")
	return nil
}
func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func updateCLI() error {
	fmt.Println("Updating jsmon-cli to the latest version...")

	cmd := exec.Command("go", "install", "github.com/rashahacks/golang-api-client-dev@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to update jsmon-cli: %v", err)
	}

	fmt.Println("Successfully updated jsmon-cli to the latest version!")
	return nil
}

func main() {
	uploadUrl := flag.String("u", "", "URL to upload for scanning")
	apiKeyFlag := flag.String("key", "", "API key for authentication")
	updateFlag := flag.Bool("ud", false, "Update jsmon-cli to the latest version")
	scanFileId := flag.String("fid", "", " File to be rescanned by fileId.")
	uploadFile := flag.String("f", "", "File to upload giving path to the file locally.")
	getAllResults := flag.String("jsi", "", "View JS Intelligence Data by domain name")
	size := flag.Int("s", 100, "Number of results to fetch (default 100)")
	fileTypes := flag.String("type", "", "files type (e.g. pdf,txt)")
	workspaceFlag := flag.String("wksp", "", "Workspace ID")
	listWorkspacesFlag := flag.Bool("workspaces", false, "List all available workspaces")
	getScannerResultsFlag := flag.Bool("secrets", false, "View Keys & Secrets by domain name")
	// cron := flag.String("cron", "", "Set cronjob.")
	// cronNotification := flag.String("channel", "", "Set cronjob notification channel.")
	// cronTime := flag.Int64("time", 0, "Set cronjob time.")
	// cronType := flag.String("type", "", "Set type[URLs, Analysis, Scanner] of cronjob.")
	// cronDomains := flag.String("domains", "", "Set domains for cronjob.")
	// cronDomainsNotify := flag.String("domainsNotify", "", "Set notify(true/false) for each domain for cronjob.")
	viewurls := flag.Bool("urls", false, "view all urls")
	//viewurlsSize := flag.Int("us", 10, "Number of URLs to fetch")
	scanDomainFlag := flag.String("d", "", "Domain to automate scan")
	wordsFlag := flag.String("w", "", "Comma-separated list of words to include in the scan")
	urlswithmultipleResponse := flag.Bool("curls", false, "View changed JS URLs.")
	getDomainsFlag := flag.Bool("domains", false, "Get all domains for the user.")

	var headers stringSliceFlag
	flag.Var(&headers, "H", "Custom headers in the format 'Key: Value' (can be used multiple times)")

	addCustomWordsFlag := flag.String("addCustomWords", "", "add custom words to the scan")
	usageFlag := flag.Bool("profile", false, "View user profile")
	viewfiles := flag.Bool("files", false, "view all files")
	viewEmails := flag.String("emails", "", "Get all emails for specified domains.")
	s3domains := flag.String("buckets", "", "get all S3Domains for specified domains")
	ips := flag.String("ips", "", "Get all IPs for specified domains")
	gql := flag.String("gqls", "", "Get graph QL operations")
	domainUrl := flag.String("domainUrls", "", "Get Domain URLs for specified domains")
	apiPath := flag.String("apis", "", "Get the APIs for specified domains")
	fileExtensionUrls := flag.String("extUrls", "", "Get URLs containing any file type.")
	socialMediaUrls := flag.String("socialUrls", "", "Get URLs for social media sites.")
	domainStatus := flag.String("domainStatuses", "", "Get the availabilty of domains")
	queryParamsUrls := flag.String("queryParamUrls", "", "Get URLs containing query params for specified domain.")
	localhostUrls := flag.String("localUrls", "", "Get URLs with localhost in the hostname.")
	filteredPortUrls := flag.String("portUrls", "", "Get URLs with port numbers in the hostname")
	s3DomainsInvalid := flag.String("bucketTakeovers", "", "Get available S3 domains (404 status).")
	compareFlag := flag.String("compare", "", "Compare two js responses by jsmon_ids (format: JSMON_ID1,JSMON_ID2)")
	reverseSearchResults := flag.String("rsearch", "", "Specify the input type (e.g., emails, domainname)")
	createWordListFlag := flag.String("wordlist", "", "creates a new word list from domains")
	searchUrlsByDomainFlag := flag.String("urlsByDomain", "", "Search URLs by domain")
	getResultByJsmonId := flag.String("jsiJsmonId", "", "Get JS Intelligence for the jsmon ID.")
	getResultByFileId := flag.String("jsiFileId", "", "Get JS Intelligence for the file ID.")
	totalAnalysisDataFlag := flag.Bool("count", false, "total count of overall analysis data")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  %s [flags]\n\n", os.Args[0])
		fmt.Println("Flags:")

		fmt.Fprintf(os.Stderr, "\nINPUT:\n")
		fmt.Fprintf(os.Stderr, "  -u <URL>          		URL to upload for scanning.\n")
		fmt.Fprintf(os.Stderr, "  -fid <fileId>         	File to be rescanned by fileId.\n")
		fmt.Fprintf(os.Stderr, "  -f <local file path>          File to upload (local path)\n")
		fmt.Fprintf(os.Stderr, "  -d <domainName>   		Domain to scan\n")

		fmt.Fprintf(os.Stderr, "\nAUTHENTICATION:\n")
		fmt.Fprintf(os.Stderr, "  -key <uuid>                   API key for authentication\n")

		fmt.Fprintf(os.Stderr, "\nUTILITY:\n")
		fmt.Fprintf(os.Stderr, "  -ud                           Update jsmon-cli to the latest version\n")

		fmt.Fprintf(os.Stderr, "\nOUTPUT:\n")
		fmt.Fprintf(os.Stderr, "  -jsi <domainName>             View JS Intelligence data by domain name\n")
		fmt.Fprintf(os.Stderr, "  -secrets                      View Keys & Secrets\n")
		fmt.Fprintf(os.Stderr, "  -urls                         View all URLs.\n")
		fmt.Fprintf(os.Stderr, "  -us int                       Number of URLs to fetch (default 10).\n")
		fmt.Fprintf(os.Stderr, "  -files                        View all files.\n")
		fmt.Fprintf(os.Stderr, "  -type <types>                 Specify file types (e.g., pdf,txt), use ',' as separator.\n")
		fmt.Fprintf(os.Stderr, "  -profile                      View user profile.\n")
		fmt.Fprintf(os.Stderr, "  -curls                        View changed JS URLs.\n")

		fmt.Fprintf(os.Stderr, "\nADDITIONAL OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  -H <Key: Value>               Custom headers (can be used multiple times).\n")
		fmt.Fprintf(os.Stderr, "  -w <words>                    Comma-separated list of words to include in the scan.\n")
		fmt.Fprintf(os.Stderr, "  -domains                      Get all domains for the user.\n")
		fmt.Fprintf(os.Stderr, "  -emails <domain>              Get all emails for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -buckets <domain>             Get all S3 domains for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -ips <domain>                 Get all IPs for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -domainUrls <domain>          Get domain URLs for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -apis <domain>                Get API paths for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -extUrls <domain>             Get URLs containing any file type.\n")
		fmt.Fprintf(os.Stderr, "  -socialUrls <domain>          Get URLs for social media sites.\n")
		fmt.Fprintf(os.Stderr, "  -domainStatuses <domain>      Get availability status of domains.\n")
		fmt.Fprintf(os.Stderr, "  -queryParamUrls <domain>      Get URLs containing query params for specified domain.\n")
		fmt.Fprintf(os.Stderr, "  -localUrls <domain>           Get URLs with localhost in the hostname.\n")
		fmt.Fprintf(os.Stderr, "  -portUrls <domain>            Get URLs with port numbers in the hostname.\n")
		fmt.Fprintf(os.Stderr, "  -bucketTakeovers <domain>     Get available S3 domains (404 status).\n")
		fmt.Fprintf(os.Stderr, "  -urlsByDomain <domain>        Search URLs by domain.\n")
		fmt.Fprintf(os.Stderr, "  -compare <ID1,ID2>            Compare two JS responses by IDs (format: ID1,ID2).\n")
		fmt.Fprintf(os.Stderr, "  -gqls <domain>                Get GraphQL operations for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -count                        Get total count of overall analysis data.\n")
		fmt.Fprintf(os.Stderr, "  -jsiJsmonId <ID>              Get automation results by jsmon ID.\n")
		fmt.Fprintf(os.Stderr, "  -jsiFileId <ID>               Get automation results by file ID.\n")

		// Automation results section
		fmt.Fprintf(os.Stderr, "\nReverse JS search:\n")
		fmt.Fprintf(os.Stderr, "  -rsearch <field>=<value>      Search by field: emails, domainname, extracteddomains, s3domains, url, extractedurls, ipv4addresses, ipv6addresses, jwttokens, gqlquery, gqlmutation, guids, apipaths, vulnerabilities, nodemodules, domainstatus, queryparamsurls, socialmediaurls, filterdporturls, gqlfragment, s3domainsinvalid, fileextensionurls, localhosturls.\n")

	}

	flag.Parse()
	if *updateFlag {
		if err := updateCLI(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}
	if *apiKeyFlag != "" {
		setAPIKey(*apiKeyFlag)
	} else {
		err := loadAPIKey()
		if err != nil {
			fmt.Println("Error loading API key:", err)
			fmt.Println("Please provide an API key using the -apikey flag.")
			os.Exit(1)
		}
	}

	if flag.NFlag() == 0 || (flag.NFlag() == 1 && *apiKeyFlag != "") {
		fmt.Println("No action specified. Use -h or --help for usage information.")
		flag.Usage()
		os.Exit(1)
	}

	if *listWorkspacesFlag {
		err := displayWorkspaces()
		if err != nil {
			fmt.Printf("Error listing workspaces: %v\n", err)
			os.Exit(1)
		}
		return
	}

	switch {
	case *scanFileId != "":
		scanFileEndpoint(*scanFileId)
	case *uploadFile != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		uploadFileEndpoint(*uploadFile, headers, *workspaceFlag)
	case *viewurls:
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}

		err := viewUrls(*size, *workspaceFlag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case *viewfiles:
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		viewFiles(*workspaceFlag)
	case *uploadUrl != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		err := uploadUrlEndpoint(*uploadUrl, headers, *workspaceFlag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	// case *rescanDomainFlag != "":
	// 	rescanDomain(*rescanDomainFlag)
	case *totalAnalysisDataFlag:
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		totalAnalysisData(*workspaceFlag)
	case *searchUrlsByDomainFlag != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		searchUrlsByDomain(*searchUrlsByDomainFlag, *workspaceFlag)
	case *urlswithmultipleResponse:
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		urlsmultipleResponse(*workspaceFlag)
	case *viewEmails != "":
		domains := strings.Split(*viewEmails, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getEmails(domains, *workspaceFlag)
	case *getResultByJsmonId != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAutomationResultsByJsmonId(strings.TrimSpace(*getResultByJsmonId), *(workspaceFlag))
	case *reverseSearchResults != "":
		parts := strings.SplitN(*reverseSearchResults, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid format for reverseSearchResults. Use field=value format.")
			return
		}

		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAutomationResultsByInput(field, value, *workspaceFlag)

	case *getResultByFileId != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAutomationResultsByFileId(strings.TrimSpace(*getResultByFileId), *workspaceFlag)
	case *s3domains != "":
		domains := strings.Split(*s3domains, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getS3Domains(domains, *workspaceFlag)
	case *ips != "":
		domains := strings.Split(*ips, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllIps(domains, *workspaceFlag)
	case *gql != "":
		domains := strings.Split(*gql, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getGqlOps(domains, *workspaceFlag)
	case *domainUrl != "":
		domains := strings.Split(*domainUrl, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getDomainUrls(domains, *workspaceFlag)
	case *apiPath != "":
		domains := strings.Split(*apiPath, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getApiPaths(domains, *workspaceFlag)
	case *getScannerResultsFlag:
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getScannerResults(*workspaceFlag)
	// case *scanUrl != "":
	// 	rescanUrlEndpoint(*scanUrl)
	// case *cron == "start":
	// 	StartCron(*cronNotification, *cronTime, *cronType, *cronDomains, *cronDomainsNotify)
	// case *cron == "stop":
	// 	StopCron()
	case *getDomainsFlag:
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getDomains(*workspaceFlag)

	case *fileExtensionUrls != "":
		extensions := strings.Split(*fileTypes, ",")
		for i, extension := range extensions {
			extensions[i] = strings.TrimSpace(extension)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllFileExtensionUrls(*fileExtensionUrls, extensions, *size, *workspaceFlag)
	case *domainStatus != "":
		// domains := strings.Split(*domainStatus, ",")
		// for i, domain := range domains {
		// 	domains[i] = strings.TrimSpace(domain)
		// }
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllDomainsStatus(*domainStatus, *size, *workspaceFlag)

	case *socialMediaUrls != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllSocialMediaUrls(*socialMediaUrls, *size, *workspaceFlag)
	case *queryParamsUrls != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllQueryParamsUrls(*queryParamsUrls, *size, *workspaceFlag)
	case *localhostUrls != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllLocalhostUrls(*localhostUrls, *size, *workspaceFlag)
	case *filteredPortUrls != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllFilteredPortUrls(*filteredPortUrls, *size, *workspaceFlag)
	case *s3DomainsInvalid != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		getAllS3DomainsInvalid(*s3DomainsInvalid, *size, *workspaceFlag)
	case *compareFlag != "":
		ids := strings.Split(*compareFlag, ",")
		if len(ids) != 2 {
			fmt.Println("Invalid format for compare. Use: JSMON_ID1,JSMON_ID2")
			os.Exit(1)
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		compareEndpoint(strings.TrimSpace(ids[0]), strings.TrimSpace(ids[1]), *workspaceFlag)
	// case *cron == "update":
	// 	UpdateCron(*cronNotification, *cronType, *cronDomains, *cronDomainsNotify, *cronTime)
	case *getAllResults != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}

		err := getAllAutomationResults(*getAllResults, *size, *workspaceFlag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case *scanDomainFlag != "":
		words := []string{}
		if *wordsFlag != "" {
			words = strings.Split(*wordsFlag, ",")
		} else {
			rootWord := extractRootWord(*scanDomainFlag)
			if rootWord != "" {
				words = []string{rootWord}
			}
		}
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Printf("Domain: %s, Words: %v\n", *scanDomainFlag, words)

		err := automateScanDomain(*scanDomainFlag, words, *workspaceFlag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case *usageFlag:

		err := callViewProfile()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case *createWordListFlag != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		domains := strings.Split(*createWordListFlag, ",")
		createWordList(domains, *workspaceFlag)
	case *addCustomWordsFlag != "":
		words := strings.Split(*addCustomWordsFlag, ",")
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		addCustomWordUser(words, *workspaceFlag)
	default:
		fmt.Println("No valid action specified.")
		flag.Usage()
		os.Exit(1)
	}
}

func extractRootWord(domain string) string {
	// Remove common TLDs and subdomains
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)

	// Remove protocol if present
	if strings.Contains(domain, "://") {
		parts := strings.Split(domain, "://")
		if len(parts) > 1 {
			domain = parts[1]
		}
	}

	// Split by dots and get the main domain part
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return domain
	}

	// Get the part before the TLD
	mainPart := parts[len(parts)-2]

	// Remove any non-alphanumeric characters
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	mainPart = reg.ReplaceAllString(mainPart, "")

	return mainPart
}
