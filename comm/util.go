package comm

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"runtime/debug"
)

var (
	specialRegular = "[^(a-zA-Z0-9\u4e00-\u9fa5)]"
	specialRegexp  = regexp.MustCompile(specialRegular)
)

func isASCIISpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// TrimString returns s without leading and trailing ASCII space.
func TrimString(s string) string {
	for len(s) > 0 && isASCIISpace(s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && isASCIISpace(s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}

//Round8 X round成最大的8的倍数
func Round8(x int) int {
	if x <= 0 {
		return 0
	}

	return (x*8 + 1) / 8
}

func PanicRecovery(quit bool) {
	var err error
	if r := recover(); r != nil {
		switch x := r.(type) {
		case string:
			err = errors.New(x)
			break
		case error:
			err = x
			break
		default:
			err = errors.New("Unknown panic")
			break
		}

		debug.PrintStack()
		log.Warn(string(debug.Stack()))
		log.Infoln("Panic :", err.Error())

		if quit {
			os.Exit(101)
		}
	}
}

func FilterSpecialSymbols(text string) string {
	return specialRegexp.ReplaceAllString(text, "")
}
