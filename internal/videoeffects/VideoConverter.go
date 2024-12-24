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
	ext := filepath.Ext(outputFile)[1:]
	filenameWithExt := outputFile[:len(outputFile)-len(ext)] + outputFormat

	go startFFmpegConversion(inputFileData, ext, progressChannel, outputFormat, filenameWithExt)
	return filenameWithExt, nil
}

func startFFmpegConversion(inputFileData []byte, inputFormat string, progressChannel chan uint8, outputFormat string, outputFile string) {
	var cmd *exec.Cmd
	//TODO: CHECK ALL FORMATS BEST SETTINGS
	if inputFormat == "mp4" && outputFormat == "webm" {
		cmd = exec.Command("ffmpeg", "-progress", "pipe:1", "-i", "pipe:0", "-y", "-crf", "30", "-c:v", "libvpx-vp9", "-b:v", "0",
			"-c:a", "libopus", "-b:a", "128k", "-cpu-used", "5", "-deadline", "realtime", "-f", outputFormat, outputFile)
	} else if inputFormat == "mp4" && outputFormat == "avi" {
		cmd = exec.Command("ffmpeg", "-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "mpeg4", "-q:v", "5",
			"-c:a", "mp3", "-b:a", "192k", "-f", outputFormat, outputFile)
	} else if inputFormat == "mp4" && outputFormat == "mkv" {
		cmd = exec.Command("ffmpeg", "-progress", "pipe:1", "-i", "pipe:0", "-y", "-map", "0", "-c", "copy", "-map_metadata", "0",
			"-f", "matroska", outputFile)
	} else if inputFormat == "mp4" && outputFormat == "mov" {
		// TODO: CHECK IF CAN SHOWCASE THIS ON PLAYER
		cmd = exec.Command("ffmpeg", "-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-crf", "18", "-c:a", "aac",
			"-b:a", "192k", "-f", outputFormat, outputFile)
	} else {
		cmd = exec.Command("ffmpeg", "-progress", "pipe:1", "-i", "pipe:0", "-y", "-preset", "veryfast", "-crf", "20", "-f", outputFormat, outputFile)
	}
	log.Println(cmd.Args)
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

	log.Printf("Output file %s generated\n", outputFile)
}
