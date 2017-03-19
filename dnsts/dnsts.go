package dnsts

import (
	"flag"
	"fmt"
	"github.com/actualben/timeserial"
	"os"
	"time"
)

const helpText string = `usage: %s [-h] [-u]

Prints a 10-digit number suitable for use as a DNS zone serial number. The
number is based on the current date and time in this format: YYYYMMDDTT. Here
YYYY represents the year including century, MM represents the zero-padded
month number, DD represents the zero-padded day of the month, and TT
represents the time in hundredths-of-a-day (864 second / ~15 minute intervals)
since midnight.

`

func main() {
	useUTC := false
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helpText, os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVar(&useUTC, "u", useUTC, "Use UTC instead of local time")
	flag.Parse()

	t := time.Now()
	if useUTC {
		t = t.UTC()
	}
	fmt.Println(timeserial.TimeSerial(t))
}
