package main

import (
	"bufio"
	"log"
	"os"

	"github.com/jingyi/sectool/pkg/networkCommunicationModule"
	"github.com/spf13/cobra"
)

var (
	ip     string
	ipFile string
)

var PingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping",
	Run: func(cmd *cobra.Command, args []string) {
		if ip != "" {
			for i := 0; i < 4; i++ {
				err := networkCommunicationModule.IcmpPing(ip)
				if err != nil {
					return
				}
			}

		} else if ipFile != "" {
			f, err := os.Open(ipFile)
			if err != nil {
				return
			}

			// 关闭文件
			logFile, err := os.OpenFile("ping.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return
			}
			defer func(logFile *os.File) {
				_ = logFile.Close()
			}(logFile)
			log.SetOutput(logFile)

			defer func(f *os.File) {
				err := f.Close()
				if err != nil {

				}
			}(f)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				err := networkCommunicationModule.IcmpPing(line)
				if err != nil {
					return
				}
			}
		}
	},
}

func init() {
	PingCmd.Flags().StringVarP(&ip, "ip", "", "", "目标ip")
	PingCmd.Flags().StringVarP(&ipFile, "ip-file", "", "", "ip文件")
}
