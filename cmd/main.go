package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "myApp", // 命令名
		Short: "Small tools",
	}

	rootCmd.AddCommand(
		UrlCmd,
		BaseCmd,
		DirScanCmd,
		PingCmd,
		deduplicateFileCmd,
		fileEncryptoCmd)
	_, err := rootCmd.ExecuteC()
	if err != nil {
		fmt.Println(err)
	}
}
