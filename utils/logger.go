package utils

import (
	"encoding/json"
	"log"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

// init is a function that is automatically executed before the main function is called.
func init() {
	makeLogFile()
	setUpLogger()
}

// makeLogFile creates a log file in the specified path.
func makeLogFile() {
	_, err := os.Create("./logs/app.log")
	if err != nil {
		log.Fatalf("failed to create the logfile: %v", err.Error())
	}
}

// setUpLogger sets up the logger configuration using the specified JSON format.
func setUpLogger() {
	var cfg zap.Config

	formatJson := []byte(
		`
		{
			"level": "debug",
			"encoding": "json",
			"outputPaths": [
				"stdout",
				"./logs/app.log"
			],
			"errorOutputPaths": [
				"stderr"
			],
			"encoderConfig": {
				"messageKey": "message",
				"levelKey": "level",
				"timeKey": "time",
				"nameKey": "logger",
				"callerKey": "caller",
				"stacktraceKey": "stacktrace",
				"lineEnding": "\n",
				"timeEncoder": "iso8601",
				"levelEncoder": "lowercase",
				"callerEncoder": "short"
			}
		}
		`)

	if err := json.Unmarshal(formatJson, &cfg); err != nil {
		panic(err)
	}

	Logger = zap.Must(cfg.Build()).Sugar()
	defer Logger.Sync()
}