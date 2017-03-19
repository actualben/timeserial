package timeserial

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

const TemplateToken string = `$TIMESERIAL`

// TemplateTokenRegexp - regular expression matching $TIMESERIAL
var TemplateTokenRegexp = regexp.MustCompile(regexp.QuoteMeta(TemplateToken))

// SerialNumberRegexp - regular expression matching date-based DNS serial nums
var SerialNumberRegexp = regexp.MustCompile(`\b20\d\d[01]\d[0123]\d{3}\b`)

// TimeSerial returns the timeserial string for the given time.Time
func TimeSerial(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d%02d", t.Year(), t.Month(), t.Day(),
		min(((t.Hour()*3600)+(t.Minute()*60)+t.Second())/864, 99))
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

// min returns the lesser of the two arguments
func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
