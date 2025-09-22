package main

import (
	"fmt"

	"github.com/jingyi/sectool/pkg/informationCollectionModule"
	"github.com/spf13/cobra"
)

var (
	url string
)

var DirScanCmd = &cobra.Command{
	Use:   "dirScan",
	Short: "dir扫描",
	Run: func(cmd *cobra.Command, args []string) {
		if url != "" {
			err := informationCollectionModule.DirScan(url)
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	DirScanCmd.Flags().StringVarP(&url, "u", "", "", "目标地址")
}
