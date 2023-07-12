package utilz

import (
	"github.com/miekg/dns"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"io/ioutil"
	"log"
	"math"
	"strings"
)

func FileContentToList(filePath string) []string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("fail to open file: " + filePath)
		return nil
	}
	contentList := strings.Split(string(fileContent), "\n")
	return contentList
}

func UniqueStrList(strList []string) []string {
	uniqList := make([]string, 0)
	tempMap := make(map[string]bool, len(strList))
	for _, v := range strList {
		if tempMap[v] == false && len(v) > 0 {
			tempMap[v] = true
			uniqList = append(uniqList, v)
		}
	}
	return uniqList
}

func DnsxClient(domain string, resolversFile string) *dnsx.DNSX {
	validResolversList := UniqueStrList(FileContentToList(resolversFile))

	DefaultOptions := dnsx.Options{
		BaseResolvers:     validResolversList,
		MaxRetries:        5,
		QuestionTypes:     []uint16{dns.TypeA},
		TraceMaxRecursion: math.MaxUint16,
		Hostsfile:         true,
	}

	dnsxClient, err := dnsx.New(DefaultOptions)
	if err != nil {
		log.Println("DnsxClient> ", err)
		return nil
	}
	return dnsxClient
}

func GetCnameIPsByDomain(url string, resolversFile string) (cname, ips []string) {

	domain, err := GetSubDomain(url)
	if err != nil {
		log.Printf("GetCnameIPsByDomain> %s getsubdomain failed, check url format.", url)
		return
	}
	dnsxClient := DnsxClient(domain, resolversFile)

	dnsxResult, _ := dnsxClient.QueryOne(domain)
	CName := dnsxResult.CNAME
	IPs := dnsxResult.A

	cname = CName
	ips = IPs

	return
}
