package utilz

import (
	"encoding/json"
	asnmap "github.com/projectdiscovery/asnmap/libs"
	"log"
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

func GetAsnInfoByIps(ips []string, proxy string) (cidr, asn, org, addr []string) {
	client, err := asnmap.NewClient()
	if proxy != "" {
		proxys := []string{proxy}
		_, err = client.SetProxy(proxys)
	}
	if err != nil {
		log.Println("GetAsnInfoByIps> asnmap new client ", err)
		return
	}

	for _, item := range ips { // Retrieve result have value break.
		data := handleInput(client, item)
		if data != nil && len(cidr) == 0 && len(asn) == 0 && len(org) == 0 && len(addr) == 0 {
			cidr = data.AsRange
			asn = append(asn, data.AsNumber)
			org = append(org, data.AsName)
			addr = append(addr, data.AsCountry)
		} else {
			log.Printf("GetAsnInfoByIps> asnmap can't get data by %s", item)
			//cidr = "Na"
			//asn = "Na"
			//org = "Na"
			//addr = "Na"
		}
	}
	return
}

func GetIpsByAsnmap(url string) (ips []string) {
	domain, err := GetSubDomain(url)
	if err != nil {
		log.Printf("GetIpsByAsnmap> %s getsubdomain failed, check url format.", url)
		return
	}
	ips, err = asnmap.ResolveDomain(domain)
	if err != nil {
		log.Fatal(err)
	}

	return
}
