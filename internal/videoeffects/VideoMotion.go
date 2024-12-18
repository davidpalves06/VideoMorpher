package videoeffects

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func ChangeVideoMotion(inputFileData []byte, outputFile string, progressChannel chan uint8, motionSpeed float32) (string, error) {
	var filter string

	if motionSpeed < 0.25 || motionSpeed > 2 {
		log.Println("Motion speed is outside the accepted speed range!")
		return "", errors.New("motion speed is outside the accepted speed range")
	}

	if motionSpeed >= 0.5 {
		var videoFilterSpeed = 1 / motionSpeed
		var audioFilterSpeed = motionSpeed
		log.Println("Video filter speed:", videoFilterSpeed, "; Audio filter speed:", audioFilterSpeed)
		filter = fmt.Sprintf("[0:v]setpts=%.2f*PTS[v];[0:a]atempo=%.2f[a]", videoFilterSpeed, audioFilterSpeed)
	} else {
		//TODO: TAKE CARE OF THIS MOTION SPEED
		log.Println(motionSpeed)
	}

	go startFFmpegMotionChange(inputFileData, outputFile, progressChannel, motionSpeed, filter)
	return outputFile, nil
}

func startFFmpegMotionChange(inputFileData []byte, outputFile string, progressChannel chan uint8, motionSpeed float32, filter string) {
	cmd := exec.Command("ffmpeg", "-loglevel", "info", "-progress", "pipe:1", "-i", "pipe:0", "-filter_complex", filter, "-map", "[v]", "-map", "[a]", "-y", "-preset", "veryfast", "-c:v", "libx264", outputFile)
	cmd.Stdin = bytes.NewReader(inputFileData)
	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()

	err := cmd.Start()

	if err != nil {
		log.Fatal("OMG")
	}

	var outputVideoDuration int64 = int64(float32(GetInputVideoDuration(stderrPipe)) / motionSpeed)

	go SendProgressPercentageThroughChannel(stdoutPipe, outputVideoDuration, progressChannel)

	err = cmd.Wait()
	if err != nil {
		log.Print("ERROR on FFmpeg")
		log.Fatal(err)
	}

	close(progressChannel)
	log.Printf("Output file %v generated\n", outputFile)
}
