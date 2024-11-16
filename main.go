package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
func main() {
	scanUrl := flag.String("scanUrl", "", "URL to be rescanned by jsmonId.")
	uploadUrl := flag.String("uploadUrl", "", "URL to upload for scanning")
	apiKeyFlag := flag.String("apikey", "", "API key for authentication")
	scanFileId := flag.String("scanFile", "", " File to be rescanned by fileId.")
	uploadFile := flag.String("uploadFile", "", "File to upload giving path to the file locally.")
	getAllResults := flag.String("getAutomationData", "", "Get all automation results")
	size := flag.Int("size", 10000, "Number of results to fetch (default 10000)")
	fileTypes := flag.String("fileTypes", "", "files type (e.g. pdf,txt)")
	getScannerResultsFlag := flag.Bool("getScannerData", false, "Get scanner results")
	cron := flag.String("cron", "", "Set cronjob.")
	cronNotification := flag.String("notifications", "", "Set cronjob notification channel.")
	cronTime := flag.Int64("time", 0, "Set cronjob time.")
	cronType := flag.String("vulnerabilitiesType", "", "Set type[URLs, Analysis, Scanner] of cronjob.")
	cronDomains := flag.String("domains", "", "Set domains for cronjob.")
	cronDomainsNotify := flag.String("domainsNotify", "", "Set notify(true/false) for each domain for cronjob.")
	viewurls := flag.Bool("urls", false, "view all urls")
	viewurlsSize := flag.Int("urlSize", 10, "Number of URLs to fetch")
	scanDomainFlag := flag.String("scanDomain", "", "Domain to automate scan")
	wordsFlag := flag.String("words", "", "Comma-separated list of words to include in the scan")
	urlswithmultipleResponse := flag.Bool("changedUrls", false, "View changed JS URLs.")
	getDomainsFlag := flag.Bool("getDomains", false, "Get all domains for the user.")
	var headers stringSliceFlag
	flag.Var(&headers, "H", "Custom headers in the format 'Key: Value' (can be used multiple times)")
	addCustomWordsFlag := flag.String("addCustomWords", "", "add custom words to the scan")
	usageFlag := flag.Bool("usage", false, "View user profile")
	viewfiles := flag.Bool("getFiles", false, "view all files")
	viewEmails := flag.String("getEmails", "", "Get all emails for specified domains.")
	s3domains := flag.String("getS3Domains", "", "get all S3Domains for specified domains")
	ips := flag.String("getIps", "", "Get all IPs for specified domains")
	gql := flag.String("getGqlOps", "", "Get graph QL operations")
	domainUrl := flag.String("getDomainUrls", "", "Get Domain URLs for specified domains")
	apiPath := flag.String("getApiPaths", "", "Get the APIs for specified domains")
	fileExtensionUrls := flag.String("getFileExtensionUrls", "", "Get URLs containing any file type.")
	socialMediaUrls := flag.String("getSocialMediaUrls", "", "Get URLs for social media sites.")
	domainStatus := flag.String("getDomainStatus", "", "Get the availabilty of domains")
	queryParamsUrls := flag.String("getQueryParamsUrls", "", "Get URLs containing query params for specified domain.")
	localhostUrls := flag.String("getLocalhostUrls", "", "Get URLs with localhost in the hostname.")
	filteredPortUrls := flag.String("getUrlsWithPorts", "", "Get URLs with port numbers in the hostname")
	s3DomainsInvalid := flag.String("getS3DomainsInvalid", "", "Get available S3 domains (404 status).")
	compareFlag := flag.String("compare", "", "Compare two js responses by jsmon_ids (format: JSMON_ID1,JSMON_ID2)")
	reverseSearchResults := flag.String("reverseSearchResults", "", "Specify the input type (e.g., emails, domainname)")
	//getResultByValue := flag.String("value", "", "Specify the input value")
	createWordListFlag := flag.String("createWordList", "", "creates a new word list from domains")
	searchUrlsByDomainFlag := flag.String("searchUrlsByDomain", "", "Search URLs by domain")
	getResultByJsmonId := flag.String("getResultByJsmonId", "", "Get automation results by jsmon ID.")
	getResultByFileId := flag.String("getResultByFileId", "", "Get automation results by file ID.")
	rescanDomainFlag := flag.String("rescanDomain", "", "Rescan all URLs for a specific domain")
	totalAnalysisDataFlag := flag.Bool("totalAnalysisData", false, "total count of overall analysis data")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  %s [flags]\n\n", os.Args[0])
		fmt.Println("Flags:")

		// Input section
		fmt.Fprintf(os.Stderr, "\nINPUT:\n")
		fmt.Fprintf(os.Stderr, "  -scanUrl <jsmonId>             	  URL to be rescanned by jsmonId.\n")
		fmt.Fprintf(os.Stderr, "  -uploadUrl <URL>          		  URL to upload for scanning.\n")
		fmt.Fprintf(os.Stderr, "  -scanFile <fileId>         		  File to be rescanned by fileId.\n")
		fmt.Fprintf(os.Stderr, "  -uploadFile <local file path>       File to upload (local path)\n")
		fmt.Fprintf(os.Stderr, "  -scanDomain <domainName>   		  Domain to automate scan\n")

		// Authentication section
		fmt.Fprintf(os.Stderr, "\nAUTHENTICATION:\n")
		fmt.Fprintf(os.Stderr, "  -apikey <XXXXXX-XXXX-XXXX-XXXX-XXXXXX>          API key for authentication\n")

		// Output section
		fmt.Fprintf(os.Stderr, "\nOUTPUT:\n")
		fmt.Fprintf(os.Stderr, "  -getAutomationData <domainName>  Get all automation results.\n")
		fmt.Fprintf(os.Stderr, "  -getScannerData                  Get scanner results.\n")
		fmt.Fprintf(os.Stderr, "  -getUrls                         View all URLs.\n")
		fmt.Fprintf(os.Stderr, "  -urlSize int                     Number of URLs to fetch (default 10).\n")
		fmt.Fprintf(os.Stderr, "  -getFiles                        View all files.\n")
		fmt.Fprintf(os.Stderr, "  -fileTypes <types>               Specify file types (e.g., pdf,txt), use ',' as separator.\n")
		fmt.Fprintf(os.Stderr, "  -usage                           View user profile.\n")
		fmt.Fprintf(os.Stderr, "  -changedUrls                     View changed JS URLs.\n")

		// Cronjob section
		// fmt.Fprintf(os.Stderr, "\nCRON JOB:\n")
		// fmt.Fprintf(os.Stderr, "  -cron <string>                   Set, update, or stop cronjob\n")
		// fmt.Fprintf(os.Stderr, "  -notifications <string>          Set cronjob notification channel\n")
		// fmt.Fprintf(os.Stderr, "  -time <int>                      Set cronjob time\n")
		// fmt.Fprintf(os.Stderr, "  -vulnerabilitiesType <type>      Set type of cronjob (URLs, Analysis, Scanner)\n")
		// fmt.Fprintf(os.Stderr, "  -domains <string>                Set domains for cronjob\n")
		// fmt.Fprintf(os.Stderr, "  -domainsNotify <true/false>      Set notification for each domain\n")

		// Additional options
		fmt.Fprintf(os.Stderr, "\nADDITIONAL OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  -H <Key: Value>                  Custom headers (can be used multiple times).\n")
		fmt.Fprintf(os.Stderr, "  -words <words>                   Comma-separated list of words to include in the scan.\n")
		fmt.Fprintf(os.Stderr, "  -getDomains                      Get all domains for the user.\n")
		fmt.Fprintf(os.Stderr, "  -getEmails <domain>              Get all emails for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -getS3Domains <domain>           Get all S3 domains for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -getIps <domain>                 Get all IPs for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -getDomainUrls <domain>          Get domain URLs for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -getApiPaths <domain>            Get API paths for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -getFileExtensionUrls <domain>   Get URLs containing any file type.\n")
		fmt.Fprintf(os.Stderr, "  -getSocialMediaUrls <domain>     Get URLs for social media sites.\n")
		fmt.Fprintf(os.Stderr, "  -getDomainStatus <domain>        Get availability status of domains.\n")
		fmt.Fprintf(os.Stderr, "  -getQueryParamsUrls <domain>     Get URLs containing query params for specified domain.\n")
		fmt.Fprintf(os.Stderr, "  -getLocalhostUrls <domain>       Get URLs with localhost in the hostname.\n")
		fmt.Fprintf(os.Stderr, "  -getUrlsWithPorts <domain>       Get URLs with port numbers in the hostname.\n")
		fmt.Fprintf(os.Stderr, "  -getS3DomainsInvalid <domain>    Get available S3 domains (404 status).\n")
		fmt.Fprintf(os.Stderr, "  -rescanDomain <domain>           Rescan all URLs for a specific domain.\n")
		fmt.Fprintf(os.Stderr, "  -searchUrlsByDomain <domain>     Search URLs by domain.\n")
		fmt.Fprintf(os.Stderr, "  -compare <ID1,ID2>               Compare two JS responses by IDs (format: ID1,ID2).\n")
		fmt.Fprintf(os.Stderr, "  -getGqlOps <domain>              Get GraphQL operations for specified domains.\n")
		fmt.Fprintf(os.Stderr, "  -totalAnalysisData               Get total count of overall analysis data.\n")
		fmt.Fprintf(os.Stderr, "  -getResultByJsmonId <ID>         Get automation results by jsmon ID.\n")
		fmt.Fprintf(os.Stderr, "  -getResultByFileId <ID>          Get automation results by file ID.\n")

		// Automation results section
		fmt.Fprintf(os.Stderr, "\nAUTOMATION RESULTS BY FIELD:\n")
		fmt.Fprintf(os.Stderr, "  -reverseSearchResults <field>=<value>\n")
		fmt.Fprintf(os.Stderr, "    Search by field : emails, domainname, extracteddomains, s3domains, url, extractedurls, ipv4addresses, ipv6addresses, jwttokens, gqlquery, gqlmutation, guids, apipaths, vulnerabilities, nodemodules, domainstatus, queryparamsurls, socialmediaurls, filterdporturls, gqlfragment, s3domainsinvalid, fileextensionurls, localhosturls.\n")

	}

	flag.Parse()

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

	switch {
	case *scanFileId != "":
		scanFileEndpoint(*scanFileId)
	case *uploadFile != "":
		uploadFileEndpoint(*uploadFile, headers)
	case *viewurls:
		viewUrls(*viewurlsSize)
	case *viewfiles:
		viewFiles()
	case *uploadUrl != "":
		uploadUrlEndpoint(*uploadUrl, headers)
	case *rescanDomainFlag != "":
		rescanDomain(*rescanDomainFlag)
	case *totalAnalysisDataFlag:
		totalAnalysisData()
	case *searchUrlsByDomainFlag != "":
		searchUrlsByDomain(*searchUrlsByDomainFlag)
	case *urlswithmultipleResponse:
		urlsmultipleResponse()
	case *viewEmails != "":
		domains := strings.Split(*viewEmails, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		getEmails(domains)
	case *getResultByJsmonId != "":
		getAutomationResultsByJsmonId(strings.TrimSpace(*getResultByJsmonId))
	case *reverseSearchResults != "":
		parts := strings.SplitN(*reverseSearchResults, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid format for reverseSearchResults. Use field=value format.")
			return
		}

		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		getAutomationResultsByInput(field, value)

	case *getResultByFileId != "":
		getAutomationResultsByFileId(strings.TrimSpace(*getResultByFileId))
	case *s3domains != "":
		domains := strings.Split(*s3domains, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		getS3Domains(domains)
	case *ips != "":
		domains := strings.Split(*ips, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		getAllIps(domains)
	case *gql != "":
		domains := strings.Split(*gql, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		getGqlOps(domains)
	case *domainUrl != "":
		domains := strings.Split(*domainUrl, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		getDomainUrls(domains)
	case *apiPath != "":
		domains := strings.Split(*apiPath, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		getApiPaths(domains)
	case *getScannerResultsFlag:
		getScannerResults()
	case *scanUrl != "":
		rescanUrlEndpoint(*scanUrl)
	case *cron == "start":
		StartCron(*cronNotification, *cronTime, *cronType, *cronDomains, *cronDomainsNotify)
	case *cron == "stop":
		StopCron()
	case *getDomainsFlag:
		getDomains()

	case *fileExtensionUrls != "":
		extensions := strings.Split(*fileTypes, ",")
		for i, extension := range extensions {
			extensions[i] = strings.TrimSpace(extension)
		}
		getAllFileExtensionUrls(*fileExtensionUrls, extensions, *size)
	case *domainStatus != "":
		// domains := strings.Split(*domainStatus, ",")
		// for i, domain := range domains {
		// 	domains[i] = strings.TrimSpace(domain)
		// }
		getAllDomainsStatus(*domainStatus, *size)

	case *socialMediaUrls != "":
		getAllSocialMediaUrls(*socialMediaUrls, *size)
	case *queryParamsUrls != "":
		getAllQueryParamsUrls(*queryParamsUrls, *size)
	case *localhostUrls != "":
		getAllLocalhostUrls(*localhostUrls, *size)
	case *filteredPortUrls != "":
		getAllFilteredPortUrls(*filteredPortUrls, *size)
	case *s3DomainsInvalid != "":
		getAllS3DomainsInvalid(*s3DomainsInvalid, *size)
	case *compareFlag != "":
		ids := strings.Split(*compareFlag, ",")
		if len(ids) != 2 {
			fmt.Println("Invalid format for compare. Use: JSMON_ID1,JSMON_ID2")
			os.Exit(1)
		}
		compareEndpoint(strings.TrimSpace(ids[0]), strings.TrimSpace(ids[1]))
	case *cron == "update":
		UpdateCron(*cronNotification, *cronType, *cronDomains, *cronDomainsNotify, *cronTime)
	case *getAllResults != "":
		getAllAutomationResults(*getAllResults, *size)
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
		fmt.Printf("Domain: %s, Words: %v\n", *scanDomainFlag, words)
		automateScanDomain(*scanDomainFlag, words)

	case *usageFlag:
		callViewProfile()
	case *createWordListFlag != "":
		domains := strings.Split(*createWordListFlag, ",")
		createWordList(domains)
	case *addCustomWordsFlag != "":
		words := strings.Split(*addCustomWordsFlag, ",")
		addCustomWordUser(words)
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

// type Args struct {
// 	Cron             string
// 	CronNotification string
// 	CronTime         int64
// 	CronType         string
// }

// func parseArgs() Args {
// 	//CRON JOB FLAGS ->
// 	cron := flag.String("cron", "", "Set cronjob.")
// 	cronNotification := flag.String("notifications", "", "Set cronjob notification.")
// 	cronTime := flag.Int64("time", 0, "Set cronjob time.")
// 	cronType := flag.String("type", "", "Set type of cronjob.")

// 	flag.Parse()

// 	return Args{
// 		Cron:             *cron,
// 		CronNotification: *cronNotification,
// 		CronTime:         *cronTime,
// 		CronType:         *cronType,
// 	}
// }
