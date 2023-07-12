package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	httpxUtilz "httpxUtilz/utilz"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

type PassiveResult struct {
	CName       string   `json:"cname"`
	IP          string   `json:"ip"`
	Cdn         int      `json:"cdn"`
	CdnByIP     bool     `json:"cdn_by_ip"`
	CdnByHeader []string `json:"cdn_by_header"`
	CdnByCidr   bool     `json:"cdn_by_cidr"`
	CdnByAsn    bool     `json:"cdn_by_asn"`
	CdnByCName  bool     `json:"cdn_by_cname"`
	Cidr        string   `json:"cidr"`
	Asn         string   `json:"asn"`
	Org         string   `json:"org"`
	Addr        string   `json:"addr"`
}

type ResponseResult struct {
	Url                    string   `json:"url"`
	Title                  string   `json:"title"`
	Server                 string   `json:"server"`
	Via                    string   `json:"via"`
	Power                  string   `json:"x-powered-by"`
	StatusCode             int      `json:"status_code"`
	Alive                  int      `json:"alive"`
	ContentLength          int64    `json:"content_length"`
	ContentLengthByAllBody int64    `json:"content_length_by_all_body"`
	ResponseHeader         []string `json:"response_header"`
}

type MatchResponseResult struct {
	MayVul map[string]string `json:"may_vul"`
}

type Result struct {
	BaseInfo    ResponseResult      `json:"base_info"`
	PassiveInfo PassiveResult       `json:"passive_info"`
	RegexInfo   MatchResponseResult `json:"regex_info"`
}

type ProcessUrlParams struct {
	URLPipe         []string
	Url             string
	Filename        string
	Proxy           string
	UseHTTPS        bool
	FollowRedirects bool
	MaxRedirects    int
	Method          string
	RandomUserAgent bool
	Headers         string
	FollowSameHost  bool
	Timeout         int
	Processes       int
	RateLimit       int
	Res             bool
	ResultFile      string
	Passive         bool
	Base            bool
	MayVul          bool
}

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("readURLsFromFile: Open File Error", err)
		return nil, err
	}
	defer file.Close()

	urls := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		urls = append(urls, url)
	}

	if err := scanner.Err(); err != nil {
		log.Println("readURLsFromFile: Read File Error", err)
		return nil, err
	}

	return urls, nil
}

func isResultEmpty(result Result) bool {
	emptyResult := Result{}

	// Compare two structs for equality using the DeepEqual function
	return reflect.DeepEqual(result, emptyResult)
}

//func saveResultsToFile(results []Result, resultFile string) {
//	// The default path for the result file is "./result.json"
//	if resultFile == "" {
//		resultFile = "./result.json"
//	}
//
//	// Create the result file
//	file, err := os.OpenFile(resultFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
//	if err != nil {
//		log.Println("saveResultsToFile> create result file error:", err)
//		return
//	}
//	defer file.Close()
//
//	// Convert the result to a JSON string.
//	jsonData, err := json.Marshal(results)
//	if err != nil {
//		log.Println("saveResultsToFile> json marshal error:", err)
//		return
//	}
//
//	// Write the JSON string to a file
//	_, err = file.Write(jsonData)
//	if err != nil {
//		log.Println("saveResultsToFile> write result to file error:", err)
//		return
//	}
//
//	log.Println("results saved to", resultFile)
//}

// writes the content of the buffer to the specified file
func WriteBufferToFile(buffer *bytes.Buffer, filePath string) error {
	// The default path for the result file is "./result.json"
	if filePath == "" {
		filePath = "./result.json"
	}
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("WriteBufferToFile> failed to create file: %w", err)
	}
	defer file.Close()

	_, err = buffer.WriteTo(file)
	return err
}

