package videoeffects

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func ChangeVideoOutputFormat(tmpFileName string, inputFileName string, progressChannel chan uint8, outputFormat string) (string, error) {
	if !is_Format_Supported(outputFormat) {
		logger.Warn().Printf("Output format %s is not supported\n", outputFormat)
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

	logger.Debug().Println(cmd.Args)
	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()

	logger.Debug().Println("Starting ffmpeg command")
	err := cmd.Start()

	if err != nil {
		logger.Error().Printf("Error starting ffmpeg command\n")
		return
	}

	var outputVideoDuration int64 = getInputVideoDuration(stderrPipe)
	logger.Debug().Printf("Output video duration : %d\n", outputVideoDuration)

	go sendProgressPercentageThroughChannel(stdoutPipe, outputVideoDuration, progressChannel)

	err = cmd.Wait()
	if err != nil {
		logger.Error().Println("Error while executing ffmpeg command")
		return
	}

	logger.Debug().Printf("Output file %s generated\n", outputFile)
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
