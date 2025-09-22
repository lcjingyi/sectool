package networkCommunicationModule

//DNS查询

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// 校验和
func CheckUdp(udpData []byte) (uint16, error) {
	var sum uint32
	if len(udpData)%2 != 0 {
		return 0, fmt.Errorf("UDP数据长度必须为偶数")
	}
	for i := 0; i < len(udpData); i = i + 2 {
		sum = sum + uint32(binary.BigEndian.Uint16(udpData[i:]))
	}

	//处理溢出
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}

	return uint16(sum), nil
}

// ip处理
func ipTransform(ip string) ([]byte, error) {
	list := strings.Split(ip, ".")
	if len(list) != 4 {
		return nil, fmt.Errorf("IP地址错误")
	}
	res := make([]byte, 4)
	for i := 0; i < 4; i++ {
		num, _ := strconv.Atoi(list[i])
		res[i] = uint8(num)
	}
	return res, nil
}

// udp数据包
func Udp(srcIP, dstIP string, srcPort, dstPort int, payload []byte) ([]byte, error) {
	udpLenth := 8 + len(payload)
	udpData := make([]byte, udpLenth)

	//源端口和目标端口
	binary.BigEndian.PutUint16(udpData[0:2], uint16(srcPort))
	binary.BigEndian.PutUint16(udpData[2:4], uint16(dstPort))

	//UDP长度
	binary.BigEndian.PutUint16(udpData[4:6], uint16(udpLenth))
	//校验和
	binary.BigEndian.PutUint16(udpData[6:], uint16(0x0000))
	//写入数据
	if len(payload) > 0 {
		copy(udpData[8:], payload)
	}

	//伪首部
	head := make([]byte, 12+udpLenth)
	tmp, _ := ipTransform(srcIP)
	copy(head[0:], tmp)
	tmp, _ = ipTransform(dstIP)
	copy(head[4:], tmp)
	binary.BigEndian.PutUint16(head[8:], 0x0011)
	binary.BigEndian.PutUint16(head[10:], uint16(udpLenth))
	copy(head[12:], udpData)
	//校验和
	checkSum, err := CheckUdp(head)
	if err != nil {
		return nil, fmt.Errorf("校验和计算错误")
	}
	binary.BigEndian.PutUint16(udpData[6:], checkSum)

	return udpData, nil
}

// 发送udp请求
func UdpRequest(udpData []byte) error {
	dnsServer := "8.8.8.8"
	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		return err
	}
	defer conn.Close()
	requestData, _ := Udp("192.168.1.1", dnsServer, 52, 64, udpData)
	//发送请求
	_, err = conn.Write(requestData)
	if err != nil {
		return err
	}

	return nil
}