func processURL(params ProcessUrlParams) (result Result) {
	config := httpxUtilz.RequestClientConfig{
		ProxyURL:        params.Proxy,
		UseHTTPS:        params.UseHTTPS,
		FollowRedirects: params.FollowRedirects,
		MaxRedirects:    params.MaxRedirects,
		Method:          params.Method,
		RandomUserAgent: params.RandomUserAgent,
		Headers: map[string]string{
			"User-Agent": params.Headers,
		},
		FollowSameHost: params.FollowSameHost,
		Timeout:        time.Duration(params.Timeout),
	}

	var (
		title                  string
		server                 string
		via                    string
		power                  string
		statusCode             int
		alive                  int
		contentLength          int64
		contentLengthByAllBody int64
		responseHeader         []string
		resp                   *httpxUtilz.Response
		err                    error
	)

	if params.Base {
		resp, err = config.GetResponseByUrl(params.Url)
		if err != nil {
			log.Println("processURL>  request error: ", err)
			return
		}

		title = config.GetTitleByResponse(resp)
		server, via, power = config.GetBannerByResponse(resp)
		statusCode = config.GetStatusByResponse(resp)
		alive = config.GetAliveByResponse(resp)
		contentLength = config.GetContentLengthByResponse(resp)
		contentLengthByAllBody = config.GetContentLengthAllBodyByResponse(resp)
		responseHeader = config.GetServerAllHeaderByResponse(resp)
	}

	baseInfo := ResponseResult{
		Url:                    params.Url,
		Title:                  title,
		Server:                 server,
		Via:                    via,
		Power:                  power,
		StatusCode:             statusCode,
		Alive:                  alive,
		ContentLength:          contentLength,
		ContentLengthByAllBody: contentLengthByAllBody,
		ResponseHeader:         responseHeader,
	}

	var (
		cdn          int
		cdnbyip      bool
		cdnbyheader  []string
		cdnbycidr    bool
		cdnbyasn     bool
		cdnbycname   bool
		passiveInfos PassiveResult
	)
	if params.Passive {
		cname, cnameIps := config.GetCNameIPByDomain(params.Url, "./data/vaildResolvers.txt")
		resolveIps := config.GetIpsByAsnmap(params.Url)
		ips := cnameIps + resolveIps
		if len(ips) == 0 {
			return Result{}
		}
		cidr, asn, org, addr := config.GetAsnInfoByIp(ips, params.Proxy)

		if len(ips) > 0 {

			if !params.Base { // not get baseinfo, but cdnbyheader need response
				resp, err = config.GetResponseByUrl(params.Url)
				if err != nil {
					log.Println("processURL>  request error: ", err)
					return
				}
			}

			cdn, cdnbyip, cdnbyheader, cdnbycidr, cdnbyasn, cdnbycname = config.GetCdnInfoByAll(
				resp, ips, "./data/cdn_header_keys.json",
				cidr, "./data/cdn_ip_cidr.json",
				asn, "./data/cdn_asn_list.json",
				cname, "./data/cdn_cname_keywords.json")
		}

		passiveInfos = PassiveResult{
			CName:       cname,
			IP:          ips,
			Cdn:         cdn,
			CdnByIP:     cdnbyip,
			CdnByHeader: cdnbyheader,
			CdnByCidr:   cdnbycidr,
			CdnByAsn:    cdnbyasn,
			CdnByCName:  cdnbycname,
			Cidr:        cidr,
			Asn:         asn,
			Org:         org,
			Addr:        addr,
		}
	}

	var (
		matchResponseResult MatchResponseResult
	)

	if params.MayVul {
		if !params.Base { // not get baseinfo, but regex matches need response
			resp, err = config.GetResponseByUrl(params.Url)
			if err != nil {
				log.Println("processURL>  request error: ", err)
				return
			}
		}
		matchResponseResult.MayVul = config.GetMayVulInfoByRespone(resp, "./data/regex_MayVul.json")
	}

	result = Result{
		BaseInfo:    baseInfo,
		PassiveInfo: passiveInfos,
		RegexInfo:   matchResponseResult,
	}
	return
}

func ProcessURLFromLine(params ProcessUrlParams) {
	// Create a wait group to wait for all Goroutines to complete
	var wg sync.WaitGroup

	// Create a channel to limit the rate of requests
	rateLimiter := time.Tick(time.Second / time.Duration(params.RateLimit))

	// Create a list of results
	results := make([]Result, 0)

	// Create a buffer to store the results temporarily
	var buffer bytes.Buffer

	// Initiate multiple Goroutines for concurrent processing
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Retrieve a token from the channel to control the rate
			<-rateLimiter

			// Perform the request and processing
			result := processURL(params)

			if !isResultEmpty(result) {
				// Add the result to the list.
				results = append(results, result)

				jsonData, err := json.Marshal(result)
				if err != nil {
					log.Println("ProcessURLFromLine> json marshal error:", err)
					return
				}

				fmt.Println(string(jsonData))

				buffer.WriteString(string(jsonData))
				buffer.WriteString("\n")
			} else {
				log.Println(params.Url + " can't get result")
			}

		}()
	}

	// Wait for all Goroutines to complete
	wg.Wait()

	// Save the results to a JSON file
	if params.Res && buffer.Len() > 0 {
		//saveResultsToFile(results, params.ResultFile)
		err := WriteBufferToFile(&buffer, params.ResultFile)
		if err != nil {
			fmt.Println("WriteBufferToFile Error:", err)
			return
		}
	}
}

