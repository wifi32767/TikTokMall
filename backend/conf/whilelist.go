package conf

import (
	"fmt"
	"regexp"
)

var WhiteListRe *regexp.Regexp

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
