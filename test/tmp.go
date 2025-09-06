package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

// 生成AES密钥(256位)并保存到文件
func generateKey(keyPath string) error {
	key := make([]byte, 32) // 256位密钥
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	return os.WriteFile(keyPath, key, 0600) // 仅当前用户可读写
}

// 从文件加载密钥
func loadKey(keyPath string) ([]byte, error) {
	return os.ReadFile(keyPath)
}

// 计算文件的SHA-256哈希值
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

// 加密文件（流式处理）
func encryptFile(inputPath, outputPath, keyPath string) error {
	// 加载密钥
	key, err := loadKey(keyPath)
	if err != nil {
		// 密钥不存在则生成新密钥
		if os.IsNotExist(err) {
			if err := generateKey(keyPath); err != nil {
				return fmt.Errorf("生成密钥失败: %v", err)
			}
			key, err = loadKey(keyPath)
			if err != nil {
				return fmt.Errorf("加载新生成的密钥失败: %v", err)
			}
		} else {
			return fmt.Errorf("加载密钥失败: %v", err)
		}
	}

	// 打开输入文件
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开输入文件失败: %v", err)
	}
	defer inFile.Close()

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outFile.Close()

	// 生成随机IV (AES-CTR模式推荐16字节IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return fmt.Errorf("生成IV失败: %v", err)
	}

	// 先写入IV到输出文件头部
	if _, err := outFile.Write(iv); err != nil {
		return fmt.Errorf("写入IV失败: %v", err)
	}

	// 计算原始文件哈希并写入
	hash, err := calculateHash(inputPath)
	if err != nil {
		return fmt.Errorf("计算文件哈希失败: %v", err)
	}
	// 写入哈希长度和哈希值
	hashBytes := []byte(hash)
	if _, err := outFile.Write([]byte{byte(len(hashBytes))}); err != nil {
		return fmt.Errorf("写入哈希长度失败: %v", err)
	}
	if _, err := outFile.Write(hashBytes); err != nil {
		return fmt.Errorf("写入哈希值失败: %v", err)
	}

	// 创建AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("创建AES加密器失败: %v", err)
	}

	// 创建CTR模式流加密器
	stream := cipher.NewCTR(block, iv)

	// 流式加密：使用自定义Writer分块处理
	writer := &cipher.StreamWriter{S: stream, W: outFile}
	if _, err := io.Copy(writer, inFile); err != nil {
		return fmt.Errorf("加密过程失败: %v", err)
	}

	fmt.Printf("文件加密成功: %s -> %s\n", inputPath, outputPath)
	fmt.Printf("密钥文件: %s (请妥善保管!)\n", keyPath)
	return nil
}

// 解密文件（流式处理）
func decryptFile(inputPath, outputPath, keyPath string) error {
	// 加载密钥
	key, err := loadKey(keyPath)
	if err != nil {
		return fmt.Errorf("加载密钥失败: %v", err)
	}

	// 打开加密文件
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开加密文件失败: %v", err)
	}
	defer inFile.Close()

	// 读取IV (AES块大小)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return fmt.Errorf("读取IV失败: %v", err)
	}

	// 读取哈希长度和哈希值
	var hashLen byte
	if _, err := inFile.Read([]byte{hashLen}); err != nil {
		return fmt.Errorf("读取哈希长度失败: %v", err)
	}
	originalHash := make([]byte, hashLen)
	if _, err := io.ReadFull(inFile, originalHash); err != nil {
		return fmt.Errorf("读取哈希值失败: %v", err)
	}

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outFile.Close()

	// 创建AES解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("创建AES解密器失败: %v", err)
	}

	// 创建CTR模式流解密器（CTR加密解密使用相同的流）
	stream := cipher.NewCTR(block, iv)

	// 流式解密：使用自定义Reader分块处理
	reader := &cipher.StreamReader{S: stream, R: inFile}
	if _, err := io.Copy(outFile, reader); err != nil {
		return fmt.Errorf("解密过程失败: %v", err)
	}

	// 验证文件完整性
	currentHash, err := calculateHash(outputPath)
	if err != nil {
		return fmt.Errorf("计算解密文件哈希失败: %v", err)
	}

	if currentHash != string(originalHash) {
		return fmt.Errorf("文件完整性验证失败，可能被篡改或密钥错误")
	}

	fmt.Printf("文件解密成功: %s -> %s\n", inputPath, outputPath)
	return nil
}

func main() {
	// 命令行参数
	operation := flag.String("op", "", "操作: encrypt 或 decrypt")
	input := flag.String("in", "", "输入文件路径")
	output := flag.String("out", "", "输出文件路径")
	key := flag.String("key", "secret.key", "密钥文件路径")
	flag.Parse()

	if *operation == "" || *input == "" || *output == "" {
		fmt.Println("用法:")
		fmt.Println("  加密: filecrypt -op encrypt -in 原始文件 -out 加密文件 -key 密钥文件")
		fmt.Println("  解密: filecrypt -op decrypt -in 加密文件 -out 解密文件 -key 密钥文件")
		os.Exit(1)
	}

	switch *operation {
	case "encrypt":
		if err := encryptFile(*input, *output, *key); err != nil {
			fmt.Printf("加密失败: %v\n", err)
			os.Exit(1)
		}
	case "decrypt":
		if err := decryptFile(*input, *output, *key); err != nil {
			fmt.Printf("解密失败: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("无效操作，使用 'encrypt' 或 'decrypt'")
		os.Exit(1)
	}
}
