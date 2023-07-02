package utilz

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//ToDoList
//	Alive      int    `bson:"alive"` Done
//	URL        string `bson:"url"`
//	Subdomain  string `bson:"subdomain"`
//	Level      int    `bson:"level"`
//	CName      string `bson:"cname"` Done
//	IP         string `bson:"ip"` Done
//	CDN        int    `bson:"cdn"` Done
//	CdnInfo	   string `bson:"cdn_info"` Done
//	Port       int    `bson:"port"`
//	Status     int    `bson:"status"` Done
//	Reason     string `bson:"reason"`
//	Title      string `bson:"title"` Done
//	Banner     string `bson:"banner"` Done
//	CIDR       string `bson:"cidr"`
//	ASN        string `bson:"asn"`
//	Org        string `bson:"org"`
//	Addr       string `bson:"addr"`
//	ISP        string `bson:"isp"`
//	Source     string `bson:"source"`
//	CreateDate string `bson:"create_date"`
//	UpdateDate string `bson:"update_date"`
//}

type Response struct {
	Raw     string
	Data    []byte
	Headers http.Header
	Status  int
}

// httpxUtilz -u hackerone.com -title -server -status-code -tech-detect -ip -cname -asn -cdn -duc

func (config *RequestClientConfig) GetResponseByUrl(targetUrl string) (*Response, error) {
	target, err := parseUrl(targetUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	client := NewRequestClient(*config)
	resp, err := client.Get(target)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Response{
		Raw:     string(body),
		Data:    body,
		Headers: resp.Header,
		Status:  resp.StatusCode,
	}, nil
}

func (config *RequestClientConfig) GetAliveByResponse(resp *Response) (alive int) {
	alive = 1
	if resp.Status == http.StatusNotFound || resp.Status == http.StatusBadGateway {
		alive = 0
	}
	return
}

func (config *RequestClientConfig) GetTitleByResponse(resp *Response) (title string) {
	title = ExtractTitle(resp)
	return
}

func (config *RequestClientConfig) GetServerByResponse(resp *Response) (Banner string) {
	server := resp.Headers.Get("Server")
	if len(server) > 0 {
		Banner = server + ","
	}
	via := resp.Headers.Get("Via")
	if len(via) > 0 {
		Banner += via + ","
	}
	power := resp.Headers.Get("X-Powered-By")
	if len(power) > 0 {
		Banner += power
	}
	return
}

func (config *RequestClientConfig) GetServerAllHeaderByResponse(resp *Response) (responseheader []string) {
	for key, values := range resp.Headers {
		for _, value := range values {
			responseheader = append(responseheader, fmt.Sprintf("%s: %s", key, value))
		}
	}

	return
}

func (config *RequestClientConfig) GetStatusByResponse(resp *Response) (Status int) {
	Status = resp.Status
	return
}

func (config *RequestClientConfig) GetCNameIPByDomain(domain string, resolversFile string) (cname string, ips string) {
	cname, ips = GetCnameIPsByDomain(domain, resolversFile)
	return
}

func (config *RequestClientConfig) GetCdnInfoByAll(resp *Response, ips, CdnHeaderfilename, cidr, CdnCidrfilename, asn, CdnAsnfilename, cname, CdnCNamefilename string) (cdn int, cdnbyip bool, cdnbyheader string, cdnbycidr, cdnbyasn, cdnbycname bool) {
	cdn, cdnbyip = GetCDNInfoByIps(ips)

	cdn, cdnbyheader = GetCDNInfoByHeader(resp, CdnHeaderfilename)

	cdn, cdnbycidr = GetCDNInfoByCidr(cidr, CdnCidrfilename)

	cdn, cdnbyasn = GetCDNInfoByAsn(asn, CdnAsnfilename)

	cdn, cdnbycname = GetCDNInfoByCName(cname, CdnCNamefilename)

	return
}

func (config *RequestClientConfig) GetAsnInfoByIp(ips string, proxy string) (cidr, asn, org, addr string) {
	cidr, asn, org, addr = GetAsnInfoByIps(ips, proxy)
	return
}
