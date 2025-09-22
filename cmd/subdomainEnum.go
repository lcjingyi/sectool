package main

import (
	"fmt"
	"os"

	"github.com/jingyi/sectool/pkg/informationCollectionModule"
	"github.com/spf13/cobra"
)

var host string

var subdomainEnum = &cobra.Command{
	Use:   "subdomain",
	Short: "子域名爆破",
	Run: func(cmd *cobra.Command, args []string) {
		if host != "" {
			err := informationCollectionModule.SubdomainEnum(host)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("请输入域名")
			os.Exit(1)
		}
	},
}

func init() {
	subdomainEnum.Flags().StringVarP(&host, "domain", "m", "", "子域名爆破")
}
