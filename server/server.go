package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/auguslxk/LCS/article"
	"github.com/auguslxk/LCS/lib"
	"github.com/micro/go-log"
)

func Run() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		artA, _ := lib.ReadFile("articleA.txt")
		artB, _ := lib.ReadFile("articleB.txt")
		articleA := article.ArticleInit(artA)
		articleB := article.ArticleInit(artB)
		article.DuplicateChecking(articleA, articleB)
		_, _ = io.WriteString(w, genratePage(articleA, articleB))
	}
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateHtml(art *article.Article, title string) string {
	html := bytes.Buffer{}

	html.WriteString("<div style='width:50%'>")
	html.WriteString(fmt.Sprintf("<h3>%s</h3>", title))
	for _, sentence := range art.Sentences {
		html.WriteString(makeSentence(&sentence))
	}

	html.WriteString("</div>")

	return html.String()
}

func wrapSentence(sentence *article.Sentence) string {
	buf := bytes.Buffer{}
	if sentence.DuplicateId == -1 {
		buf.WriteString("<div>")
		buf.WriteString(string(sentence.Content))
		buf.WriteString("</div>")
		return buf.String()
	}
	buf.WriteString("<div><b style='color:black'>" + strconv.Itoa(sentence.DuplicateId) + "</b>")

	segment := []int{0}
	begin := 0
	end := begin
	//fmt.Printf("before = %+v\n", sentence.DupIndex)
	for end < len(sentence.DupIndex) {
		if end+1 < len(sentence.DupIndex) && sentence.DupIndex[end+1] == sentence.DupIndex[end]+1 {
			end++
		} else {
			segment = append(segment, sentence.DupIndex[begin], sentence.DupIndex[end]+1)
			begin = end + 1
			end = begin
		}
	}
	segment = append(segment, len(sentence.Content))
	//fmt.Printf("end= %+v\n", segment)
	color := ""
	for i := 0; i < len(segment)-1; i++ {
		if segment[i] == segment[i+1] {
			if color == "" {
				color = "red"
			} else {
				color = ""
			}
			continue
		}
		buf.WriteString("<span style='color:" + color + "'>")
		buf.WriteString(string(sentence.Content[segment[i]:segment[i+1]]))
		if color == "" {
			color = "red"
		} else {
			color = ""
		}
		buf.WriteString("</span>")
	}
	buf.WriteString("</div>")
	return buf.String()
}

func makeSentence(sentence *article.Sentence) string {
	buf := bytes.Buffer{}
	content := wrapSentence(sentence)
	for _, char := range content {
		if char == '\n' || char == '\r' {
			buf.WriteString("<br>")
		} else {
			buf.WriteRune(char)
		}
	}
	return buf.String()
}

func genratePage(arta, artb *article.Article) string {
	page := bytes.Buffer{}

	page.WriteString("<html><head><link rel='icon' href='favicon.ico' type='image/x-icon' /></head>")

	aDupRatio := fmt.Sprintf("<h1>A文章重复率%f%%</h1>", float64(arta.DuplicateSize*100)/float64(arta.CharSize))
	bDupRatio := fmt.Sprintf("<h1>B文章重复率%f%%</h1>", float64(artb.DuplicateSize*100)/float64(artb.CharSize))

	page.WriteString(aDupRatio)
	page.WriteString(bDupRatio)

	page.WriteString("<div style='display:flex'>")

	page.WriteString(generateHtml(arta, "文章A"))
	page.WriteString(generateHtml(artb, "文章B"))

	page.WriteString("</div></html>")

	return page.String()
}
