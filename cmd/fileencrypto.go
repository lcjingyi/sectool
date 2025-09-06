package main

import (
	"fmt"

	"github.com/jingyi/sectool/pkg/dataProcessModule/fileProcessModule"
	"github.com/spf13/cobra"
)

var p fileProcessModule.FilePath

var (
	fileEncrypto bool
	fileDecrypto bool
)

var fileEncryptoCmd = &cobra.Command{
	Use:   "fileEncrypto",
	Short: "文件加解密",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO:加密
		if fileEncrypto {
			err := fileProcessModule.EncryptoFile(p)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("加密成功")
		} else if fileDecrypto {
			//TODO:解密
			err := fileProcessModule.DecryptoFile(p.Inputpath, p.KeyFilePath)
			if err != nil {
				fmt.Errorf("解密命令错误")
			}
		}
	},
}

func init() {
	fileEncryptoCmd.Flags().BoolVarP(&fileEncrypto, "Encrypto", "e", false, "加密文件")
	fileEncryptoCmd.Flags().BoolVarP(&fileDecrypto, "Decrypto", "d", false, "解密文件")
	fileEncryptoCmd.Flags().StringVarP(&p.Inputpath, "inputpath", "", "", "输入文件路径")
	fileEncryptoCmd.Flags().StringVarP(&p.Outputpath, "outputpath", "", "", "输出文件路径")
	fileEncryptoCmd.Flags().StringVarP(&p.KeyFilePath.Path, "keypath", "", "", "密钥文件路径")
}