func ProcessURLFromGroup(params ProcessUrlParams) {
	// Create a wait group to wait for all Goroutines to complete
	var wg sync.WaitGroup

	// Create a channel to limit the rate of requests
	rateLimiter := time.Tick(time.Second / time.Duration(params.RateLimit))

	// Create a list of results
	results := make([]Result, 0)

	// Create a buffer to store the results temporarily
	var buffer bytes.Buffer

	urls, err := readURLsFromFile(params.Filename)
	if err != nil {
		log.Println("ProcessURLFromGroup> failed to read URLs from file:", err)
		return
	}

	// Create a Goroutine pool to limit the concurrency
	pool := &sync.Pool{
		New: func() interface{} {
			return make(chan struct{}, params.Processes)
		},
	}

	// Initiate multiple Goroutines for concurrent processing
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Retrieve a token from the channel to control the rate
			<-rateLimiter

			//Retrieve a Goroutine from the Goroutine pool
			token := pool.Get().(chan struct{})

			// Perform the request and processing
			params.Url = url
			result := processURL(params)

			if !isResultEmpty(result) {
				// Add the result to the list
				results = append(results, result)

				jsonData, err := json.Marshal(result)
				if err != nil {
					log.Println("ProcessURLFromGroup> json marshal error:", err)
					return
				}

				fmt.Println(string(jsonData))

				buffer.WriteString(string(jsonData))
				buffer.WriteString("\n")
			} else {
				log.Println(url + " can't get result")
			}

			// Release the Goroutine back to the Goroutine pool
			pool.Put(token)
		}(url)
	}

	// Wait for all Goroutines to complete
	wg.Wait()

	// Save the results to a JSON file
	if params.Res && buffer.Len() > 0 {
		//saveResultsToFile(results, params.ResultFile)
		err := WriteBufferToFile(&buffer, params.ResultFile)
		if err != nil {
			fmt.Println("WriteBufferToFile Error:", err)
			return
		}
	}
}

func ProcessURLFromPipe(params ProcessUrlParams) {
	// Create a wait group to wait for all Goroutines to complete
	var wg sync.WaitGroup

	// Create a channel to limit the rate of requests
	rateLimiter := time.Tick(time.Second / time.Duration(params.RateLimit))

	// Create a list of results
	results := make([]Result, 0)

	// Create a buffer to store the results temporarily
	var buffer bytes.Buffer

	// Create a Goroutine pool to limit the concurrency
	pool := &sync.Pool{
		New: func() interface{} {
			return make(chan struct{}, params.Processes)
		},
	}

	// Initiate multiple Goroutines for concurrent processing
	for _, url := range params.URLPipe {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Retrieve a token from the channel to control the rate
			<-rateLimiter

			// Retrieve a Goroutine from the Goroutine pool
			token := pool.Get().(chan struct{})

			// Perform the request and processing
			params.Url = url
			result := processURL(params)

			if !isResultEmpty(result) {
				// Add the result to the list
				results = append(results, result)

				jsonData, err := json.Marshal(result)
				if err != nil {
					log.Println("ProcessURLFromPipe> json marshal error:", err)
					return
				}

				fmt.Println(string(jsonData))

				buffer.WriteString(string(jsonData))
				buffer.WriteString("\n")
			} else {
				log.Println(url + " can't get result")
				return
			}

			// Release the Goroutine back to the Goroutine pool
			pool.Put(token)
		}(url)
	}

	// Wait for all Goroutines to complete
	wg.Wait()

	// Save the results to a JSON file
	if params.Res && buffer.Len() > 0 {
		//saveResultsToFile(results, params.ResultFile)
		err := WriteBufferToFile(&buffer, params.ResultFile)
		if err != nil {
			fmt.Println("WriteBufferToFile Error:", err)
			return
		}
	}
}
