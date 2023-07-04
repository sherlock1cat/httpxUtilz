package main

import (
	"bufio"
	"flag"
	"fmt"
	"httpxUtilz/cmd"
	"os"
)

var params cmd.ProcessUrlParams

func init() {
	flag.StringVar(&params.Url, "url", "", "URL to process.")
	flag.StringVar(&params.Filename, "urls", "", "File URLs to process.")
	flag.StringVar(&params.Proxy, "proxy", "", "Proxy URL.")
	flag.BoolVar(&params.UseHTTPS, "usehttps", true, "Initiate an HTTPS request.")
	flag.BoolVar(&params.FollowRedirects, "followredirects", true, "Perform a URL request redirection.")
	flag.IntVar(&params.MaxRedirects, "maxredirects", 10, "Maximum number of redirections.")
	flag.StringVar(&params.Method, "method", "GET", "The default request method is GET.")
	flag.BoolVar(&params.RandomUserAgent, "randomuseragent", true, "Whether to use a random User-Agent header.")
	flag.StringVar(&params.Headers, "headers", "", "Customize the request headers.")
	flag.BoolVar(&params.FollowSameHost, "followsamehost", true, "Follow Same Host.")
	flag.IntVar(&params.Timeout, "timeout", 10, "Request url timeout.")
	flag.IntVar(&params.Processes, "processes", 1, "Number of processes.")
	flag.IntVar(&params.RateLimit, "rateLimit", 50, "Rate limit.")
	flag.BoolVar(&params.Res, "res", false, "Default not save result.")
	flag.StringVar(&params.ResultFile, "resultFile", "", "Default save to ./result.json.")
	flag.BoolVar(&params.Passive, "passive", false, "Default not get passive info data.")
	flag.Parse()
}

func main() {
	// Check if the standard input is connected to the terminal
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			params.URLPipe = append(params.URLPipe, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Unable to read from the pipe input:", err)
			return
		}
		cmd.ProcessURLFromPipe(params)
	} else {
		if params.Url != "" {
			cmd.ProcessURLFromLine(params)
		} else if params.Filename != "" {
			cmd.ProcessURLFromGroup(params)
		} else {
			flag.Usage()
		}
	}
}
