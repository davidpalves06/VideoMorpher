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
	if inputFormat == "mp4" {
		params := getMP4CommandParameters(outputFormat, outputFile)
		cmd = exec.Command("ffmpeg", params...)
	} else if inputFormat == "avi" {
		params := getAVICommandParameters(outputFormat, outputFile)
		cmd = exec.Command("ffmpeg", params...)
	} else {
		cmd = exec.Command("ffmpeg", "-progress", "pipe:1", "-i", "pipe:0", "-y", "-c", "copy", "-f", outputFormat, outputFile)
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

func getMP4CommandParameters(outputFormat string, outputFile string) []string {
	if outputFormat == "webm" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-crf", "30", "-c:v", "libvpx-vp9", "-b:v", "0",
			"-c:a", "libopus", "-b:a", "128k", "-cpu-used", "5", "-deadline", "realtime", "-f", outputFormat, outputFile}
	} else if outputFormat == "avi" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "mpeg4", "-q:v", "5",
			"-c:a", "mp3", "-b:a", "192k", "-f", outputFormat, outputFile}
	} else if outputFormat == "mkv" {
		// TODO: CHECK IF CAN SHOWCASE THIS ON PLAYER
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-map", "0", "-c", "copy", "-map_metadata", "0",
			"-f", "matroska", outputFile}
	} else if outputFormat == "mov" {
		// TODO: CHECK IF CAN SHOWCASE THIS ON PLAYER
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-crf", "18", "-c:a", "aac",
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

func getAVICommandParameters(outputFormat string, outputFile string) []string {
	if outputFormat == "webm" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-crf", "30", "-c:v", "libvpx-vp9", "-b:v", "0",
			"-c:a", "libopus", "-b:a", "128k", "-cpu-used", "5", "-deadline", "realtime", "-f", outputFormat, outputFile}
	} else if outputFormat == "mp4" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-preset", "veryfast",
			"-crf", "23", "-c:a", "aac", "-b:a", "128k", "-f", outputFormat, outputFile}
	} else if outputFormat == "mkv" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-crf", "23", "-c:a", "aac", "-b:a", "192k", "-map_metadata", "0",
			"-f", "matroska", outputFile}
	} else if outputFormat == "mov" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libx264", "-crf", "23",
			"-c:a", "aac", "-b:a", "192k", "-c:s", "mov_text", "-map_metadata", "0", "-f", outputFormat, outputFile}
	} else if outputFormat == "ogv" {
		return []string{"-progress", "pipe:1", "-i", "pipe:0", "-y", "-c:v", "libtheora", "-q:v", "7", "-c:a", "libvorbis",
			"-q:a", "5", "-f", outputFormat, outputFile}
	} else if outputFormat == "flv" {
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
