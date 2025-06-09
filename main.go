package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type stringSliceFlag []string

var (
	silentFlag               bool
	uploadUrl                *string
	apiKeyFlag               *string
	updateFlag               *bool
	scanFileId               *string
	uploadFile               *string
	getAllResults            *string
	size                     *int
	workspaceFlag            *string
	listWorkspacesFlag       *bool
	getScannerResultsFlag    *bool
	query                    *string
	workspaceShort           *string
	workspaceLong            *string
	viewurls                 *bool
	scanDomainFlag           *string
	wordsFlag                *string
	urlswithmultipleResponse *bool
	getDomainsFlag           *bool
	headers                  stringSliceFlag
	addCustomWordsFlag       *string
	usageFlag                *bool
	viewfiles                *bool
	reverseSearchResults     *string
	createWordListFlag       *string
	searchUrlsByDomainFlag   *string
	getResultByJsmonId       *string
	getResultByFileId        *string
	totalAnalysisDataFlag    *bool
)

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

	body, err := io.ReadAll(resp.Body)
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
		fmt.Println("[INF] Use -wksp to list the workspaces")
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

func showAvailableWorkspaces() error {
	workspaces, err := getWorkspaces()
	if err != nil {
		fmt.Println(err)
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

func init() {
	// Remove the default -h / --help flag
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.CommandLine.Usage = func() {} // or keep it empty if you want no help at all

	// Define all flags
	flag.BoolVar(&silentFlag, "st", false, "Run in silent mode (no banner output)")
	uploadUrl = flag.String("u", "", "URL to upload for scanning")
	apiKeyFlag = flag.String("key", "", "API key for authentication")
	updateFlag = flag.Bool("ud", false, "Update jsmon-cli to the latest version")
	scanFileId = flag.String("fid", "", " File to be rescanned by fileId.")
	uploadFile = flag.String("f", "", "File to upload giving path to the file locally.")
	getAllResults = flag.String("jsi", "", "View JS Intelligence Data by domain name")
	size = flag.Int("s", 100, "Number of results to fetch (default 100)")
	workspaceFlag = flag.String("wksp", "", "Workspace ID")
	listWorkspacesFlag = flag.Bool("workspaces", false, "List all available workspaces")
	getScannerResultsFlag = flag.Bool("secrets", false, "View Keys & Secrets by domain name")
	query = flag.String("query", "", "Enable query builder functionality")
	workspaceShort = flag.String("cw", "", "Create a new workspace (Example: -cw nandini)")
	workspaceLong = flag.String("createWorkspace", "", "Create a new workspace (Example: -createWorkspace nandini)")
	viewurls = flag.Bool("urls", false, "view all urls")
	scanDomainFlag = flag.String("d", "", "Domain to automate scan")
	wordsFlag = flag.String("w", "", "Comma-separated list of words to include in the scan")
	urlswithmultipleResponse = flag.Bool("curls", false, "View changed JS URLs.")
	getDomainsFlag = flag.Bool("domains", false, "Get all domains for the user.")
	flag.Var(&headers, "H", "Custom headers in the format 'Key: Value' (can be used multiple times)")
	addCustomWordsFlag = flag.String("addCustomWords", "", "add custom words to the scan")
	usageFlag = flag.Bool("profile", false, "View user profile")
	viewfiles = flag.Bool("files", false, "view all files")
	reverseSearchResults = flag.String("rsearch", "", "Specify the input type (e.g., emails, domainname)")
	createWordListFlag = flag.String("wordlist", "", "creates a new word list from domains")
	searchUrlsByDomainFlag = flag.String("urlsByDomain", "", "Search URLs by domain")
	getResultByJsmonId = flag.String("jsiJsmonId", "", "Get JS Intelligence for the jsmon ID.")
	getResultByFileId = flag.String("jsiFileId", "", "Get JS Intelligence for the file ID.")
	totalAnalysisDataFlag = flag.Bool("count", false, "total count of overall analysis data")

	flag.Parse()
}

func main() {
	if !silentFlag {
		showBanner()
		displayVersion()
	}

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
	case *workspaceShort != "":
		createWorkspace(*workspaceShort)
	case *workspaceLong != "":
		createWorkspace(*workspaceLong)
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
				fmt.Printf("es: %v\n", err)
			}
			os.Exit(1)
		}
		err := uploadUrlEndpoint(*uploadUrl, headers, *workspaceFlag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
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
	case *query != "":
		if *workspaceFlag == "" {
			fmt.Println("No workspace specified. Use -workspaces to list available workspaces and provide a workspace ID using the -wksp flag.")
			err := displayWorkspaces()
			if err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
			}
			os.Exit(1)
		}
		// constructedQuery := fmt.Sprintf("field = %s, sub = %v, domain = %s", *field, *sub, *domain)
		queryBuilder(*workspaceFlag, *query)
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
		// fmt.Printf("Domain: %s, Words: %v\n", *scanDomainFlag, words)

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
