package informationCollectionModule

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func PortScan(ip string, portStart, portEnd int) {

	var openPorts []int
	//portRange := portEnd - portStart + 1
	numGoroutine := 100
	var wg sync.WaitGroup
	var mu sync.Mutex
	const maxBuff = 1000
	ports := make(chan int, maxBuff)
	fmt.Println(time.Now())

	for i := 0; i < numGoroutine; i++ {
		wg.Add(1)
		go func(ip string, port <-chan int, openPorts *[]int) {
			defer wg.Done()
			for p := range ports {
				_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, p), time.Second*1)
				if err == nil {
					mu.Lock()
					fmt.Printf("%d is open on ip %s \n", p, ip)
					*openPorts = append(*openPorts, p)
					mu.Unlock()
				}

			}
		}(ip, ports, &openPorts)
	}
	for p := portStart; p <= portEnd; p++ {
		ports <- p
	}
	close(ports)
	wg.Wait()
	fmt.Println(time.Now())
}
