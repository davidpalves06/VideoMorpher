package videoeffects

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func ChangeVideoMotion(tmpFileName string, outputFile string, progressChannel chan uint8, motionSpeed float32) (string, error) {
	var filter string

	if motionSpeed < 0.25 || motionSpeed > 10 {
		logger.Warn().Println("Motion speed is outside the accepted range")
		return "", errors.New("motion speed is outside the accepted speed range")
	}

	if motionSpeed >= 0.5 {
		var videoFilterSpeed = 1 / motionSpeed
		var audioFilterSpeed = motionSpeed
		filter = fmt.Sprintf("[0:v]setpts=%.2f*PTS[v];[0:a]atempo=%.2f[a]", videoFilterSpeed, audioFilterSpeed)
	} else {
		var videoFilterSpeed = 1 / motionSpeed
		var audioFilterSpeed = motionSpeed / 0.5
		filter = fmt.Sprintf("[0:v]setpts=%.2f*PTS[v];[0:a]atempo=0.5,atempo=%.2f[a]", videoFilterSpeed, audioFilterSpeed)
	}

	logger.Debug().Printf("FFmpeg filter: %s\n", filter)

	go startFFmpegMotionChange(tmpFileName, outputFile, progressChannel, motionSpeed, filter)
	return outputFile, nil
}

func startFFmpegMotionChange(tmpFileName string, outputFile string, progressChannel chan uint8, motionSpeed float32, filter string) {
	inputFormat := filepath.Ext(outputFile)[1:]
	var cmd *exec.Cmd
	params := getMotionCommandParameters(tmpFileName, inputFormat, filter, outputFile)
	if len(params) != 0 {
		cmd = exec.Command("ffmpeg", params...)
	} else {
		cmd = exec.Command("ffmpeg", "-loglevel", "info", "-progress", "pipe:1", "-i", "pipe:0", "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y", "-preset", "veryfast", outputFile)
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
	outputVideoDuration = int64(float32(float32(outputVideoDuration) / motionSpeed))
	logger.Debug().Printf("Output video duration : %d\n", outputVideoDuration)

	go sendProgressPercentageThroughChannel(stdoutPipe, outputVideoDuration, progressChannel)

	err = cmd.Wait()
	if err != nil {
		logger.Error().Println("Error while executing ffmpeg command")
		return
	}

	log.Printf("Output file %v generated\n", outputFile)
	removeTempFile(tmpFileName)
}

func getMotionCommandParameters(tmpFileName string, inputFormat string, filter string, outputFile string) []string {
	if inputFormat == "webm" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "libvpx-vp9", "-b:v", "0", "-c:a", "libopus",
			"-b:a", "128k", "-cpu-used", "5", "-deadline", "realtime",
			"-crf", "23", outputFile}
	} else if inputFormat == "mp4" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "libx264", "-preset", "veryfast",
			"-crf", "23", "-c:a", "aac", "-b:a", "128k", outputFile}
	} else if inputFormat == "avi" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "mpeg4", "-q:v", "5",
			"-c:a", "mp3", "-b:a", "192k", outputFile}
	} else if inputFormat == "mkv" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-map_metadata", "0", outputFile}
	} else if inputFormat == "mov" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "libx264", "-crf", "23", "-c:a", "aac",
			"-b:a", "192k", outputFile}
	} else if inputFormat == "ogv" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "libtheora", "-q:v", "7", "-c:a", "libvorbis",
			"-q:a", "5", "-f", inputFormat, outputFile}
	} else if inputFormat == "flv" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "libx264", "-preset", "veryfast",
			"-crf", "23", "-c:a", "aac", "-b:a", "128k", "-f", inputFormat, outputFile}
	} else if inputFormat == "mpeg" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "mpeg2video", "-b:v", "5000k",
			"-c:a", "ac3", "-b:a", "384k", "-f", inputFormat, outputFile}
	} else if inputFormat == "nut" {
		return []string{"-progress", "pipe:1", "-i", tmpFileName, "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y",
			"-c:v", "libx264", "-crf", "23",
			"-c:a", "aac", "-b:a", "128k", "-c:s", "copy", "-f", inputFormat, outputFile}
	}
	return []string{}
}
