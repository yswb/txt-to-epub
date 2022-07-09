package txt_to_epub

import (
	"bufio"
	"fmt"
	"github.com/bmaupin/go-epub"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

// TxtToEpub 把 txt 文件转为 epub 文件，txt 文件必须为 UTF-8 编码
//filePath：txt文件路径，author：作者，titleReg：每章标题的正则，用户分割章节
//第一章之前的内容都会作为序言，书籍结尾会合并到最后一章
func TxtToEpub(filePath, author, titleReg string) {
	source, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败", err)
		panic(err)
	}
	defer func() {
		err = source.Close()
		if err != nil {
			log.Fatal("出现错误", err)
		}
	}()
	//文件名，去除后缀
	fileName := strings.TrimSuffix(source.Name(), path.Ext(source.Name()))

	buf := bufio.NewReader(source)
	//标题最大长度
	maxTitleLen := 50
	//标题正则
	reg := regexp.MustCompile(titleReg)
	//扫描标题
	hasTitle := ScanTitle(buf, reg)
	if hasTitle {
		log.Println("开始处理")
		//重置reader到开始位置
		buf.Reset(source)
		_, err = source.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal("出现错误", err)
		}
	} else {
		log.Fatal("没有匹配到任何章节")
	}
	title, body := "", ""
	//统计信息
	word, chapter := 0, 0

	book := epub.NewEpub(fileName)
	book.SetAuthor(author)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("读取完毕")
				if title == "" {
					log.Fatal("没有匹配到任何章节")
				}
				//处理最后一章
				AddSection(book, title, body)
				break
			} else {
				log.Fatal("读取出现错误：", err)
			}
		}
		line = strings.TrimLeft(line, " 　")
		length := len([]rune(line))
		//忽略空行
		if length < 1 {
			continue
		}
		if length < maxTitleLen && reg.MatchString(line) {
			chapter++
			//保存上一个章节
			if title != "" {
				AddSection(book, title, body)
			}
			//第一章之前的内容
			if title == "" && body != "" {
				AddSection(book, "序言or简介", body)
			}
			title = line
			body = fmt.Sprintf("<h3>%s</h3>\n", line)
		} else {
			word += length
			line = fmt.Sprintf("<p>%s</p>\n", line)
			body = body + line
		}

	}

	savePath := path.Join(path.Dir(filePath), fileName+".epub")
	log.Println("开始输出文件：", savePath)
	err = book.Write(savePath)
	if err != nil {
		log.Fatal("文件保存失败", err)
	}
	log.Printf("处理完成,共%d章，%d字\n", chapter, word)

}
func AddSection(book *epub.Epub, title, body string) {
	_, err := book.AddSection(body, title, "", "")
	log.Println("当前章节：", title)
	if err != nil {
		log.Fatal("章节处理失败", err)
	}
}

func ScanTitle(buf *bufio.Reader, reg *regexp.Regexp) bool {
	hasTitle := false
	log.Println("开始匹配章节……")
	//先扫描一遍看能不能匹配到章节，找到就停止
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("读取出现错误：", err)
			}
		}
		if line != "" && reg.MatchString(line) {
			hasTitle = true
			break
		}
	}
	return hasTitle
}
