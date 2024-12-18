package videoeffects

import (
	"fmt"
)

func ChangeVideoFormat(inputFileData []byte, inputFileName string, outputDirectory string, progressChannel chan uint8, outputFormat string) (string, error) {
	var outputFile = fmt.Sprintf("%sFormatChange-%s", outputDirectory, inputFileName)
	return ChangeVideoOutputFormat(inputFileData, outputFile, progressChannel, outputFormat)
}

func ChangeVideoMotionSpeed(inputFileData []byte, inputFileName string, outputDirectory string, progressChannel chan uint8, motionSpeed float32) (string, error) {
	var outputFile = fmt.Sprintf("%s%.2fSpeed-%s", outputDirectory, motionSpeed, inputFileName)
	return ChangeVideoMotion(inputFileData, outputFile, progressChannel, motionSpeed)
}
