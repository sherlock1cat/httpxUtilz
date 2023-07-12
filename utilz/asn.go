package utilz

import (
	"encoding/json"
	asnmap "github.com/projectdiscovery/asnmap/libs"
	"log"
	"strings"
)

type AsnData struct {
	AsNumber  string   `json:"as_number"`
	AsName    string   `json:"as_name"`
	AsCountry string   `json:"as_country"`
	AsRange   []string `json:"as_range"`
}

func handleInput(client *asnmap.Client, item string) *AsnData {
	results, err := client.GetData(item)
	if err != nil {
		log.Println(err)
		return nil
	}
	output, err := asnmap.GetFormattedDataInJson(results)
	if err != nil {
		log.Println(err)
		return nil
	}
	var data AsnData
	if len(output) > 0 {
		//log.Println(fmt.Sprintf("handleInput: %s: %s", item, string(output)))

		err := json.Unmarshal([]byte(output), &data)
		if err != nil {
			log.Printf("handleInput: Failed to parse JSON：%v", err)
			return nil
		}
	}
	return &data
}

func GetAsnInfoByIps(ips string, proxy string) (cidr, asn, org, addr string) {
	client, err := asnmap.NewClient()
	if proxy != "" {
		proxys := []string{proxy}
		_, err = client.SetProxy(proxys)
	}
	if err != nil {
		log.Println("GetAsnInfoByIps> asnmap new client ", err)
		return
	}
	items := strings.Split(ips, ",")

	for _, item := range items { // Retrieve the first result.
		data := handleInput(client, item)
		if data != nil {
			cidr = strings.Join(data.AsRange, ",")
			asn = data.AsNumber
			org = data.AsName
			addr = data.AsCountry
		} else {
			log.Printf("GetAsnInfoByIps> asnmap can't get data by %s", ips)
			cidr = "Na"
			asn = "Na"
			org = "Na"
			addr = "Na"
		}
		return
	}
	return
}

func GetIpsByAsnmap(url string) (ips string) {
	domain, err := GetSubDomain(url)
	if err != nil {
		log.Printf("GetIpsByAsnmap> %s getsubdomain failed, check url format.", url)
		return
	}
	resolvedIps, err := asnmap.ResolveDomain(domain)
	if err != nil {
		log.Fatal(err)
	}
	ips = strings.Join(resolvedIps, ",")
	return
}
