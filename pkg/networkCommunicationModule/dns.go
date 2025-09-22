package networkCommunicationModule

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// DNS头部结构
type dnsHeader struct {
	id      uint16 // 会话标识
	flags   uint16 // 标志
	qdcount uint16 // 问题数
	ancount uint16 // 回答资源记录数
	nscount uint16 // 权威名称服务器数
	arcount uint16 // 附加资源记录数
}

// DNS问题结构
type dnsQuestion struct {
	qname  []byte // 查询名
	qtype  uint16 // 查询类型
	qclass uint16 // 查询类
}

// 主函数
func lmain() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: ", os.Args[0], " <域名>")
		fmt.Println("示例: ", os.Args[0], " example.com")
		os.Exit(1)
	}

	domain := os.Args[1]
	dnsServer := "8.8.8.8:53" // Google DNS服务器

	// 构建DNS查询包
	packet, err := buildDNSQuery(domain)
	if err != nil {
		fmt.Println("构建DNS查询包失败:", err)
		os.Exit(1)
	}

	// 发送DNS查询并接收响应
	response, err := sendDNSQuery(packet, dnsServer)
	if err != nil {
		fmt.Println("发送DNS查询失败:", err)
		os.Exit(1)
	}

	// 解析DNS响应
	err = parseDNSResponse(response)
	if err != nil {
		fmt.Println("解析DNS响应失败:", err)
		os.Exit(1)
	}
}

// 构建DNS查询包
func buildDNSQuery(domain string) ([]byte, error) {
	// 创建头部
	header := dnsHeader{
		id:      0x1234, // 随机会话ID
		flags:   0x0100, // 标准查询
		qdcount: 1,      // 1个查询
		ancount: 0,
		nscount: 0,
		arcount: 0,
	}

	// 创建问题
	question, err := buildQuestion(domain)
	if err != nil {
		return nil, err
	}

	// 计算总长度
	totalLength := 12 + len(question.qname) + 4 // 头部12字节 + 问题 + 类型2字节 + 类2字节

	// 构建数据包
	packet := make([]byte, totalLength)
	offset := 0

	// 写入头部
	binary.BigEndian.PutUint16(packet[offset:], header.id)
	offset += 2
	binary.BigEndian.PutUint16(packet[offset:], header.flags)
	offset += 2
	binary.BigEndian.PutUint16(packet[offset:], header.qdcount)
	offset += 2
	binary.BigEndian.PutUint16(packet[offset:], header.ancount)
	offset += 2
	binary.BigEndian.PutUint16(packet[offset:], header.nscount)
	offset += 2
	binary.BigEndian.PutUint16(packet[offset:], header.arcount)
	offset += 2

	// 写入问题
	copy(packet[offset:], question.qname)
	offset += len(question.qname)
	binary.BigEndian.PutUint16(packet[offset:], question.qtype)
	offset += 2
	binary.BigEndian.PutUint16(packet[offset:], question.qclass)

	return packet, nil
}

// 构建DNS问题部分
func buildQuestion(domain string) (dnsQuestion, error) {
	// 将域名转换为DNS格式（每个部分前加长度）
	labels := splitDomain(domain)
	var qname []byte

	for _, label := range labels {
		if len(label) > 63 { // DNS标签最大长度为63
			return dnsQuestion{}, fmt.Errorf("域名标签过长: %s", label)
		}
		qname = append(qname, byte(len(label)))
		qname = append(qname, []byte(label)...)
	}
	qname = append(qname, 0x00) // 域名结束标志

	return dnsQuestion{
		qname:  qname,
		qtype:  0x0001, // A记录
		qclass: 0x0001, // IN类（互联网）
	}, nil
}

// 分割域名为标签
func splitDomain(domain string) []string {
	var labels []string
	start := 0
	for i := 0; i < len(domain); i++ {
		if domain[i] == '.' {
			if i > start {
				labels = append(labels, domain[start:i])
			}
			start = i + 1
		}
	}
	if start < len(domain) {
		labels = append(labels, domain[start:])
	}
	return labels
}

