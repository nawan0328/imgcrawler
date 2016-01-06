package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	var a string // 키워드 입력 받을 변수
	for {
		fmt.Scanf("%s", &a)
		word := a
		client := &http.Client{}

		resp, err := client.Get("http://www.google.co.kr/search?q=" + word)
		if err != nil {
			fmt.Println(err)
			return
		}

		google, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		sp_google := strings.Split(string(google), "<h3 class=\"r\"><a href=\"/url?q=")

		for i, _ := range sp_google {
			if i+1 < len(sp_google) {
				link_google := strings.Split(sp_google[i+1], "&amp;sa=U")
				go target(link_google[0], parseUrl(link_google[0]))
			}
		}
	}

}

func target(targetUrl, host string) {

	res, err := http.Get(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	web, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	sp_web := strings.Split(string(web), "<img src=\"")

	for i, _ := range sp_web {
		if i+1 < len(sp_web) {
			link_web := strings.Split(sp_web[i+1], "\"")
			//fmt.Println(host)
			//fmt.Println(link_web[0])
			compare := strings.Split(link_web[0], "/")
			if compare[0] == "" {
				//fmt.Println(strings.Replace(link_web[0], "/", host, 1))
				go fileGet(strings.Replace(link_web[0], "/", host, 1), i)
			} else {
				go fileGet(link_web[0], i)
			}
		}
	}
}

func parseUrl(targetUrl string) (host string) {

	u, err := url.Parse(targetUrl)
	if err != nil {
		fmt.Println(err)
	}
	host = u.Scheme + "://" + u.Host
	return
}

func fileGet(src string, count int) {

	u, err := url.Parse(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.Get(src)
	if err != nil {
		//fmt.Println("ERROR")
		return
	}
	defer resp.Body.Close()
	out, err := os.Create(strings.Replace(u.Host, ".", "", -1) + "_" + strconv.Itoa(count) + ".jpg")
	io.Copy(out, resp.Body)
	fmt.Println("get image = " + strings.Replace(u.Host, ".", "", -1) + "_" + strconv.Itoa(count) + ".jpg")
}
