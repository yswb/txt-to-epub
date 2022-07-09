package txt_to_epub

import (
	"github.com/bmaupin/go-epub"
	"testing"
)

func TestTxtToEpub(t *testing.T) {
	TxtToEpub("example.txt", "未知", `^第\S+章\s*(\S+)`)
}

func TestEpub(t *testing.T) {
	e := epub.NewEpub("书籍名称")
	e.SetAuthor("yswb")
	section1Body := `<h1>第一章</h1>
<p>这是内容</p>`
	_, _ = e.AddSection(section1Body, "第一章", "", "")

	err := e.Write("My EPUB.epub")
	if err != nil {
		panic(err)
	}
}
