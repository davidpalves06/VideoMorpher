package videoeffects

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"path/filepath"
)

func ChangeVideoOutputFormat(inputFileData []byte, outputFile string, progressChannel chan uint8, outputFormat string) (string, error) {

	if !is_Format_Supported(outputFormat) {
		return "", errors.New("output format is not supported")
	}
	ext := filepath.Ext(outputFile)
	filenameWithExt := outputFile[:len(outputFile)-len(ext)] + "." + outputFormat

	go startFFmpegConversion(inputFileData, progressChannel, outputFormat, filenameWithExt)
	return filenameWithExt, nil
}

func startFFmpegConversion(inputFileData []byte, progressChannel chan uint8, outputFormat string, outputFile string) {
	//TODO: CHECK BEST CONVERSION PARAMS FROM FORMAT TO FORMAT
	cmd := exec.Command("ffmpeg", "-progress", "pipe:1", "-i", "pipe:0", "-y", "-preset", "veryfast", "-c:v", "libx264", "-c:a", "copy", "-f", outputFormat, outputFile)
	cmd.Stdin = bytes.NewReader(inputFileData)
	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()
	err := cmd.Start()

	if err != nil {
		log.Fatal("OMG")
	}

	var outputVideoDuration int64 = GetInputVideoDuration(stderrPipe)

	go SendProgressPercentageThroughChannel(stdoutPipe, outputVideoDuration, progressChannel)

	err = cmd.Wait()
	if err != nil {
		log.Print("ERROR on FFmpeg")
		log.Fatal(err)
	}

	close(progressChannel)
	log.Printf("Output file %s generated\n", outputFile)
}
