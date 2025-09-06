package main

import (
	"github.com/jingyi/sectool/pkg/informationCollectionModule"
	"github.com/spf13/cobra"
)

var (
	url  string
	path string
)

var DirScanCmd = &cobra.Command{
	Use:   "dirScan",
	Short: "dir扫描",
	Run: func(cmd *cobra.Command, args []string) {
		if url != "" && path != "" {
			informationCollectionModule.DirBlasting(url, path)
		} else if url != "" && path == "" {
			informationCollectionModule.DirBlasting(url, path)
		}
	},
}

func init() {
	DirScanCmd.Flags().StringVarP(&url, "u", "", "", "目标地址")
	DirScanCmd.Flags().StringVarP(&path, "p", "", "", "字典地址")
}
