package utilz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

func MatchResponseWithJSONRules(response string, rulesFiles string) (matches map[string]string) {
	var rules map[string]string

	jsonRules, err := ioutil.ReadFile(rulesFiles)
	if err != nil {
		log.Fatal("MatchResponseWithJSONRules> Failed to process the JSON file：", err)
		return
	}

	err = json.Unmarshal([]byte(jsonRules), &rules)
	if err != nil {
		fmt.Println("MatchResponseWithJSONRules> Error parsing JSON rules:", err)
		return nil
	}

	matches = make(map[string]string)

	for key, regexPattern := range rules {
		re := regexp.MustCompile(regexPattern)
		match := re.FindString(response)
		if match != "" {
			matches[key] = match
		}
	}

	return matches
}
