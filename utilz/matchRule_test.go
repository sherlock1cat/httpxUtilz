package utilz

import "testing"

func TestMatchResponseWithJSONRules(t *testing.T) {

	response := `Sample response string
	Some random text
	JDBC Connection: jdbc:mysql://localhost:3306
	AccessKeyId: abcdef123456
	Some other text
	Amazon AWS URL: s3.amazonaws.com
	Some more random text
	`
	rulesFile := "../data/regex_MayVul.json"

	matches := make(map[string]string)
	matches = MatchResponseWithJSONRules(response, rulesFile)

	// Verify if the matching results meet the expected criteria.
	expectedMatches := map[string]string{
		"JDBC Connection": "jdbc:mysql://localhost:3306",
		"OSS":             "AccessKeyId",
		"Amazon AWS URL":  "s3.amazonaws.com",
	}

	for key, expectedValue := range expectedMatches {
		actualValue, ok := matches[key]
		if !ok {
			t.Errorf("Expected key '%s' not found in matches", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected value '%s' for key '%s', but got '%s'", expectedValue, key, actualValue)
		}
	}
}
