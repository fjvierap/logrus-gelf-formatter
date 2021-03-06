package formatter

// Portions of this files were taken from https://github.com/sirupsen/logrus
// Copyrighted by Simon Eskildsen under the MIT license

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Defines a log format type that wil output line separated JSON objects
// in the GELF format.
type GelfTimestampFormatter struct{}

// Format formats the log entry to GELF JSON
func (f *GelfTimestampFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(fields, len(entry.Data)+6)
	blacklist := []string{"_id", "id", "timestamp", "version", "level"}
	var timestamp float64

	for k, v := range entry.Data {

		if contains(k, blacklist) {
			continue
		}

		if k == "_timestamp" {
			timestamp = v.(float64)
			continue
		}

		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			data["_"+k] = v.Error()
		default:
			data["_"+k] = v
		}
	}

	data["version"] = "1.1"
	data["short_message"] = entry.Message
	if timestamp != 0 {
		data["timestamp"] = timestamp
	} else {
		data["timestamp"] = round((float64(entry.Time.UnixNano())/float64(1000000))/float64(1000), 4)
	}
	data["level"] = entry.Level
	data["level_name"] = entry.Level.String()
	data["_pid"] = os.Getpid()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}

	return append(serialized, '\n'), nil
}