// 发送DNS查询并接收响应
func sendDNSQuery(packet []byte, dnsServer string) ([]byte, error) {
	// 创建UDP连接
	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 设置超时
	// conn.SetDeadline(time.Now().Add(5 * time.Second))

	// 发送查询包
	_, err = conn.Write(packet)
	if err != nil {
		return nil, err
	}

	// 接收响应
	buffer := make([]byte, 512) // DNS响应最大长度通常为512字节
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// 解析DNS响应
func parseDNSResponse(response []byte) error {
	if len(response) < 12 {
		return fmt.Errorf("响应数据包太短")
	}

	// 解析头部
	header := dnsHeader{
		id:      binary.BigEndian.Uint16(response[0:2]),
		flags:   binary.BigEndian.Uint16(response[2:4]),
		qdcount: binary.BigEndian.Uint16(response[4:6]),
		ancount: binary.BigEndian.Uint16(response[6:8]),
		nscount: binary.BigEndian.Uint16(response[8:10]),
		arcount: binary.BigEndian.Uint16(response[10:12]),
	}

	// 检查响应标志
	qr := (header.flags >> 15) & 1 // 1表示响应
	rcode := header.flags & 0x000F // 响应码
	if qr != 1 {
		return fmt.Errorf("不是一个DNS响应")
	}
	if rcode != 0 {
		return fmt.Errorf("DNS查询失败，响应码: %d", rcode)
	}

	fmt.Printf("DNS响应: ID=0x%x, 回答数=%d\n", header.id, header.ancount)

	// 跳过问题部分
	offset := 12
	for i := 0; i < int(header.qdcount); i++ {
		// 跳过域名
		for {
			if offset >= len(response) {
				return fmt.Errorf("响应格式错误")
			}
			length := int(response[offset])
			offset++
			if length == 0 {
				break
			}
			offset += length
		}
		offset += 4 // 跳过qtype和qclass
	}

	// 解析回答部分
	for i := 0; i < int(header.ancount); i++ {
		if offset >= len(response) {
			break
		}

		// 跳过域名（简化处理，实际可能包含指针）
		offset = skipName(response, offset)

		// 解析类型、类、TTL和数据长度
		if offset+10 > len(response) {
			return fmt.Errorf("响应格式错误")
		}
		rtype := binary.BigEndian.Uint16(response[offset : offset+2])
		_ = binary.BigEndian.Uint16(response[offset+2 : offset+4])
		ttl := binary.BigEndian.Uint32(response[offset+4 : offset+8])
		dataLen := binary.BigEndian.Uint16(response[offset+8 : offset+10])
		offset += 10

		// 解析数据
		if offset+int(dataLen) > len(response) {
			return fmt.Errorf("响应数据长度错误")
		}

		// 如果是A记录，解析IP地址
		if rtype == 1 { // A记录
			if dataLen == 4 {
				ip := net.IPv4(response[offset], response[offset+1], response[offset+2], response[offset+3])
				fmt.Printf("IP地址: %s, TTL: %d\n", ip.String(), ttl)
			}
		} else if rtype == 5 { // CNAME记录
			cname, newOffset := parseName(response, offset)
			fmt.Printf("CNAME: %s, TTL: %d\n", cname, ttl)
			offset = newOffset - offset + int(dataLen) // 调整偏移量
		}

		offset += int(dataLen)
	}

	return nil
}

// 跳过域名（处理指针）
func skipName(response []byte, offset int) int {
	for {
		if offset >= len(response) {
			return offset
		}

		// 检查是否是指针 (最高两位为1)
		if (response[offset] & 0xC0) == 0xC0 {
			offset += 2 // 指针占2字节
			break
		}

		length := int(response[offset])
		offset++
		if length == 0 {
			break
		}
		offset += length
	}
	return offset
}

// 解析域名（处理指针）
func parseName(response []byte, offset int) (string, int) {
	var name string
	originalOffset := offset
	seenPointer := false

	for {
		if offset >= len(response) {
			break
		}

		// 检查是否是指针 (最高两位为1)
		if (response[offset] & 0xC0) == 0xC0 {
			if !seenPointer {
				originalOffset = offset + 2 // 记录指针后的偏移量
				seenPointer = true
			}
			// 计算指针指向的偏移量 (去除最高两位)
			pointer := int(binary.BigEndian.Uint16(response[offset:offset+2]) & 0x3FFF)
			offset = pointer
			continue
		}

		length := int(response[offset])
		offset++
		if length == 0 {
			break
		}

		if len(name) > 0 {
			name += "."
		}
		name += string(response[offset : offset+length])
		offset += length
	}

	if !seenPointer {
		originalOffset = offset
	}

	return name, originalOffset
}
