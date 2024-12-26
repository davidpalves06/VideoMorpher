package videoeffects

import (
	"fmt"
)

func ChangeVideoFormat(tempFileName string, inputFileName string, outputDirectory string, progressChannel chan uint8, outputFormat string) (string, error) {
	var outputFile = fmt.Sprintf("%sFormatChange-%s", outputDirectory, inputFileName)
	return ChangeVideoOutputFormat(tempFileName, outputFile, progressChannel, outputFormat)
}

func ChangeVideoMotionSpeed(tempFileName string, inputFileName string, outputDirectory string, progressChannel chan uint8, motionSpeed float32) (string, error) {
	var outputFile = fmt.Sprintf("%s%.2fSpeed-%s", outputDirectory, motionSpeed, inputFileName)
	return ChangeVideoMotion(tempFileName, outputFile, progressChannel, motionSpeed)
}
