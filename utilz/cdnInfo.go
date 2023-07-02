package utilz

import (
	"encoding/json"
	"github.com/projectdiscovery/cdncheck"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func ReadJSONFile(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("ReadJSONFile: Failed to read the file：%v", err)
		return nil, err
	}

	var value []string
	err = json.Unmarshal(data, &value)
	if err != nil {
		log.Fatalf("ReadJSONFile: Failed to parse JSON：%v", err)
		return nil, err
	}

	return value, nil
}

func ReadCNameJSONFile(filename string) (map[string]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("ReadJSONFile: Failed to read the file：%v", err)
		return nil, err
	}

	var value map[string]string
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		log.Fatalf("ReadJSONFile: Failed to parse JSON：%v", err)
		return nil, err
	}

	return value, nil
}

func GetCDNInfoByIps(ips string) (cdn int, CdnInfo string) {
	client := cdncheck.New()
	ipList := strings.Split(ips, ",")

	for _, ipStr := range ipList {
		ip := net.ParseIP(strings.TrimSpace(ipStr))
		if ip == nil {
			//log.Println("GetCDNInfoByIps: Invalid IP address: ", ipStr)
			continue
		}
		matched, val, err := client.CheckCDN(ip)
		if err != nil {
			//log.Println("GetCDNInfoByIps: ", err)
			continue
		}

		if matched {
			//log.Println(ip, "is a", val)
			cdn = 1
			CdnInfo = val
			return
		} else {
			//log.Println(ip, "is not a CDN")
			cdn = 0
			CdnInfo = ""
			continue
		}
	}
	return
}

func GetCDNInfoByHeader(resp *Response, CdnHeaderfilename string) (cdn int, CdnInfo string) {
	cdnHeaders, err := ReadJSONFile(CdnHeaderfilename)
	if err != nil {
		log.Fatal("GetCDNInfoByHeader: Failed to process the JSON file：", err)
		return
	}
	for _, header := range cdnHeaders {
		if cdn == 1 {
			if value := resp.Headers.Get(header); value != "" {
				CdnInfo += value + ","
			}
		}
		if value := resp.Headers.Get(header); value != "" {
			cdn = 1
			CdnInfo = "CDN_Headers: " + header + ","
		}
	}
	return
}

func GetCDNInfoByCidr(cidr string, CdnCidrfilename string) (cdn int, CdnInfo string) {
	cdnCidrs, err := ReadJSONFile(CdnCidrfilename)
	if err != nil {
		log.Fatal("GetCDNInfoByCidr: Encountered an error while processing the JSON file：", err)
		return
	}
	cidrRange := strings.Split(cidr, ",")
	for _, checkCidr := range cidrRange {
		for _, value := range cdnCidrs {
			if value == checkCidr && value != "" {
				cdn = 1
				CdnInfo = "CDN_Cidr: true" + ","
				return
			}
		}
	}

	return
}

func GetCDNInfoByAsn(asn, CdnAsnfilename string) (cdn int, CdnInfo string) {
	cdnAsn, err := ReadJSONFile(CdnAsnfilename)
	if err != nil {
		log.Fatal("GetCDNInfoByAsn: Failed to handle the JSON file：", err)
		return
	}

	asnRange := strings.Split(asn, ",")
	for _, checkAsn := range asnRange {
		for _, value := range cdnAsn {
			if value == checkAsn && value != "" {
				cdn = 1
				CdnInfo = "CDN_Asn: true" + ","
				return
			}
		}
	}

	return
}

func GetCDNInfoByCName(cname, CdnCNamefilename string) (cdn int, CdnInfo string) {
	cnameMap, err := ReadCNameJSONFile(CdnCNamefilename)
	if err != nil {
		log.Fatal("GetCDNInfoByAsn: 处理Json文件失败：", err)
		return
	}

	cnameRange := strings.Split(cname, ",")
	for _, checkCName := range cnameRange {
		Info, ok := cnameMap[checkCName]
		if !ok {
			continue
		}
		if cdn == 0 {
			cdn = 1
			CdnInfo = "Cdn_cname: " + Info + ","
		} else {
			CdnInfo += Info + ","
		}

	}

	return
}