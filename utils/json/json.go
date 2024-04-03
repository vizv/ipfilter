package json

import (
	"fmt"
	"regexp"
)

const reJsonTemplatePrefix = `"%s"\s*:\s*`
const reJsonTemplateString = reJsonTemplatePrefix + `"([^"]+)"`
const reJsonTemplateBoolean = reJsonTemplatePrefix + `(true|false)`

func GetJsonValueString(json string, key string) string {
	re := regexp.MustCompile(fmt.Sprintf(reJsonTemplateString, regexp.QuoteMeta(key)))
	match := re.FindStringSubmatch(json)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

func GetJsonValueBoolean(json string, key string) bool {
	re := regexp.MustCompile(fmt.Sprintf(reJsonTemplateBoolean, regexp.QuoteMeta(key)))
	match := re.FindStringSubmatch(json)
	if len(match) == 2 {
		return match[1] == "true"
	}
	return false
}
