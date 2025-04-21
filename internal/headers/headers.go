package headers

import (
	"strings"
	"fmt"
	"regexp"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return map[string]string{}
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
	i := strings.Index(string(data), crlf)
	// clrf not found, not enough data yet
	if i == -1 {
		return 0, false, nil
	}
	// clrf at beginning, headers are done
	if i == 0 {
		//so, the crlf should be consumed...
		return len(crlf), true, nil
	}

	rawHeader := strings.Split(string(data), crlf)[0]
	trimmed := strings.TrimSpace(rawHeader)
	splitHeader := strings.Split(trimmed, " ")
	if len(splitHeader) != 2 {
		return 0, false, fmt.Errorf("Invalid header length")
	}
	rawFieldName := splitHeader[0]
	fieldValue := splitHeader[1]
	if strings.Index(rawFieldName, ":") != len(rawFieldName)-1 || len(strings.Split(rawFieldName, " ")) != 1 {
		return 0, false, fmt.Errorf("Invalid field name")
	}

	fieldName := strings.ToLower(strings.Split(rawFieldName, ":")[0])
	re := regexp.MustCompile("([^A-Za-z0-9!#$%&'*+.^_`|~-])+")
	if re.Match([]byte(fieldName)) {
		return 0, false, fmt.Errorf("Invalid field name")
	}

	if _, ok := (*h)[fieldName]; ok {
		(*h)[fieldName] = (*h)[fieldName] + ", " + fieldValue
	} else {
		(*h)[fieldName] = fieldValue
	}
	
	return len(rawHeader) + len(crlf), false, nil
}