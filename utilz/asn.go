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
		log.Fatal(err)
	}
	output, err := asnmap.GetFormattedDataInJson(results)
	if err != nil {
		log.Fatal(err)
	}
	var data AsnData
	if len(output) > 0 {
		//log.Println(fmt.Sprintf("handleInput: %s: %s", item, string(output)))

		err := json.Unmarshal([]byte(output), &data)
		if err != nil {
			log.Fatalf("handleInput: Failed to parse JSON：%v", err)
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
		log.Fatal("GetAsnInfoByIps: asnmap new client ", err)
	}
	items := strings.Split(ips, ",")

	for _, item := range items { // Retrieve the first result.
		data := handleInput(client, item)

		cidr = strings.Join(data.AsRange, ",")
		asn = data.AsNumber
		org = data.AsName
		addr = data.AsCountry

		return cidr, asn, org, addr
	}
	return
}
