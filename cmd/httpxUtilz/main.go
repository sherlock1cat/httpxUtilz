package main

import (
	"bufio"
	"flag"
	"fmt"
	"httpxUtilz/cmd"
	"os"
)

var (
	url             string
	urls            string
	proxy           string
	usehttps        bool
	followredirects bool
	maxredirects    int
	method          string
	randomuseragent bool
	headers         string
	followsamehost  bool
	timeout         int
	processes       int
	rateLimit       int
	res             bool
	resultFile      string
)

func init() {
	flag.StringVar(&url, "url", "", "URL to process.")
	flag.StringVar(&urls, "urls", "", "File URLs to process.")
	flag.StringVar(&proxy, "proxy", "", "Proxy URL.")
	flag.BoolVar(&usehttps, "usehttps", true, "Initiate an HTTPS request.")
	flag.BoolVar(&followredirects, "followredirects", true, "Perform a URL request redirection.")
	flag.IntVar(&maxredirects, "maxredirects", 10, "Maximum number of redirections.")
	flag.StringVar(&method, "method", "GET", "The default request method is GET.")
	flag.BoolVar(&randomuseragent, "randomuseragent", true, "Whether to use a random User-Agent header.")
	flag.StringVar(&headers, "headers", "", "Customize the request headers.")
	flag.BoolVar(&followsamehost, "followsamehost", true, "Follow Same Host.")
	flag.IntVar(&processes, "processes", 1, "Number of processes.")
	flag.IntVar(&rateLimit, "rateLimit", 100, "Rate limit.")
	flag.BoolVar(&res, "res", false, "Default not save result.")
	flag.StringVar(&resultFile, "resultFile", "", "Default save to ./result.json.")
	flag.Parse()
}

func main() {
	// Check if the standard input is connected to the terminal
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		var urlPipe []string
		for scanner.Scan() {
			urlPipe = append(urlPipe, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Unable to read from the pipe input:", err)
			return
		}
		cmd.ProcessURLFromPipe(urlPipe, proxy, usehttps, followredirects, maxredirects, method, randomuseragent, headers, followredirects, timeout, processes, rateLimit, res, resultFile)
	} else {
		if url != "" {
			cmd.ProcessURLFromLine(url, proxy, usehttps, followredirects, maxredirects, method, randomuseragent, headers, followredirects, timeout, rateLimit, res, resultFile)
		} else if urls != "" {
			cmd.ProcessURLFromGroup(urls, proxy, usehttps, followredirects, maxredirects, method, randomuseragent, headers, followredirects, timeout, processes, rateLimit, res, resultFile)
		} else {
			flag.Usage()
		}
	}
}
