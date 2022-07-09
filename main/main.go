package main

import (
	"flag"
	"log"
	"strings"
	tte "txt-to-epub"
)

func main() {
	var source = flag.String("s", "", "txt 源文件")
	var author = flag.String("a", "未知", "作者")
	var titleReg = flag.String("reg", `^第\S+章\s*.+`, "章节标题的正则")
	flag.Parse()
	if *source == "" {
		log.Fatal("请输出正确的文件名称")
	}
	name := strings.ToLower(*source)
	if !strings.HasSuffix(name, ".txt") {
		log.Fatal("只支持 txt 文件")
	}

	tte.TxtToEpub(name, *author, *titleReg)
}
