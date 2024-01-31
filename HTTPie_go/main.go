package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"net/http"
	Url "net/url"
	"os"

	"github.com/gookit/color"
)

// 存储所有的子命令
var commands = make(map[string]Request)

type Request interface {
	parse(args []string)
}

// get子命令
type GetCommand struct {
	url string
}

func (g *GetCommand) parse(args []string) {
	subCommand := flag.NewFlagSet("get", flag.ExitOnError)
	subCommand.StringVar(&g.url, "url", "", "Request URL")
	subCommand.Parse(args)
	validateUrl(&g.url)
	color.Grayf("get %s 。。。。。\n", g.url)
	parseRes("get", &g.url, nil)
}

// post子命令
type PostCommand struct {
	url  string
	args string
}

func (p *PostCommand) parse(args []string) {
	subCommand := flag.NewFlagSet("post", flag.ExitOnError)
	subCommand.StringVar(&p.url, "url", "", "Request URL")
	subCommand.StringVar(&p.args, "args", "", "Request args")
	subCommand.Parse(args)
	validateUrl(&p.url)
	color.Grayf("post %s 。。。。。\n", p.url)
	parseRes("post", &p.url, &p.args)
}

func main() {
	commands["get"] = &GetCommand{}
	commands["post"] = &PostCommand{}
	// 判断是否有命令
	if len(os.Args) < 2 {
		color.Redln("expected 'get' or 'post' subcommands")
		os.Exit(1)
	}
	// 选择子命令并解析
	commands[os.Args[1]].parse(os.Args[2:])
}

// 发送并解析请求返回
func parseRes(method string, url *string, args *string) {
	var res *http.Response
	var err error
	switch method {
	case "get":
		res, err = http.Get(*url)
	case "post":
		res, err = http.Post(*url, "application/json", bytes.NewBufferString(*args))
	}
	if err != nil {
		color.Redln("request failed: " + err.Error())
		os.Exit(1)
	}
	// 创建缓冲区，将所有的数据一起输出
	output := bufio.NewWriter(os.Stdout)
	output.WriteString(color.Blue.Sprintf("%s %s\n\n", res.Proto, res.Status))
	// 将请求头一行行打印
	for name, value := range res.Header {
		output.WriteString(color.Green.Sprintf("%s: ", name) + color.Sprintf("%s\n", value[0]))
	}
	output.WriteString("\n")
	body := make([]byte, 1024)
	// 读取所有返回数据
	for {
		n, err := res.Body.Read(body)
		if err != nil && err != io.EOF {
			color.Redln("read body failed: " + err.Error())
			os.Exit(1)
		}
		if n < 1024 {
			break
		}
		color.Yellow.Sprintf("%s", body[:n])
		output.WriteString(color.Yellow.Sprintf("%s", body[:n]))
	}
	output.Flush()
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
