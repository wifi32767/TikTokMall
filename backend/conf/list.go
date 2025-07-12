package conf

import (
	"fmt"
	"regexp"
)

var (
	WhiteListRe     *regexp.Regexp
	ProtectedListRe *regexp.Regexp
)

func initWhiteList() {
	combinedParttern := ""
	for _, parttern := range conf.WhiteList {
		combinedParttern += parttern + "|"
	}
	if len(combinedParttern) > 1 {
		combinedParttern = combinedParttern[:len(combinedParttern)-1]
	}
	fmt.Println("WhiteListRe:", combinedParttern)
	WhiteListRe = regexp.MustCompile(combinedParttern)
}

func initProtectedList() {
	combinedParttern := ""
	for _, parttern := range conf.ProtectedList {
		combinedParttern += parttern + "|"
	}
	if len(combinedParttern) > 1 {
		combinedParttern = combinedParttern[:len(combinedParttern)-1]
	}
	fmt.Println("ProtectedListRe:", combinedParttern)
	ProtectedListRe = regexp.MustCompile(combinedParttern)
}
