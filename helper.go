package lichv

import "regexp"

func IsMatch(text string, filter string) bool {
	reg := regexp.MustCompile(filter)
	result := reg.FindAllString(text, -1)
	if len(result) > 0 {
		return true
	}else{
		return false
	}
}
