package informationCollectionModule

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func requestSend(url string, dirs <-chan string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for d := range dirs {
		resp, err := http.Get(fmt.Sprintf("%s/%s", url, d))
		if err == nil {
			if resp.StatusCode == 200 {
				fmt.Println(d)
			}
		}
	}

}

func DirScan(url string) error {

	fmt.Println("start...")
	//打开文件
	f, err := os.Open("../script/dir.txt")
	if err != nil {
		return fmt.Errorf("文件打开失败")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var wg sync.WaitGroup
	var mu sync.Mutex
	const maxBuffer = 1000
	const numGoroutine = 100
	var buffer = make(chan string, maxBuffer) //管道

	for range numGoroutine {
		wg.Add(1)
		go requestSend(url, buffer, &wg, &mu)
	}

	for scanner.Scan() {
		buffer <- scanner.Text()
	}
	close(buffer)
	wg.Wait()
	return nil

}
