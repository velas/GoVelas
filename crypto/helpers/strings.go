package helpers

import (
	"regexp"
	"strings"
)

func smallBeutifiers(in string) string {
	rIP := NewCaseInsensitiveReplacer("ip", "IP")
	out := rIP.Replace(in)

	rCSS := NewCaseInsensitiveReplacer("css", "CSS")
	out = rCSS.Replace(out)

	rURL := NewCaseInsensitiveReplacer("url", "URL")
	out = rURL.Replace(out)

	rAPI := NewCaseInsensitiveReplacer("api", "API")
	out = rAPI.Replace(out)

	rSSL := NewCaseInsensitiveReplacer("ssl", "SSL")
	out = rSSL.Replace(out)
	return out
}

func SnakeCaseToCamelCase(inputUnderScoreStr string) string {
	isToUpper := false

	camelCase := ""
	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	camelCase = smallBeutifiers(camelCase)
	return camelCase
}

type CaseInsensitiveReplacer struct {
	toReplace   *regexp.Regexp
	replaceWith string
}

func NewCaseInsensitiveReplacer(toReplace, replaceWith string) *CaseInsensitiveReplacer {
	return &CaseInsensitiveReplacer{
		toReplace:   regexp.MustCompile("(?i)" + toReplace),
		replaceWith: replaceWith,
	}
}

func (cir *CaseInsensitiveReplacer) Replace(str string) string {
	return cir.toReplace.ReplaceAllString(str, cir.replaceWith)
}
