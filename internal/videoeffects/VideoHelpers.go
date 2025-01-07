package videoeffects

import (
	"bufio"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

var ActiveCommands map[string]exec.Cmd = make(map[string]exec.Cmd)
var VideoCommandWaitGroup sync.WaitGroup

func is_Format_Supported(format string) bool {
	var supportedFormats = map[string]bool{
		"mp4": true, "avi": true, "mkv": true, "mov": true, "webm": true, "ogv": true, "flv": true, "mpeg": true, "nut": true,
	}
	return supportedFormats[format]
}

func parseDuration(parts []string) int64 {
	logger.Debug().Printf("Parsing duration for parts: %v\n", parts)
	hours, _ := time.ParseDuration(parts[0] + "h")
	minutes, _ := time.ParseDuration(parts[1] + "m")

	secParts := strings.Split(parts[2], ".")
	seconds, _ := time.ParseDuration(secParts[0] + "s")
	var milliseconds time.Duration
	if len(secParts) > 1 {
		milliseconds, _ = time.ParseDuration(secParts[1] + "ms")
	}

	duration := hours + minutes + seconds + milliseconds
	logger.Debug().Printf("Parsed duration : %d\n", duration)
	return duration.Microseconds()
}

func getInputVideoDuration(stderrPipe io.ReadCloser) int64 {
	logger.Debug().Println("Reading stderr from ffmpeg to get input video duration")
	scanner := bufio.NewScanner(stderrPipe)
	var inputVideoDuration int64 = -1
	for scanner.Scan() {
		var cmdOutput = strings.TrimSpace(strings.ToLower(scanner.Text()))
		if strings.Contains(cmdOutput, "duration") {
			duration := strings.Split(strings.Split(cmdOutput, ",")[0], ":")[1:]
			if len(duration) < 3 {
				continue
			}
			parsedDuration := parseDuration(duration)
			inputVideoDuration = parsedDuration
			break
		}
	}
	logger.Debug().Printf("Input video duration: %d\n", inputVideoDuration)
	return inputVideoDuration
}

func sendProgressPercentageThroughChannel(stdoutPipe io.ReadCloser, outputVideoDuration int64, progressChannel chan uint8) {
	var progressPercentage uint8 = 0
	logger.Debug().Println("Reading stdout from ffmpeg to get processing progress")
	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		var cmdOutput = strings.TrimSpace(scanner.Text())
		if strings.Contains(cmdOutput, "out_time_ms") {
			us_Output_time, err := strconv.ParseInt(strings.Split(cmdOutput, "=")[1], 10, 64)
			if err != nil {
				logger.Warn().Printf("Error while parsing progression for log %s. Error : %s\n", cmdOutput, err.Error())
			} else {
				progressPercentage = uint8(math.Round(float64(us_Output_time) / float64(outputVideoDuration) * 100))
				logger.Debug().Printf("Sending Process Percentage through channel: %d\n", progressPercentage)
				progressChannel <- progressPercentage
			}
		}
	}
	logger.Debug().Println("Closing progress channel")
}

func removeTempFile(tmpFileName string) {
	logger.Debug().Printf("Removing temporary file %s\n", tmpFileName)
	err := os.Remove(tmpFileName)
	if err != nil {
		logger.Warn().Printf("Could not remove temporary file %s : %s\n", tmpFileName, err.Error())
	} else {
		logger.Debug().Printf("Temporary file %s sucessfully removed\n", tmpFileName)
	}
}
