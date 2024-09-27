package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Text string
	Href string
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}

	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var ret string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += strings.Join(strings.Fields(text(c)), " ")
	}

	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}

	ret.Text = text(n)

	return ret
}

func Parse(r io.Reader) ([]Link, error) {
	doc, htmlParseErr := html.Parse(r)

	if htmlParseErr != nil {
		return nil, htmlParseErr
	}

	var links []Link

	nodes := linkNodes(doc)

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func main() {
	htmlFile := flag.String("file", "1.html", "HTML file from which links will be extracted")

	flag.Parse()

	file, err := os.Open(*htmlFile)

	if err != nil {
		log.Fatal(err)
	}

	links, err := Parse(file)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", links)
}
