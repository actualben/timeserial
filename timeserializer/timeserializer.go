package timeserializer

import (
	"flag"
	"fmt"
	"github.com/actualben/timeserial/timeserial"
	"os"
	"regexp"
)

const helpText string = `usage: %s [-h] [-e] [-i] [-n] [-u] [file ...]

For each file given on the command line, replace any instances of the template
text %s with a 10-digit number suitable for use as a DNS zone
serial number. The number is based on the current date and time in this
format: YYYYMMDDTT. Here YYYY represents the year including century, MM
represents the zero-padded month number, DD represents the zero-padded day of
the month, and TT represents the time in hundredths-of-a-day
(864 second / ~15 minute intervals) since midnight.

This command is designed to expand the $TIMESERIAL template variable in DNS
zone files as they are placed into read-only container filesystems

`

func main() {
	// Argument processing
	inPlace := false
	replaceExisting := false
	useNow := false
	useUTC := false
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helpText, os.Args[0], timeserial.TemplateToken)
		flag.PrintDefaults()
	}
	flag.BoolVar(&inPlace, "i", inPlace,
		"Edit the given files in-place")
	flag.BoolVar(&replaceExisting, "e", replaceExisting, fmt.Sprintf(
		"Replace existing serial numbers, instead of replacing %s",
		timeserial.TemplateToken))
	flag.BoolVar(&useNow, "n", useNow,
		"Use the currrent time, instead of file mtimes")
	flag.BoolVar(&useUTC, "u", useUTC,
		"Use the UTC time zone instead of localtime")
	flag.Parse()

	// init
	var regex *regexp.Regexp
	if replaceExisting {
		regex = timeserial.SerialNumberRegexp
	} else {
		regex = timeserial.TemplateTokenRegexp
	}

	// Things doing
	for _, filename := range flag.Args() {
		matched, content, err := timeserial.TimeSerialize(filename, regex,
			useNow, useUTC, inPlace)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if inPlace {
			if matched {
				fmt.Println(filename)
			}
			continue
		}

		fmt.Println(content)
	}
}
