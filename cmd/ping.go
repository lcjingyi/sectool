package main

import (
	"github.com/jingyi/sectool/pkg"
	"github.com/spf13/cobra"
)

var ip string

var PingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping",
	Run: func(cmd *cobra.Command, args []string) {
		if ip != "" {
			pkg.IcmpPing(ip)
		}
	},
}

func init() {
	PingCmd.Flags().StringVarP(&ip, "ip", "", "", "目标ip")
}
