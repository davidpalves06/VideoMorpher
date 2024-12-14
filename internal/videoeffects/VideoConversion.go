package videoeffects

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os/exec"
)

func VideoConversion(inputFileData []byte, outputFile string) {

	cmd := exec.Command("ffmpeg", "-report", "-i", "pipe:0", "-filter_complex", `[0:v]setpts=0.5*PTS[v];[0:a]atempo=2.0[a]`, "-map", "[v]", "-map", "[a]", "-y", "-preset", "veryfast", outputFile)
	log.Println(cmd.Args)
	cmd.Stdin = bytes.NewReader(inputFileData)
	stderrPipe, _ := cmd.StderrPipe()

	err := cmd.Start()

	if err != nil {
		log.Fatal("OMG")
	}

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Printf("STDERR: %s\n", scanner.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Print("ERROR on FFmpeg")
		log.Fatal(err)
	}

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
