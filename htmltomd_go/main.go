package main
import (
	"flag"
	"io"
	"os"
	"fmt"
	"net/http"
	"github.com/TruthHun/html2md"
)
func main(){
	if len(os.Args) < 2{
		fmt.Println("not found 'url' and 'fileName'")
		os.Exit(1)
	}
	var url string
	var fileName string
	flag.StringVar(&url,"url","","要请求的地址")
	flag.StringVar(&fileName,"fileName","","要保存的文件名")
	flag.Parse()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	htmlStr := ""
	arr := make([]byte,1024)
	for {
		n, err := resp.Body.Read(arr)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			os.Exit(1)
		}
		if n < 1024 {
			htmlStr += string(arr[:n])
			break
		}
		htmlStr += string(arr[:n])
	}
	md := html2md.Convert(htmlStr)
	create, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = create.WriteString(md)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_ = create.Close()

}