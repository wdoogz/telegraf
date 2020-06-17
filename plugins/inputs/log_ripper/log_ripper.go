package log_ripper

import (
	"io/ioutil"
	"log"
	"regexp"

	"github.com/influxdata/telegraf"
)

// RipperStruct is the base struct for the config file
type RipperStruct struct {
	LogFiles      []string `toml:log_files`
	RegexOverride []string `toml:regex`
}

// Description gives a description
func (_ *RipperStruct) Description() string {
	return "Parses log files for Faults"
}

var ripperSampleConfig = `
	## By default log_ripper will parse logs for errors
	## You may bring your own regex to scrape logs for what means most to you
	# log_files = ["/var/log/messages"]
	# regex = ["(Error|Failure)"]
`

// SampleConfig generates a sample template
func (_ *RipperStruct) SampleConfig() string {
	return ripperSampleConfig
}

var baseRegex = "[eE][rR]{2}[oO][rR]"

// LogFileStats does this
func (rs *RipperStruct) LogFileStats(acc telegraf.Accumulator) error {
	var totalError int = 0
	for _, lFile := range rs.LogFiles {
		totalError = parseLogFile(lFile)
		tags := map[string]string{
			"FilePath": lFile,
		}
		fields := map[string]interface{}{
			"total_errors": totalError,
		}
		acc.AddFields("logErrors", fields, tags)
	}
	return nil
}

func parseLogFile(filename string) int {
	logFile, _ := ioutil.ReadFile(filename)
	r, err := regexp.Compile("[eE][rR]{2}[oO][rR]")
	if err != nil {
		log.Fatal(err)
	}
	results := r.FindAllString(string(logFile), -1)
	return len(results)
}
