package videoeffects

import (
	"errors"
	"log"
	"os/exec"
	"path/filepath"
)

func ChangeVideoOutputFormat(tmpFileName string, inputFileName string, progressChannel chan uint8, outputFormat string) (string, error) {

	if !is_Format_Supported(outputFormat) {
		return "", errors.New("output format is not supported")
	}
	ext := filepath.Ext(inputFileName)[1:]
	filenameWithExt := inputFileName[:len(inputFileName)-len(ext)] + outputFormat

	go startFFmpegConversion(tmpFileName, progressChannel, outputFormat, filenameWithExt)
	return filenameWithExt, nil
}

func startFFmpegConversion(tmpFileName string, progressChannel chan uint8, outputFormat string, outputFile string) {
	var cmd *exec.Cmd
	params := getConversionCommandParameters(tmpFileName, outputFormat, outputFile)
	if len(params) != 0 {
		cmd = exec.Command("ffmpeg", params...)
	} else {
		cmd = exec.Command("ffmpeg", "-progress", "pipe:1", "-i", tmpFileName, "-y", "-c", "copy", "-f", outputFormat, outputFile)
	}
	log.Println(cmd.Args)
	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()
	err := cmd.Start()

	if err != nil {
		log.Fatal("OMG")
	}

	var outputVideoDuration int64 = getInputVideoDuration(stderrPipe)

	go sendProgressPercentageThroughChannel(stdoutPipe, outputVideoDuration, progressChannel)

	err = cmd.Wait()
	if err != nil {
		log.Print("ERROR on FFmpeg")
		log.Fatal(err)
	}

	log.Printf("Output file %s generated\n", outputFile)
	removeTempFile(tmpFileName)
}

func getConversionCommandParameters(tmpFileName string, outputFormat string, outputFile string) []string {
	if outputFormat == "webm" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-crf", "30", "-c:v", "libvpx-vp9", "-b:v", "0",
			"-c:a", "libopus", "-b:a", "128k", "-cpu-used", "5", "-deadline", "realtime", "-f", outputFormat, outputFile}
	} else if outputFormat == "mp4" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-c:v", "libx264", "-preset", "veryfast",
			"-crf", "23", "-c:a", "aac", "-b:a", "128k", "-f", outputFormat, outputFile}
	} else if outputFormat == "avi" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-c:v", "mpeg4", "-q:v", "5",
			"-c:a", "mp3", "-b:a", "192k", "-f", outputFormat, outputFile}
	} else if outputFormat == "mkv" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-map", "0", "-c", "copy", "-map_metadata", "0",
			"-f", "matroska", outputFile}
	} else if outputFormat == "mov" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-c:v", "libx264", "-crf", "23", "-c:a", "aac",
			"-b:a", "192k", "-f", outputFormat, outputFile}
	} else if outputFormat == "ogv" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-c:v", "libtheora", "-q:v", "7", "-c:a", "libvorbis",
			"-q:a", "5", "-f", outputFormat, outputFile}
	} else if outputFormat == "flv" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-c:v", "libx264", "-preset", "veryfast",
			"-crf", "23", "-c:a", "aac", "-b:a", "128k", "-f", outputFormat, outputFile}
	} else if outputFormat == "mpeg" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-c:v", "mpeg2video", "-b:v", "5000k",
			"-c:a", "ac3", "-b:a", "384k", "-f", outputFormat, outputFile}
	} else if outputFormat == "nut" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-y", "-c:v", "libx264", "-crf", "23",
			"-c:a", "aac", "-b:a", "128k", "-c:s", "copy", "-f", outputFormat, outputFile}
	}
	return []string{}
}
