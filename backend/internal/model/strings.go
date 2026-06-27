package model

import (
	"strings"
)

type Strings []string

func (s *Strings) Scan(src any) error {
	if src == nil {
		*s = nil
		return nil
	}
	var raw string
	switch v := src.(type) {
	case []byte:
		raw = string(v)
	case string:
		raw = v
	default:
		*s = Strings{}
		return nil
	}
	raw = strings.TrimSpace(raw)
	if len(raw) < 2 || raw[0] != '{' || raw[len(raw)-1] != '}' {
		*s = Strings{}
		return nil
	}
	inner := raw[1 : len(raw)-1]
	if inner == "" {
		*s = Strings{}
		return nil
	}
	var result []string
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		if c == '"' {
			i++
			var elem strings.Builder
			for i < len(inner) {
				if inner[i] == '\\' {
					i++
					if i < len(inner) {
						elem.WriteByte(inner[i])
					}
				} else if inner[i] == '"' {
					break
				} else {
					elem.WriteByte(inner[i])
				}
				i++
			}
			result = append(result, elem.String())
		} else if c == ',' || c == ' ' {
			continue
		} else {
			start := i
			for i < len(inner) && inner[i] != ',' {
				i++
			}
			result = append(result, strings.TrimSpace(inner[start:i]))
			i--
		}
	}
	*s = Strings(result)
	return nil
}
