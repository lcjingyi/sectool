package main

import (
	"fmt"

	"github.com/jingyi/sectool/pkg/informationCollectionModule"
	"github.com/spf13/cobra"
)

var (
	ipHost    string
	portStart int
	portEnd   int
)

var portScanfCmd = &cobra.Command{
	Use:   "portScanf",
	Short: "端口扫描",
	Run: func(cmd *cobra.Command, args []string) {
		if ipHost != "" && portStart != 0 && portEnd != 0 && portStart < portEnd {
			informationCollectionModule.PortScan(ipHost, portStart, portEnd)
		} else {
			fmt.Errorf("输入错误")
		}
	},
}

func init() {
	portScanfCmd.Flags().StringVarP(&ipHost, "ip", "i", "", "The ip address of the host to scan")
	portScanfCmd.Flags().IntVarP(&portStart, "start", "s", 0, "The start port of the range to scan")
	portScanfCmd.Flags().IntVarP(&portEnd, "end", "e", 0, "The end port of the range to scan")
}
