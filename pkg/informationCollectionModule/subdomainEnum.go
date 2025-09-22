package informationCollectionModule

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func SubdomainEnum(host string) error {
	fmt.Println("start to ...")
	f, err := os.Open("../script/subdomain.txt")
	if err != nil {
		return fmt.Errorf("文件打开失败")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var rightSub []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	const numGoroutine = 100
	const maxbuffer = 1000
	var buffer = make(chan string, maxbuffer)

	for range numGoroutine {
		wg.Add(1)
		go func(host string, subs <-chan string) {
			defer wg.Done()
			for sub := range subs {
				ips, err := http.Get(fmt.Sprintf("http://%s.%s", sub, host))
				if err == nil && ips.StatusCode == 200 {
					mu.Lock()
					fmt.Println(fmt.Sprintf("%s.%s", sub, host))
					rightSub = append(rightSub, sub)
					mu.Unlock()
				}
			}
		}(host, buffer)
	}

	for scanner.Scan() {
		buffer <- scanner.Text()
	}
	close(buffer)
	wg.Wait()
	return fmt.Errorf("end...")
}
