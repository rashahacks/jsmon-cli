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

type Flags struct {
	ScanUrl              string
	UploadUrl            string
	ApiKey               string
	ScanFileId           string
	UploadFile           string
	GetAllResults        string
	Size                 int
	FileTypes            string
	GetScannerResults    bool
	Cron                 string
	CronNotification     string
	CronTime             int64
	CronType             string
	CronDomains          string
	CronDomainsNotify    string
	ViewUrls             bool
	ViewUrlsSize         int
	ScanDomain           string
	Words                string
	UrlsWithMultipleResp bool
	GetDomains           bool
	Headers              stringSliceFlag
	AddCustomWords       string
	Usage                bool
	ViewFiles            bool
	ViewEmails           string
	S3Domains            string
	Ips                  string
	Gql                  string
	DomainUrl            string
	ApiPath              string
	FileExtensionUrls    string
	SocialMediaUrls      string
	DomainStatus         string
	QueryParamsUrls      string
	LocalhostUrls        string
	FilteredPortUrls     string
	S3DomainsInvalid     string
	Compare              string
	ReverseSearchResults string
	CreateWordList       string
	SearchUrlsByDomain   string
	GetResultByJsmonId   string
	GetResultByFileId    string
	RescanDomain         string
	TotalAnalysisData    bool
}

func parseFlags() *Flags {
	f := &Flags{}

	flag.StringVar(&f.ScanUrl, "scanUrl", "", "URL to be rescanned by jsmonId.")
	flag.StringVar(&f.UploadUrl, "uploadUrl", "", "URL to upload for scanning")
	flag.StringVar(&f.ApiKey, "apikey", "", "API key for authentication")
	flag.StringVar(&f.ScanFileId, "scanFile", "", "File to be rescanned by fileId.")
	flag.StringVar(&f.UploadFile, "uploadFile", "", "File to upload giving path to the file locally.")
	flag.StringVar(&f.GetAllResults, "getAutomationData", "", "Get all automation results")
	flag.IntVar(&f.Size, "size", 10000, "Number of results to fetch (default 10000)")
	flag.StringVar(&f.FileTypes, "fileTypes", "", "files type (e.g. pdf,txt)")
	flag.BoolVar(&f.GetScannerResults, "getScannerData", false, "Get scanner results")
	flag.StringVar(&f.Cron, "cron", "", "Set cronjob.")
	flag.StringVar(&f.CronNotification, "notifications", "", "Set cronjob notification channel.")
	flag.Int64Var(&f.CronTime, "time", 0, "Set cronjob time.")
	flag.StringVar(&f.CronType, "vulnerabilitiesType", "", "Set type[URLs, Analysis, Scanner] of cronjob.")
	flag.StringVar(&f.CronDomains, "domains", "", "Set domains for cronjob.")
	flag.StringVar(&f.CronDomainsNotify, "domainsNotify", "", "Set notify(true/false) for each domain for cronjob.")
	flag.BoolVar(&f.ViewUrls, "urls", false, "view all urls")
	flag.IntVar(&f.ViewUrlsSize, "urlSize", 10, "Number of URLs to fetch")
	flag.StringVar(&f.ScanDomain, "scanDomain", "", "Domain to automate scan")
	flag.StringVar(&f.Words, "words", "", "Comma-separated list of words to include in the scan")
	flag.BoolVar(&f.UrlsWithMultipleResp, "changedUrls", false, "View changed JS URLs.")
	flag.BoolVar(&f.GetDomains, "getDomains", false, "Get all domains for the user.")
	flag.Var(&f.Headers, "H", "Custom headers in the format 'Key: Value' (can be used multiple times)")
	flag.StringVar(&f.AddCustomWords, "addCustomWords", "", "add custom words to the scan")
	flag.BoolVar(&f.Usage, "usage", false, "View user profile")
	flag.BoolVar(&f.ViewFiles, "getFiles", false, "view all files")
	flag.StringVar(&f.ViewEmails, "getEmails", "", "Get all emails for specified domains.")
	flag.StringVar(&f.S3Domains, "getS3Domains", "", "get all S3Domains for specified domains")
	flag.StringVar(&f.Ips, "getIps", "", "Get all IPs for specified domains")
	flag.StringVar(&f.Gql, "getGqlOps", "", "Get graph QL operations")
	flag.StringVar(&f.DomainUrl, "getDomainUrls", "", "Get Domain URLs for specified domains")
	flag.StringVar(&f.ApiPath, "getApiPaths", "", "Get the APIs for specified domains")
	flag.StringVar(&f.FileExtensionUrls, "getFileExtensionUrls", "", "Get URLs containing any file type.")
	flag.StringVar(&f.SocialMediaUrls, "getSocialMediaUrls", "", "Get URLs for social media sites.")
	flag.StringVar(&f.DomainStatus, "getDomainStatus", "", "Get the availabilty of domains")
	flag.StringVar(&f.QueryParamsUrls, "getQueryParamsUrls", "", "Get URLs containing query params for specified domain.")
	flag.StringVar(&f.LocalhostUrls, "getLocalhostUrls", "", "Get URLs with localhost in the hostname.")
	flag.StringVar(&f.FilteredPortUrls, "getUrlsWithPorts", "", "Get URLs with port numbers in the hostname")
	flag.StringVar(&f.S3DomainsInvalid, "getS3DomainsInvalid", "", "Get available S3 domains (404 status).")
	flag.StringVar(&f.Compare, "compare", "", "Compare two js responses by jsmon_ids (format: JSMON_ID1,JSMON_ID2)")
	flag.StringVar(&f.ReverseSearchResults, "reverseSearchResults", "", "Specify the input type (e.g., emails, domainname)")
	flag.StringVar(&f.CreateWordList, "createWordList", "", "creates a new word list from domains")
	flag.StringVar(&f.SearchUrlsByDomain, "searchUrlsByDomain", "", "Search URLs by domain")
	flag.StringVar(&f.GetResultByJsmonId, "getResultByJsmonId", "", "Get automation results by jsmon ID.")
	flag.StringVar(&f.GetResultByFileId, "getResultByFileId", "", "Get automation results by file ID.")
	flag.StringVar(&f.RescanDomain, "rescanDomain", "", "Rescan all URLs for a specific domain")
	flag.BoolVar(&f.TotalAnalysisData, "totalAnalysisData", false, "total count of overall analysis data")

	flag.Parse()

	return f
}

