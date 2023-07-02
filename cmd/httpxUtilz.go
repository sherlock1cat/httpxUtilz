package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	httpxUtilz "httpxUtilz/utilz"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

type Result struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Banner      string `json:"banner"`
	StatusCode  int    `json:"statusCode"`
	CName       string `json:"cname"`
	IP          string `json:"ip"`
	Alive       int    `json:"alive"`
	Cdn         int    `json:"cdn"`
	CdnByIP     bool   `json:"cdn_by_ip"`
	CdnByHeader string `json:"cdn_by_header"`
	CdnByCidr   bool   `json:"cdn_by_cidr"`
	CdnByAsn    bool   `json:"cdn_by_asn"`
	CdnByCName  bool   `json:"cdn_by_cname"`
	Cidr        string `json:"cidr"`
	Asn         string `json:"asn"`
	Org         string `json:"org"`
	Addr        string `json:"addr"`
}

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("readURLsFromFile: Open File Error", err)
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
		log.Fatal("readURLsFromFile: Read File Error", err)
		return nil, err
	}

	return urls, nil
}

func isResultEmpty(result Result) bool {
	emptyResult := Result{}

	// Compare two structs for equality using the DeepEqual function
	return reflect.DeepEqual(result, emptyResult)
}

func saveResultsToFile(results []Result, resultFile string) {
	// The default path for the result file is "./result.json"
	if resultFile == "" {
		resultFile = "./result.json"
	}

	// Create the result file
	file, err := os.OpenFile(resultFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("create result file error:", err)
		return
	}
	defer file.Close()

	// Convert the result to a JSON string.
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Println("json marshal error:", err)
		return
	}

	// Write the JSON string to a file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("write result to file error:", err)
		return
	}

	log.Println("results saved to", resultFile)
}

func processURL(url, proxy string, usehttps bool, followredirects bool, maxredirects int, method string, randomuseragent bool, headers string, followsamehost bool, timeout int) (result Result) {
	config := httpxUtilz.RequestClientConfig{
		ProxyURL:        proxy,
		UseHTTPS:        usehttps,
		FollowRedirects: followredirects,
		MaxRedirects:    maxredirects,
		Method:          method,
		RandomUserAgent: randomuseragent,
		Headers: map[string]string{
			"User-Agent": headers,
		},
		FollowSameHost: followsamehost,
		Timeout:        time.Duration(timeout),
	}

	resp, err := config.GetResponseByUrl(url)
	if err != nil {
		log.Println("request error: ", err)
		return
	}

	title := config.GetTitleByResponse(resp)
	Banner := config.GetServerByResponse(resp)
	statuscode := config.GetStatusByResponse(resp)
	cname, ips := config.GetCNameIPByDomain(url, "./data/vaildResolvers.txt")
	alive := config.GetAliveByResponse(resp)
	if len(ips) == 0 {
		return Result{}
	}
	cidr, asn, org, addr := config.GetAsnInfoByIp(ips, proxy)

	var (
		cdn         int
		cdnbyip     bool
		cdnbyheader string
		cdnbycidr   bool
		cdnbyasn    bool
		cdnbycname  bool
	)
	if len(ips) > 0 {
		cdn, cdnbyip, cdnbyheader, cdnbycidr, cdnbyasn, cdnbycname = config.GetCdnInfoByAll(
			resp, ips, "./data/cdn_header_keys.json",
			cidr, "./data/cdn_ip_cidr.json",
			asn, "./data/cdn_asn_list.json",
			cname, "./data/cdn_cname_keywords.json")
	}

	result = Result{
		Url:         url,
		Title:       title,
		Banner:      Banner,
		StatusCode:  statuscode,
		CName:       cname,
		IP:          ips,
		Alive:       alive,
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
	return
}

func ProcessURLFromLine(url, proxy string, usehttps bool, followredirects bool, maxredirects int, method string, randomuseragent bool, headers string, followsamehost bool, timeout int, rateLimit int, res bool, resultFile string) {
	// Create a wait group to wait for all Goroutines to complete
	var wg sync.WaitGroup

	// Create a channel to limit the rate of requests
	rateLimiter := time.Tick(time.Second / time.Duration(rateLimit))

	// Create a list of results
	results := make([]Result, 0)

	// Initiate multiple Goroutines for concurrent processing
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Retrieve a token from the channel to control the rate
			<-rateLimiter

			// Perform the request and processing
			result := processURL(url, proxy, usehttps, followredirects, maxredirects, method,
				randomuseragent, headers, followsamehost, timeout)

			if !isResultEmpty(result) {
				// Add the result to the list.
				results = append(results, result)

				jsonData, err := json.Marshal(result)
				if err != nil {
					log.Println("json marshal error:", err)
					return
				}

				fmt.Println(string(jsonData))
			} else {
				log.Println(url + " can't get result")
			}

		}()
	}

	// Wait for all Goroutines to complete
	wg.Wait()

	// Save the results to a JSON file
	if res {
		saveResultsToFile(results, resultFile)
	}

}

