package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
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
const token string = `$TIMESERIAL`

// TemplateTokenRegexp - regular expression matching $TIMESERIAL
var TemplateTokenRegexp = regexp.MustCompile(regexp.QuoteMeta(token))

// SerialNumberRegexp - regular expression matching date-based DNS serial nums
var SerialNumberRegexp = regexp.MustCompile(`\b20\d\d[01]\d[0123]\d{3}\b`)

// fileSearch returns a boolean indicating that the contents matched the regex
// and a string containing the entire contents of the file
func fileSearch(filename string, r *regexp.Regexp) (bool, string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, "", err
	}
	contents := string(bytes)

	matched := false
	if r.FindIndex(bytes) != nil {
		matched = true
	}

	return matched, contents, nil
}

func fileWrite(filename string, contents string, mode os.FileMode) {
	err := ioutil.WriteFile(filename, []byte(contents), mode)
	if err != nil {
		fmt.Print(err)
	}
}

// TimeSerialize install time-based DNS serial numbers in the given DNS zone
// file
func TimeSerialize(filename string, regex *regexp.Regexp, useNow bool,
	useUTC bool, inPlace bool) (matched bool, contents string, err error) {
	// long function declarations suck
	matched, contents, err = fileSearch(filename, regex)
	if err != nil || !matched {
		return false, contents, err
	}

	// should we stat the file? inPlace needs Mode, !useNow needs ModTime
	var stats os.FileInfo
	if inPlace || !useNow {
		// collect some info about the file
		stats, err = os.Stat(filename)
		if err != nil {
			return true, contents, err
		}
	}

	// calculate the time to use
	var t time.Time
	if useNow {
		t = time.Now()
	} else {
		t = stats.ModTime()
	}
	if useUTC {
		t = t.UTC()
	}

	// do the text replacement
	replacement := regex.ReplaceAllLiteralString(contents, TimeSerial(t))

	if inPlace && replacement != contents {
		fileWrite(filename, replacement, stats.Mode())
	}

	return true, replacement, nil
}

func main() {
	// Argument processing
	inPlace := false
	replaceExisting := false
	useNow := false
	useUTC := false
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helpText, os.Args[0], token)
		flag.PrintDefaults()
	}
	flag.BoolVar(&inPlace, "i", inPlace,
		"Edit the given files in-place")
	flag.BoolVar(&replaceExisting, "e", replaceExisting, fmt.Sprintf(
		"Replace existing serial numbers, instead of replacing %s", token))
	flag.BoolVar(&useNow, "n", useNow,
		"Use the currrent time, instead of file mtimes")
	flag.BoolVar(&useUTC, "u", useUTC,
		"Use the UTC time zone instead of localtime")
	flag.Parse()

	// init
	var regex *regexp.Regexp
	if replaceExisting {
		regex = SerialNumberRegexp
	} else {
		regex = TemplateTokenRegexp
	}

	// Things doing
	for _, filename := range flag.Args() {
		matched, content, err := TimeSerialize(filename, regex, useNow,
			useUTC, inPlace)
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
