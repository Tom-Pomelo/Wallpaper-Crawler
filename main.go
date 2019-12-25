package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"
	"strconv"
	"math/rand"
	"sync"
)

var (
	url   = "https://zhuanlan.zhihu.com/p/25021104"
	reImg = `<img[\s\S]*?data-original="(http[\s\S]*?jpg)"`
	ch    = make(chan int, 10)
	wg    = sync.WaitGroup{}
)

func handleError(err error, at string) {
	if err != nil {
		fmt.Println(err, at)
	}
}

func getHTML(url string) string {
	response, err := http.Get(url)
	defer response.Body.Close()
	handleError(err, url)
	bytes, _ := ioutil.ReadAll(response.Body)
	html := string(bytes)
	return html
}

func getImgUrl() []string {
	html := getHTML(url)
	re := regexp.MustCompile(reImg)
	rets := re.FindAllStringSubmatch(html, -1)
	cnt := len(rets)
	fmt.Println("# of img:", cnt)
	imgUrls := make([]string, 0)
	for _, ret := range rets {
		imgUrl := ret[1]
		imgUrls = append(imgUrls, imgUrl)
	}
	return imgUrls
}

func downloadImg(url string) {
	response, err := http.Get(url)
	handleError(err, url)
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	file := "./img/" + strconv.Itoa(int(time.Now().UnixNano())) + ".jpg"
	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		handleError(err, url)
	}
}

func downloadImgAsync(url string) {
	wg.Add(1)
	go func() {
		ch <- rand.Int()
		downloadImg(url)
		<-ch
		wg.Done()
	}()
	wg.Wait()
}

func main() {
	imgUrls := getImgUrl()
	for _, imgUrl := range imgUrls {
		downloadImgAsync(imgUrl)
	}
}
