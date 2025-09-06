package informationCollectionModule

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Scanning(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 0 {
		fmt.Printf("URL: %s, 状态码: %d\n", url, resp.StatusCode)
	}
}

func DirBlasting(url, path string) {
	url = strings.TrimRight(url, "/")
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		request := url + "/" + line
		Scanning(request)
	}
}
