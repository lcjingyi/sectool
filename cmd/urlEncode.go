package main

import (
	"fmt"

	"github.com/jingyi/sectool/pkg"
	"github.com/spf13/cobra"
)

var (
	urlEncodeStr string
	urlDecodeStr string
)

var UrlCmd = &cobra.Command{
	Use:   "url",
	Short: "url编码",
	Run: func(cmd *cobra.Command, args []string) {
		if urlEncodeStr != "" {
			fmt.Println(pkg.MyUrlEncode(urlEncodeStr))
		} else if urlDecodeStr != "" {
			fmt.Println(pkg.MyUrlDncode(urlDecodeStr))
		} else {
			fmt.Println("error")
		}
	},
}

func init() {
	UrlCmd.Flags().StringVarP(&urlEncodeStr, "encode", "", "", "url编码")
	UrlCmd.Flags().StringVarP(&urlDecodeStr, "decode", "", "", "url解码")
}
