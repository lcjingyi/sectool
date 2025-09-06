package fileProcessModule

import (
	"bufio"
	"os"
)

// DeduplicateFile 文件去重模块
func DeduplicateFile(inputPath, outputPath string) error {
	// 打开输入文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func(inputFile *os.File) {
		_ = inputFile.Close()
	}(inputFile)

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	seen := make(map[string]bool)
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	//刷新输出缓冲区确保输出完整
	defer func(writer *bufio.Writer) {
		_ = writer.Flush()
	}(writer)

	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			if _, err := writer.WriteString(line + "\n"); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
