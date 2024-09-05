package main

import (
	"flag"
	"fmt"
	"os"
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
	scanUrl := flag.String("scanUrl", "", "URL or scan ID to rescan")
	uploadUrl := flag.String("uploadUrl", "", "URL to upload for scanning")
	apiKeyFlag := flag.String("apikey", "", "API key for authentication")
	scanFileId := flag.String("scanFile", "", "File ID to scan")
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
	urlswithmultipleResponse := flag.Bool("urlswithmultipleresponse", false, "check for urls with multiple responses.")
	getDomainsFlag := flag.Bool("getDomains", false, "Get all domains for the user")
	var headers stringSliceFlag
	flag.Var(&headers, "H", "Custom headers in the format 'Key: Value' (can be used multiple times)")

	usageFlag := flag.Bool("usage", false, "View user profile")
	viewfiles := flag.Bool("getFiles", false, "view all files")
	viewEmails := flag.String("getEmails", "", "view all Emails for specified domains")
	s3domains := flag.String("getS3Domains", "", "get all S3Domains for specified domains")
	ips := flag.String("getIps", "", "get all IPs for specified domains")
	gql := flag.String("getGqlOps", "", "get graph QL operations")
	domainUrl := flag.String("getDomainUrls", "", "get Domain URLs for specified domains")
	apiPath := flag.String("getApiPaths", "", "get the APIs for specified domains")
	fileExtensionUrls := flag.String("getFileExtensionUrls", "", "get the urls containing any type of file")
	socialMediaUrls := flag.String("getSocialMediaUrls", "", "get the urls for the social media sites")
	domainStatus := flag.String("getDomainStatus", "" , "get the availabilty of domains")
	queryParamsUrls := flag.String("getQueryParamsUrls", "", "get the urls containing query params for the specified domain")
	localhostUrls := flag.String("getLocalhostUrls", "", "get the urls which has localhost in the hostname for the specified domain")
	filteredPortUrls := flag.String("getUrlsWithPort", "", "get the urls containing a port number in the hostname for the specified domain")
	s3DomainsInvalid := flag.String("getS3DomainsInvalid", "", "get the S3 domains which are available (404 status code) for the specified domain")
	compareFlag := flag.String("compare", "", "Compare two js responses by jsmon_ids (format: JSMON_ID1,JSMON_ID2)")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")

		fmt.Fprintf(os.Stderr, "INPUT:\n")
		fmt.Fprintf(os.Stderr, "  -scanUrl string         URL or scan ID to rescan\n")
		fmt.Fprintf(os.Stderr, "  -uploadUrl string       URL to upload for scanning\n")
		fmt.Fprintf(os.Stderr, "  -scanFile string        File ID to scan\n")
		fmt.Fprintf(os.Stderr, "  -uploadFile string      File to upload (local path)\n")
		fmt.Fprintf(os.Stderr, "  -scanDomain string      Domain to automate scan\n")

		fmt.Fprintf(os.Stderr, "\nAUTHENTICATION:\n")
		fmt.Fprintf(os.Stderr, "  -apikey string          API key for authentication\n")

		fmt.Fprintf(os.Stderr, "\nOUTPUT:\n")
		fmt.Fprintf(os.Stderr, "  -getAutomationData string  Get all automation results\n")
		fmt.Fprintf(os.Stderr, "  -getScannerData            Get scanner results\n")
		fmt.Fprintf(os.Stderr, "  -getUrls                   View all URLs\n")
		fmt.Fprintf(os.Stderr, "  -urlSize int               Number of URLs to fetch (default 10)\n")
		fmt.Fprintf(os.Stderr, "  -getFiles                  View all files\n")
		fmt.Fprintf(os.Stderr, "  -usage                  View user profile\n")
		fmt.Fprintf(os.Stderr, "  -urlswithmultipleresponse  View user profile\n")

		fmt.Fprintf(os.Stderr, "\nCRON JOB:\n")
		fmt.Fprintf(os.Stderr, "  -cron string            Set, update, or stop cronjob\n")
		fmt.Fprintf(os.Stderr, "  -notifications string   Set cronjob notification channel\n")
		fmt.Fprintf(os.Stderr, "  -time int               Set cronjob time\n")
		fmt.Fprintf(os.Stderr, "  -vulnerabilitiesType    Set type of cronjob (URLs, Analysis, Scanner)\n")
		fmt.Fprintf(os.Stderr, "  -domains string         Set domains for cronjob\n")
		fmt.Fprintf(os.Stderr, "  -domainsNotify string   Set notify (true/false) for each domain\n")

		fmt.Fprintf(os.Stderr, "\nADDITIONAL OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  -H string               Custom headers (Key: Value, can be used multiple times)\n")
		fmt.Fprintf(os.Stderr, "  -words string           Comma-separated list of words to include in the scan\n")
		fmt.Fprintf(os.Stderr, "  -getDomains             Get all domains for the user\n")
		fmt.Fprintf(os.Stderr, "  -getEmails string          View all Emails for specified domains\n")
		fmt.Fprintf(os.Stderr, "  -getS3Domains string       Get all S3 Domains for specified domains\n")
		fmt.Fprintf(os.Stderr, "  -getIps string             Get all IPs for specified domains\n")
		fmt.Fprintf(os.Stderr, "  -getDomainUrls string      Get Domain URLs for specified domains\n")
		fmt.Fprintf(os.Stderr, "  -getApiPaths string             	Get the APIs for specified domains\n")
		fmt.Fprintf(os.Stderr, "  -getFileExtensionUrls string     	Get the urls containing any type of file\n")
		fmt.Fprintf(os.Stderr, "  -getSocialMediaUrls string       	Get the urls for the social media sites\n")
		fmt.Fprintf(os.Stderr, "  -getDomainStatus       			Get the availabilty of domains\n")
		fmt.Fprintf(os.Stderr, "  -getQueryParamsUrls       		Get the urls containing query params\n")
		fmt.Fprintf(os.Stderr, "  -getLocalhostUrls       			Get the urls containing localhost in their urls (includes local ip address)\n")
		fmt.Fprintf(os.Stderr, "  -getUrlsWithPorts       			Get the urls containing port number in their hostname\n")
		fmt.Fprintf(os.Stderr, "  -getS3DomainsInvalid       		Get the s3 bucket urls which are available (having 404 status code)\n")

		fmt.Fprintf(os.Stderr, "  -compare string         Compare two JS responses by JSMON_IDs (format: ID1,ID2)\n")
	}

	// Handle API key
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

	// Check if any action flags were provided
	if flag.NFlag() == 0 || (flag.NFlag() == 1 && *apiKeyFlag != "") {
		fmt.Println("No action specified. Use -h or --help for usage information.")
		flag.Usage()
		os.Exit(1)
	}

	// Execute the appropriate function based on the provided flag
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
	case *urlswithmultipleResponse:
		urlsmultipleResponse()
	case *viewEmails != "":
		domains := strings.Split(*viewEmails, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		getEmails(domains)
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
		}
		fmt.Printf("Domain: %s, Words: %v\n", *scanDomainFlag, words) // Add this line for debugging
		automateScanDomain(*scanDomainFlag, words)
	case *usageFlag:
		callViewProfile()
	default:
		fmt.Println("No valid action specified.")
		flag.Usage()
		os.Exit(1)
	}
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
