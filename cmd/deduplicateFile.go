package main

import (
	"github.com/jingyi/sectool/pkg/dataProcessModule/fileProcessModule"
	"github.com/spf13/cobra"
)

//实现文件去重功能命令行

// 需要输入两个参数：
// 1. 待去重文件路径
// 2. 去重后文件路径
// 命令行参数：
// ./deduplicateFile -f /path/to/file1 -t /path/to/file2
var (
	inputFilePath  string
	outputFilePath string
)

var deduplicateFileCmd = &cobra.Command{
	Use:   "deduplicateFile",
	Short: "文件去重功能",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: 实现文件去重功能
		if inputFilePath != "" && outputFilePath != "" {
			//TODO:调用文件去重函数
			_ = fileProcessModule.DeduplicateFile(inputFilePath, outputFilePath)
		}
	},
}

func init() {
	deduplicateFileCmd.Flags().StringVarP(&inputFilePath, "file", "f", "", "待去重文件路径")
	deduplicateFileCmd.Flags().StringVarP(&outputFilePath, "target", "t", "", "去重后文件路径")
}
