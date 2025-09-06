package main

import (
	"fmt"

	"github.com/jingyi/sectool/pkg/dataProcessModule/codecModule"
	"github.com/spf13/cobra"
)

var (
	b32Encode string
	b32Decode string
	b64Encode string
	b64Decode string
)

var BaseCmd = &cobra.Command{
	Use:   "base",
	Short: "base编码",
	Example: `
	#base64编码
	cmd base --encode-b64 hello`,
	Run: func(cmd *cobra.Command, args []string) {
		if b64Encode != "" {
			fmt.Println(codecModule.B64Encode(b64Encode))
		} else if b64Decode != "" {
			fmt.Println(codecModule.B64Decode(b64Decode))
		} else if b32Encode != "" {
			fmt.Println(codecModule.B32Eecode(b32Encode))
		} else if b32Decode != "" {
			fmt.Println(codecModule.B32Decode(b32Decode))
		}

	},
}

func init() {
	BaseCmd.Flags().StringVarP(&b64Encode, "encode-b64", "", "", "base64编码")
	BaseCmd.Flags().StringVarP(&b64Decode, "decode-b64", "", "", "base64解码")
	BaseCmd.Flags().StringVarP(&b32Encode, "encode-b32", "", "", "base32编码")
	BaseCmd.Flags().StringVarP(&b32Decode, "decode-b32", "", "", "base32解码")
}
