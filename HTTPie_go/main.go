package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	Url "net/url"
	"os"

	"github.com/gookit/color"
)

func main() {
	// 请求地址
	var url string
	// 判断是否有命令
	if len(os.Args) < 2 {
		color.Redln("expected 'get' or 'post' subcommands")
		os.Exit(1)
	}
	// 选择子命令
	switch os.Args[1] {
	case "get":
		// get子命令
		getCommand := flag.NewFlagSet("get", flag.ExitOnError)
		getCommand.StringVar(&url, "url", "", "要请求的网络地址")
		// 解析参数
		getCommand.Parse(os.Args[2:])
		validateUrl(&url)
		color.Bluef("get %s 。。。。。\n", url)
		parseRes("get", &url)
	case "post":
		// post子命令
		postCommand := flag.NewFlagSet("post", flag.ExitOnError)
		postCommand.StringVar(&url, "url", "", "要请求的网络地址")
		postUrl := color.Blue.Sprintf("post %s", url)
		color.Println(postUrl)
		postCommand.Parse(os.Args[2:])
	}
}

// 发送并解析请求返回
func parseRes(method string, url *string) {
	var res *http.Response
	var err error
	switch method {
	case "get":
		res, err = http.Get(*url)
	case "post":
		res, err = http.Post(*url, "application/json", nil)
	}
	if err != nil {
		color.Redln("request failed: " + err.Error())
		os.Exit(1)
	}
	color.Blueln(res.Proto, res.Status)
	fmt.Println()
	for name, value := range res.Header {
		color.Greenf("%s: ", name)
		color.Printf("%s\n", value[0])
	}
	fmt.Println()
	body := make([]byte, 1024)
	for {
		n, err := res.Body.Read(body)
		if err != nil && err != io.EOF {
			color.Redln("read body failed: " + err.Error())
		}
		if n < 1024 {
			break
		}
		color.Yellowln(string(body))
	}
	defer res.Body.Close()

}

// 验证url是否合法
func validateUrl(url *string) {
	_, err := Url.ParseRequestURI(*url)
	if err != nil {
		parseErr := color.Red.Sprint("invalid url: " + *url)
		color.Println(parseErr)
		os.Exit(1)
	}
}
