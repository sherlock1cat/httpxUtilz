package utilz

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Raw                    string
	Data                   []byte
	Headers                http.Header
	Status                 int
	ContentLength          int64
	ContentLengthByAllBody int64
}

func (config *RequestClientConfig) GetResponseByUrl(targetUrl string) (*Response, error) {
	target, err := parseUrl(targetUrl)
	if err != nil {
		log.Println("GetResponseByUrl: ", err)
		return nil, err
	}
	client := NewRequestClient(*config)
	resp, err := client.Get(target)
	if err != nil {
		log.Println("GetResponseByUrl: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("GetResponseByUrl: ", err)
		return nil, err
	}

	return &Response{
		Raw:                    string(body),
		Data:                   body,
		Headers:                resp.Header,
		Status:                 resp.StatusCode,
		ContentLength:          resp.ContentLength,
		ContentLengthByAllBody: int64(len(body)),
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
	if len(title) == 0 {
		title = "NA"
	}
	return
}

func (config *RequestClientConfig) GetBannerByResponse(resp *Response) (server, via, power string) {
	if len(resp.Headers.Get("Server")) > 0 {
		server = resp.Headers.Get("Server")
		if len(server) == 0 {
			server = "NA"
		}
	}

	if len(resp.Headers.Get("Via")) > 0 {
		via = resp.Headers.Get("Via")
		if len(via) == 0 {
			via = "NA"
		}
	}

	if len(resp.Headers.Get("X-Powered-By")) > 0 {
		power = resp.Headers.Get("X-Powered-By")
		if len(power) == 0 {
			power = "NA"
		}
	}

	return
}

func (config *RequestClientConfig) GetServerAllHeaderByResponse(resp *Response) (responseHeader []string) {
	for key, values := range resp.Headers {
		for _, value := range values {
			responseHeader = append(responseHeader, fmt.Sprintf("%s: %s", key, value))
		}
	}

	return
}

func (config *RequestClientConfig) GetStatusByResponse(resp *Response) (Status int) {
	Status = resp.Status
	return
}

func (config *RequestClientConfig) GetContentLengthByResponse(resp *Response) (contentLength int64) {
	contentLength = resp.ContentLength
	return
}

func (config *RequestClientConfig) GetContentLengthAllBodyByResponse(resp *Response) (contentLengthByAllBody int64) {
	contentLengthByAllBody = resp.ContentLengthByAllBody
	return
}

func (config *RequestClientConfig) GetCNameIPByDomain(domain string, resolversFile string) (cname string, ips string) {
	cname, ips = GetCnameIPsByDomain(domain, resolversFile)
	if len(cname) == 0 {
		cname = "NA"
	}
	return
}

func (config *RequestClientConfig) GetCdnInfoByAll(resp *Response, ips, CdnHeaderfilename, cidr, CdnCidrfilename, asn, CdnAsnfilename, cname, CdnCNamefilename string) (cdn int, cdnbyip bool, cdnbyheader []string, cdnbycidr, cdnbyasn, cdnbycname bool) {
	_, cdnbyip = GetCDNInfoByIps(ips)

	_, cdnbyheader = GetCDNInfoByHeader(resp, CdnHeaderfilename)

	_, cdnbycidr = GetCDNInfoByCidr(cidr, CdnCidrfilename)

	_, cdnbyasn = GetCDNInfoByAsn(asn, CdnAsnfilename)

	_, cdnbycname = GetCDNInfoByCName(cname, CdnCNamefilename)

	if (cdnbyip) || (len(cdnbyheader) > 0) || (cdnbycidr) || (cdnbyasn) || (cdnbycname) {
		cdn = 1
	}

	return
}

func (config *RequestClientConfig) GetAsnInfoByIp(ips string, proxy string) (cidr, asn, org, addr string) {
	cidr, asn, org, addr = GetAsnInfoByIps(ips, proxy)
	return
}

func (config *RequestClientConfig) GetMayVulInfoByRespone(resp *Response, rulesFiles string) (vulMatches map[string]string) {
	vulMatches = MatchResponseWithJSONRules(resp.Raw, rulesFiles)

	return
}
