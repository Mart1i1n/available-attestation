package pointset

import "strings"

var BlockPointSet = []string{
	"BlockDelayForReceiveBlock",
	"BlockBeforeBroadCast",
	"BlockAfterBroadCast",
	"BlockBeforeSign",
	"BlockAfterSign",
	"BlockBeforePropose",
	"BlockAfterPropose",
}

var AttestPointSet = []string{
	"AttestBeforeBroadCast",
	"AttestAfterBroadCast",
	"AttestBeforeSign",
	"AttestAfterSign",
	"AttestBeforePropose",
	"AttestAfterPropose",
}

func GetPointByName(name string) string {
	for _, p := range BlockPointSet {
		if strings.ToLower(p) == strings.ToLower(name) {
			return p
		}
	}
	for _, p := range AttestPointSet {
		if strings.ToLower(p) == strings.ToLower(name) {
			return p
		}
	}
	return ""
}
