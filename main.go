package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"log"
	"net/http"
)

func getBookBaseInfo(doc *goquery.Document) (bookName, author, description, latestChapterName string) {
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, _ := s.Attr("property")
		content, _ := s.Attr("content")
		switch property {
		case "og:novel:book_name":
			bookName = content
		case "og:novel:author":
			author = content
		case "og:description":
			description = content
		case "og:novel:latest_chapter_name":
			latestChapterName = content
		}
	})

	log.Printf("%s\n%s\n%s\n%s\n", bookName, author, description, latestChapterName)
	return
}

func getChapterLink(doc *goquery.Document) []string {
	links := make([]string, 0, 100)
	doc.Find("dd a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		links = append(links, href)
	})
	return links
}

func getText(doc *goquery.Document) string {
	text := ""
	//标题
	doc.Find("body div div div div h1").Each(func(i int, s *goquery.Selection) {
		log.Println(s.Text())
		return
	})
	//内容
	doc.Find("body div div div div").Each(func(i int, s *goquery.Selection) {
		if id, ok := s.Attr("id"); ok {
			if id == "content" {
				text = s.Text()
				return
			}
		}
	})
	log.Println(text)
	return text
}

func main() {
	var (
		dec     = mahonia.NewDecoder("gbk")
		baseUrl = "https://www.bequge.com"
	)

	res, err := http.Get(baseUrl + "/0_124/")
	if err != nil {
		log.Fatal(err)
		return
	}

	body := dec.NewReader(res.Body)

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
		return
	}
	getBookBaseInfo(doc)
	links := getChapterLink(doc)
	for _, value := range links {
		res, err = http.Get(baseUrl + value)
		if err != nil {
			log.Fatal(err)
			return
		}

		body = dec.NewReader(res.Body)
		chapterDoc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			log.Fatal(err)
			return
		}
		getText(chapterDoc)
		return
	}
}
