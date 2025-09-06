package fileProcessModule

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

//生成AES密钥(256)保存到文件

type KeyFilePath struct {
	Path string
}

// 生成AES密钥并保存到文件
func generateKey(keyPath KeyFilePath) error {
	if keyPath.Path == "" {
		_, err := os.Create("key")
		if err != nil {
			return fmt.Errorf("创建密钥文件失败")
		}
		keyPath.Path = "key"
	}

	key := make([]byte, 32)
	_, err := rand.Read(key) //rand.Read 函数用于从加密安全的伪随机数生成器中读取随机字节，并将这些字节存储到指定的字节切片中。
	if err != nil {
		return err
	}
	return os.WriteFile(keyPath.Path, key, 0660)
}

// 加载密钥
func loadKey(keyPath KeyFilePath) ([]byte, error) {
	if keyPath.Path == "" {
		keyPath.Path = "key"
	}
	return os.ReadFile(keyPath.Path)
}

// 计算文件SHA-256
func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

type FilePath struct {
	KeyFilePath
	Inputpath  string
	Outputpath string
}

// 文件加密
func EncryptoFile(path FilePath) error {
	key, err := loadKey(path.KeyFilePath)
	if err != nil {
		//密钥不存在生成新的密钥
		err := generateKey(KeyFilePath{
			Path: "",
		})
		key, err = loadKey(KeyFilePath{
			Path: "",
		})
		if err != nil {
			return fmt.Errorf("密钥生成失败")
		}
	}

	//打开输入文件
	if _, err := os.Stat(path.Inputpath); os.IsNotExist(err) {
		return fmt.Errorf("输入文件不存在")
	}
	inFile, err := os.OpenFile(path.Inputpath, os.O_RDONLY, 0660)
	if err != nil {
		return fmt.Errorf("打开输入文件失败")
	}
	defer inFile.Close()

	//创建输出文件
	if path.Outputpath == "" {
		outFile, err := os.Create("output.enc")
		if err != nil {
			return err
		}
		path.Outputpath = "output.enc"
		outFile.Close()
	}
	outFile, err := os.OpenFile(path.Outputpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer outFile.Close()

	//生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	//将IV保存在输出文件头部
	if _, err := outFile.Write(iv); err != nil {
		return err
	}

	//计算原始文件哈希并写入
	hash, err := calculateHash(path.Inputpath)
	if err != nil {
		return fmt.Errorf("计算文件哈希失败")
	}

	//写入哈希值和哈希长度
	hashBytes := []byte(hash)
	if _, err := outFile.Write([]byte{byte(len(hashBytes))}); err != nil {
		return err
	}
	if _, err := outFile.Write(hashBytes); err != nil {
		return err
	}

	//创建AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("创建失败")
	}

	//创建CTR模式流加密器
	stream := cipher.NewCTR(block, iv)

	//流加密
	writer := &cipher.StreamWriter{S: stream, W: outFile}
	if _, err := io.Copy(writer, inFile); err != nil {
		return fmt.Errorf("加密过程失败")
	}

	fmt.Printf("文件加密成功：%s -> %s\n", path.Inputpath, path.Outputpath)
	fmt.Printf("密钥文件：%s\n", path.KeyFilePath.Path)
	return nil
}

// 文件解密
func DecryptoFile(inputFilePath string, keyPath KeyFilePath) error {
	//加载密钥
	key, err := loadKey(keyPath)
	if len(key) != 32 {
		return fmt.Errorf("密码读取错误")
	} else if err != nil {
		return fmt.Errorf("密钥加载失败")
	} else {
		fmt.Println("密钥加载成功")
	}

	//打开要解密文件
	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("打开加密文件失败")
	}
	defer file.Close()

	//创建输出文件
	outputFile, err := os.Create("ouput.txt")
	if err != nil {
		return fmt.Errorf("输出文件创建失败")
	}
	defer outputFile.Close()

	//加载参数IV
	iv := make([]byte, 16)
	n, err := file.Read(iv)
	if err != nil {
		return fmt.Errorf("IV读取失败")
	}
	if n != 16 {
		return fmt.Errorf("IV长度错误")
	}

	//读取sha-256
	hash := make([]byte, 32)
	n, err = file.Read(hash)
	if err != nil {
		return fmt.Errorf("hash值读取失败")
	}
	if n != 32 {
		return fmt.Errorf("hash长度错误")
	}

	//创建AES
	blcok, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("AES创建失败")
	}
	stream := cipher.NewCTR(blcok, iv)

	//流解密
	reader := &cipher.StreamReader{S: stream, R: file}
	if _, err := io.Copy(outputFile, reader); err != nil {
		return fmt.Errorf("解密失败")
	}
	//验证完整性
	tmp, err := calculateHash("output.txt")
	if err != nil {
		return fmt.Errorf("sha-256计算失败")
	}

	if bytes.Equal([]byte(tmp), hash) {
		fmt.Println("解密完成，文件无破损")
	}
	return nil
}