func main() {
	flags := parseFlags()

	if flags.ApiKey != "" {
		setAPIKey(flags.ApiKey)
	} else {
		err := loadAPIKey()
		if err != nil {
			fmt.Println("Error loading API key:", err)
			fmt.Println("Please provide an API key using the -apikey flag.")
			os.Exit(1)
		}
	}

	if flag.NFlag() == 0 || (flag.NFlag() == 1 && flags.ApiKey != "") {
		fmt.Println("No action specified. Use -h or --help for usage information.")
		flag.Usage()
		os.Exit(1)
	}

	err := executeCommand(flags)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func executeCommand(flags *Flags) error {
	switch {
	case flags.ScanFileId != "":
		return scanFileEndpoint(flags.ScanFileId)
	case flags.UploadFile != "":
		return uploadFileEndpoint(flags.UploadFile, flags.Headers)
	case flags.ViewUrls:
		return viewUrls(flags.ViewUrlsSize)
	case flags.ViewFiles:
		return viewFiles()
	case flags.UploadUrl != "":
		return uploadUrlEndpoint(flags.UploadUrl, flags.Headers)
	case flags.RescanDomain != "":
		return rescanDomain(flags.RescanDomain)
	case flags.TotalAnalysisData:
		return totalAnalysisData()
	case flags.SearchUrlsByDomain != "":
		return searchUrlsByDomain(flags.SearchUrlsByDomain)
	case flags.UrlsWithMultipleResp:
		return urlsmultipleResponse()
	case flags.ViewEmails != "":
		domains := splitAndTrim(flags.ViewEmails)
		return getEmails(domains)
	case flags.GetResultByJsmonId != "":
		return getAutomationResultsByJsmonId(strings.TrimSpace(flags.GetResultByJsmonId))
	case flags.ReverseSearchResults != "":
		return handleReverseSearchResults(flags.ReverseSearchResults)
	case flags.GetResultByFileId != "":
		return getAutomationResultsByFileId(strings.TrimSpace(flags.GetResultByFileId))
	case flags.S3Domains != "":
		domains := splitAndTrim(flags.S3Domains)
		return getS3Domains(domains)
	case flags.Ips != "":
		domains := splitAndTrim(flags.Ips)
		return getAllIps(domains)
	case flags.Gql != "":
		domains := splitAndTrim(flags.Gql)
		return getGqlOps(domains)
	case flags.DomainUrl != "":
		domains := splitAndTrim(flags.DomainUrl)
		return getDomainUrls(domains)
	case flags.ApiPath != "":
		domains := splitAndTrim(flags.ApiPath)
		return getApiPaths(domains)
	case flags.GetScannerResults:
		return getScannerResults()
	case flags.ScanUrl != "":
		return rescanUrlEndpoint(flags.ScanUrl)
	case flags.Cron == "start":
		return StartCron(flags.CronNotification, flags.CronTime, flags.CronType, flags.CronDomains, flags.CronDomainsNotify)
	case flags.Cron == "stop":
		return StopCron()
	case flags.GetDomains:
		return getDomains()
	case flags.FileExtensionUrls != "":
		extensions := splitAndTrim(flags.FileTypes)
		return getAllFileExtensionUrls(flags.FileExtensionUrls, extensions, flags.Size)
	case flags.DomainStatus != "":
		return getAllDomainsStatus(flags.DomainStatus, flags.Size)
	case flags.SocialMediaUrls != "":
		return getAllSocialMediaUrls(flags.SocialMediaUrls, flags.Size)
	case flags.QueryParamsUrls != "":
		return getAllQueryParamsUrls(flags.QueryParamsUrls, flags.Size)
	case flags.LocalhostUrls != "":
		return getAllLocalhostUrls(flags.LocalhostUrls, flags.Size)
	case flags.FilteredPortUrls != "":
		return getAllFilteredPortUrls(flags.FilteredPortUrls, flags.Size)
	case flags.S3DomainsInvalid != "":
		return getAllS3DomainsInvalid(flags.S3DomainsInvalid, flags.Size)
	case flags.Compare != "":
		return handleCompare(flags.Compare)
	case flags.Cron == "update":
		return UpdateCron(flags.CronNotification, flags.CronType, flags.CronDomains, flags.CronDomainsNotify, flags.CronTime)
	case flags.GetAllResults != "":
		return getAllAutomationResults(flags.GetAllResults, flags.Size)
	case flags.ScanDomain != "":
		words := splitAndTrim(flags.Words)
		fmt.Printf("Domain: %s, Words: %v\n", flags.ScanDomain, words)
		return automateScanDomain(flags.ScanDomain, words)
	case flags.Usage:
		return callViewProfile()
	case flags.CreateWordList != "":
		domains := splitAndTrim(flags.CreateWordList)
		return createWordList(domains)
	case flags.AddCustomWords != "":
		words := splitAndTrim(flags.AddCustomWords)
		return addCustomWordUser(words)
	default:
		return fmt.Errorf("no valid action specified")
	}
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

func handleReverseSearchResults(input string) error {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format for reverseSearchResults. Use field=value format")
	}

	field := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	return getAutomationResultsByInput(field, value)
}

func handleCompare(input string) error {
	ids := strings.Split(input, ",")
	if len(ids) != 2 {
		return fmt.Errorf("invalid format for compare. Use: JSMON_ID1,JSMON_ID2")
	}
	return compareEndpoint(strings.TrimSpace(ids[0]), strings.TrimSpace(ids[1]))
}

// ... (rest of the code remains the same)
