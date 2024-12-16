package videoeffects

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func VideoConversion(inputFileData []byte, outputFile string, progressChannel chan uint8) {

	cmd := exec.Command("ffmpeg", "-loglevel", "info", "-progress", "pipe:1", "-i", "pipe:0", "-filter_complex", `[0:v]setpts=0.5*PTS[v];[0:a]atempo=2.0[a]`, "-map", "[v]", "-map", "[a]", "-y", "-preset", "veryfast", outputFile)
	cmd.Stdin = bytes.NewReader(inputFileData)
	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()

	err := cmd.Start()

	if err != nil {
		log.Fatal("OMG")
	}

	var outputVideoDuration int64 = -1
	var progressPercentage uint8 = 0
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			var cmdOutput = strings.TrimSpace(scanner.Text())
			if strings.Contains(cmdOutput, "Duration") {
				duration, _ := strings.CutPrefix(strings.Split(cmdOutput, ",")[0], "Duration:")
				parsedDuration := parseDuration(duration)
				outputVideoDuration = parsedDuration.Microseconds() / 2
				log.Printf("Estimated Output Video duration: %s, In Millisecond: %v\n", parsedDuration, outputVideoDuration)
				break
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			var cmdOutput = strings.TrimSpace(scanner.Text())
			if strings.Contains(cmdOutput, "out_time_us") {
				us_Output_time, _ := strconv.ParseInt(strings.Split(cmdOutput, "=")[1], 10, 64)
				progressPercentage = uint8(math.Round(float64(us_Output_time) / float64(outputVideoDuration) * 100))
				progressChannel <- progressPercentage
				log.Printf("PROGRESS: %d%%\n", progressPercentage)
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Print("ERROR on FFmpeg")
		log.Fatal(err)
	}

	close(progressChannel)
	log.Printf("Output file %v generated\n", outputFile)

}

func GetFileBytes(file io.Reader) []byte {
	var fileBytesBuffer []byte = make([]byte, 1024*1024)
	var isEoF bool = false
	var fileBytes []byte = make([]byte, 0, 1024*1024)

	for !isEoF {
		bytesRead, err := file.Read(fileBytesBuffer)
		if err == io.EOF {
			isEoF = true
		} else {
			fileBytes = append(fileBytes, fileBytesBuffer[:bytesRead]...)
		}
	}
	return fileBytes
}

func parseDuration(timeStr string) time.Duration {
	parts := strings.Split(timeStr, ":")

	hours, _ := time.ParseDuration(parts[0] + "h")
	minutes, _ := time.ParseDuration(parts[1] + "m")

	secParts := strings.Split(parts[2], ".")
	seconds, _ := time.ParseDuration(secParts[0] + "s")
	var milliseconds time.Duration
	if len(secParts) > 1 {
		milliseconds, _ = time.ParseDuration(secParts[1] + "ms")
	}

	duration := hours + minutes + seconds + milliseconds
	return duration
}
