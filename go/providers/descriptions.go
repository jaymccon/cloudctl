package providers

import "strings"

func GetCmdDetails(provider string) (name string, short string, long string) {
	switch provider {
	case "aws":
		name = "aws"
		short = "TODO: short"
		long = `TODO: long`
	default:
		name = provider
		short = ""
		long = ``
	}
	return strings.ToLower(name), short, long
}
