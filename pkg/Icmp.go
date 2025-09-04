package pkg

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// 计算 ICMP 校验和
func checkSum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

// 发送 ICMP Echo 请求
func IcmpPing(dst string) error {
	conn, err := net.Dial("ip4:icmp", dst)
	if err != nil {
		return fmt.Errorf("创建连接失败: %v", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭连接失败:", err)
		}
	}(conn)

	// 构造 ICMP Echo Request 包
	var icmp [8 + 32]byte                                            // 8字节头+32字节数据
	icmp[0] = 8                                                      // Type: 8(Echo)
	icmp[1] = 0                                                      // Code: 0
	binary.BigEndian.PutUint16(icmp[4:], uint16(os.Getpid()&0xffff)) // ID
	binary.BigEndian.PutUint16(icmp[6:], 1)                          // Seq
	copy(icmp[8:], []byte("hello, icmp!"))                           // Data

	// 计算校验和
	cs := checkSum(icmp[:])
	binary.BigEndian.PutUint16(icmp[2:], cs)

	start := time.Now()
	_, err = conn.Write(icmp[:])
	if err != nil {
		return fmt.Errorf("发送失败: %v", err)
	}

	// 设置超时
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var resp [512]byte
	n, err := conn.Read(resp[:])
	if err != nil {
		return fmt.Errorf("接收失败: %v", err)
	}

	// 判断响应类型
	if n < 20+8 || resp[20] != 0 { // 20字节IP头+8字节ICMP头，Type=0(Echo Reply)
		return fmt.Errorf("无效响应")
	}

	fmt.Printf("收到来自 %s 的回复，耗时 %v\n", dst, time.Since(start))
	return nil
}
