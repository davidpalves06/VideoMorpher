package videoeffects

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

func ChangeVideoMotion(inputFileData []byte, outputFile string, progressChannel chan uint8, motionSpeed float32) (string, error) {
	var filter string

	if motionSpeed < 0.25 || motionSpeed > 10 {
		log.Println("Motion speed is outside the accepted speed range!")
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

	go startFFmpegMotionChange(inputFileData, outputFile, progressChannel, motionSpeed, filter)
	return outputFile, nil
}

func startFFmpegMotionChange(inputFileData []byte, outputFile string, progressChannel chan uint8, motionSpeed float32, filter string) {
	inputFormat := filepath.Ext(outputFile)[1:]
	//TODO: CHECK BEST SETTINGS PER FORMAT
	var cmd *exec.Cmd
	params := getMotionCommandParameters(inputFormat, outputFile)
	if len(params) != 0 {
		cmd = exec.Command("ffmpeg", params...)
	} else {
		cmd = exec.Command("ffmpeg", "-loglevel", "info", "-progress", "pipe:1", "-i", "pipe:0", "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y", "-preset", "veryfast", outputFile)
	}
	cmd.Stdin = bytes.NewReader(inputFileData)
	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()

	err := cmd.Start()

	if err != nil {
		log.Fatal("OMG")
	}

	var outputVideoDuration int64 = -1

	if inputFormat != "ogv" {
		outputVideoDuration = getInputVideoDuration(stderrPipe)
	} else {
		outputVideoDuration, _ = getOggDurationMs(inputFileData)
	}
	outputVideoDuration = int64(float32(float32(outputVideoDuration) / motionSpeed))
	go sendProgressPercentageThroughChannel(stdoutPipe, outputVideoDuration, progressChannel)

	err = cmd.Wait()
	if err != nil {
		log.Print("ERROR on FFmpeg")
		log.Fatal(err)
	}

	log.Printf("Output file %v generated\n", outputFile)
}

func getMotionCommandParameters(outputFormat string, outputFile string) []string {
	if outputFormat == "webm" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-crf", "30", "-c:v", "libvpx-vp9", "-b:v", "0",
			"-c:a", "libopus", "-b:a", "128k", "-cpu-used", "5", "-deadline", "realtime", "-f", outputFormat, outputFile}
	} else if outputFormat == "mp4" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-preset", "veryfast",
			"-crf", "23", "-c:a", "aac", "-b:a", "128k", "-f", outputFormat, outputFile}
	} else if outputFormat == "avi" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "mpeg4", "-q:v", "5",
			"-c:a", "mp3", "-b:a", "192k", "-f", outputFormat, outputFile}
	} else if outputFormat == "mkv" {
		// TODO: CHECK IF CAN SHOWCASE THIS ON PLAYER
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-map", "0", "-c", "copy", "-map_metadata", "0",
			"-f", "matroska", outputFile}
	} else if outputFormat == "mov" {
		// TODO: CHECK IF CAN SHOWCASE THIS ON PLAYER
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-crf", "23", "-c:a", "aac",
			"-b:a", "192k", "-f", outputFormat, outputFile}
	} else if outputFormat == "ogv" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libtheora", "-q:v", "7", "-c:a", "libvorbis",
			"-q:a", "5", "-f", outputFormat, outputFile}
	} else if outputFormat == "flv" {
		// TODO: CHECK IF CAN SHOWCASE THIS ON PLAYER
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-preset", "veryfast",
			"-crf", "23", "-c:a", "aac", "-b:a", "128k", "-f", outputFormat, outputFile}
	} else if outputFormat == "mpeg" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "mpeg2video", "-b:v", "5000k",
			"-c:a", "ac3", "-b:a", "384k", "-f", outputFormat, outputFile}
	} else if outputFormat == "nut" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-crf", "23",
			"-c:a", "aac", "-b:a", "128k", "-c:s", "copy", "-f", outputFormat, outputFile}
	}
	return []string{}
}
