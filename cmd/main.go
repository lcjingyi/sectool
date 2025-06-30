package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "myApp", // 命令名
		Short: "Small tools for encryption, decryption, encoding and decoding",
	}

	rootCmd.AddCommand(UrlCmd, BaseCmd, DirScanCmd, PingCmd)
	rootCmd.ExecuteC()
}
