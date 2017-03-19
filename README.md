# timeserial

**This is all sub-alpha quality. DON'T USE IT YET!**

Generate timestamps suitable for use in a DNS zone file. There are two programs in this directory:

timeserial prints the serial for the current time to STDOUT

timeserializer searches through the given files for instances of the string `$TIMESERIAL` and replaces them with the serial representing the modification date of the file they appear in. This is intended to be used as a run step in the Dockerfile of our DNS containers. 

##  Usage

### timeserial

	usage: timeserial [-h] [-u]
	
	Prints a 10-digit number suitable for use as a DNS zone serial number. The
	number is based on the current date and time in this format: YYYYMMDDTT. Here
	YYYY represents the year including century, MM represents the zero-padded
	month number, DD represents the zero-padded day of the month, and TT
	represents the time in hundredths-of-a-day (864 second / ~15 minute intervals)
	since midnight.
	
	  -u	Use the UTC time zone instead of localtime

### timeserializer

	usage: timeserializer [-h] [-e] [-i] [-n] [-u]
	
	For each file given on the command line, replace any instances of the template
	text $TIMESERIAL with a 10-digit number suitable for use as a DNS zone
	serial number. The number is based on the current date and time in this
	format: YYYYMMDDTT. Here YYYY represents the year including century, MM
	represents the zero-padded month number, DD represents the zero-padded day of
	the month, and TT represents the time in hundredths-of-a-day
	(864 second / ~15 minute intervals) since midnight.
	
	This command is designed to expand the $TIMESERIAL template variable in DNS
	zone files as they are placed into read-only container filesystems
	
	  -e	Replace existing serial numbers, instead of replacing the template
	  -i	Edit the given files in-place
	  -n	Use the currrent time, instead of file mtimes
	  -u	Use the UTC time zone instead of localtime