func ProcessURLFromGroup(filename, proxy string, usehttps bool, followredirects bool, maxredirects int, method string, randomuseragent bool, headers string, followsamehost bool, timeout int, processes int, rateLimit int, res bool, resultFile string) {
	// Create a wait group to wait for all Goroutines to complete
	var wg sync.WaitGroup

	// Create a channel to limit the rate of requests
	rateLimiter := time.Tick(time.Second / time.Duration(rateLimit))

	// Create a list of results
	results := make([]Result, 0)

	urls, err := readURLsFromFile(filename)
	if err != nil {
		log.Println("failed to read URLs from file:", err)
		return
	}

	// Create a Goroutine pool to limit the concurrency
	pool := &sync.Pool{
		New: func() interface{} {
			return make(chan struct{}, processes)
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
			defer pool.Put(token)

			// Perform the request and processing
			result := processURL(url, proxy, usehttps, followredirects, maxredirects, method,
				randomuseragent, headers, followsamehost, timeout)

			if !isResultEmpty(result) {
				// Add the result to the list
				results = append(results, result)

				jsonData, err := json.Marshal(result)
				if err != nil {
					log.Println("json marshal error:", err)
					return
				}

				fmt.Println(string(jsonData))
			} else {
				log.Println(url + " can't get result")
			}

			// Release the Goroutine back to the Goroutine pool
			token <- struct{}{}
		}(url)
	}

	// Wait for all Goroutines to complete
	wg.Wait()

	// Save the results to a JSON file
	if res {
		saveResultsToFile(results, resultFile)
	}
}

func ProcessURLFromPipe(urlPipe []string, proxy string, usehttps bool, followredirects bool, maxredirects int, method string, randomuseragent bool, headers string, followsamehost bool, timeout int, processes int, rateLimit int, res bool, resultFile string) {
	// Create a wait group to wait for all Goroutines to complete
	var wg sync.WaitGroup

	// Create a channel to limit the rate of requests
	rateLimiter := time.Tick(time.Second / time.Duration(rateLimit))

	// Create a list of results
	results := make([]Result, 0)

	// Create a Goroutine pool to limit the concurrency
	pool := &sync.Pool{
		New: func() interface{} {
			return make(chan struct{}, processes)
		},
	}

	// Initiate multiple Goroutines for concurrent processing
	for _, url := range urlPipe {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Retrieve a token from the channel to control the rate
			<-rateLimiter

			// Retrieve a Goroutine from the Goroutine pool
			token := pool.Get().(chan struct{})
			defer pool.Put(token)

			// Perform the request and processing
			result := processURL(url, proxy, usehttps, followredirects, maxredirects, method,
				randomuseragent, headers, followsamehost, timeout)

			if !isResultEmpty(result) {
				// Add the result to the list
				results = append(results, result)

				jsonData, err := json.Marshal(result)
				if err != nil {
					log.Println("json marshal error:", err)
					return
				}

				fmt.Println(string(jsonData))
			} else {
				log.Println(url + " can't get result")
			}

			// Release the Goroutine back to the Goroutine pool
			token <- struct{}{}
		}(url)
	}

	// Wait for all Goroutines to complete
	wg.Wait()

	// Save the results to a JSON file
	if res {
		saveResultsToFile(results, resultFile)
	}
}
