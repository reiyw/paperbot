package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Paper struct {
	Id       string
	Title    string
	Authors  []string
	Volume   string
	Venue    string
	Year     int
	PdfUrl   string
	HtmlUrl  string
	AbstText string
	AbstUrl  string
	BibText  string
	BibUrl   string
	Comment  string
	Preserver
}

func Request(url string) (*Paper, error) {
	preserver, err := DetectPreserver(url)
	if err != nil {
		return nil, err
	}
	switch preserver {
	case Arxiv:
		return FromArxivUrl(url)
	case Aclweb:
		return FromAclweb(url)
	case OpenReview:
		return FromOpenreview(url)
	default:
		return nil, fmt.Errorf("notimplemented")
	}
}

func FromArxivId(id string) (*Paper, error) {
	var paper Paper
	paper.Id = id
	paper.AbstUrl = fmt.Sprintf("https://arxiv.org/abs/%s", paper.Id)
	paper.PdfUrl = fmt.Sprintf("https://arxiv.org/pdf/%s.pdf", paper.Id)
	paper.HtmlUrl = fmt.Sprintf("https://www.arxiv-vanity.com/papers/%s/", paper.Id)

	res, err := http.Get(paper.AbstUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	paper.Title, _ = doc.Find(`meta[name="citation_title"]`).Attr("content")

	authorStr := doc.Find(".authors").Text()
	authors := strings.Split(strings.Replace(authorStr, "Authors:", "", 1), ",")
	paper.Authors = []string{}
	for _, author := range authors {
		paper.Authors = append(paper.Authors, strings.TrimSpace(author))
	}

	citationDate, _ := doc.Find(`meta[name="citation_date"]`).Attr("content")
	parsedDate, _ := dateparse.ParseAny(citationDate)
	paper.Year = parsedDate.Year()

	abstText := doc.Find(".abstract").Text()
	abstText = strings.Replace(abstText, "Abstract:", "", 1)
	abstText = strings.Replace(abstText, "\n", " ", -1)
	paper.AbstText = strings.TrimSpace(abstText)

	comment := doc.Find(".comments").Text()
	comment = strings.Replace(comment, "\n", " ", -1)
	paper.Comment = strings.TrimSpace(comment)

	paper.Preserver = Arxiv

	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &paper, nil

}

func FromArxivUrl(rawurl string) (*Paper, error) {
	// https://arxiv.org/pdf/1811.01458v1.pdf
	//                    id ^^^^^^^^^^^^
	split := strings.Split(rawurl, "/")
	id := strings.Split(split[len(split)-1], ".pdf")[0]
	return FromArxivId(id)
}

func FromAclweb(rawurl string) (*Paper, error) {
	// https://aclweb.org/anthology/D16-1112.pdf
	//                           id ^^^^^^^^
	var paper Paper
	split := strings.Split(rawurl, "/")
	id := strings.Split(split[len(split)-1], ".pdf")[0]
	id = strings.Split(id, ".bib")[0]
	paper.Id = strings.ToUpper(id)
	paper.AbstUrl = fmt.Sprintf("https://aclanthology.info/papers/%s/%s", paper.Id, strings.ToLower(paper.Id))
	paper.PdfUrl = fmt.Sprintf("http://aclweb.org/anthology/%s", paper.Id)
	paper.BibUrl = fmt.Sprintf("http://aclweb.org/anthology/%s.bib", paper.Id)

	res, err := http.Get(paper.AbstUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	paper.Title, _ = doc.Find(`meta[name="citation_title"]`).Attr("content")

	doc.Find(`meta[name="citation_author"]`).Each(func(_ int, s *goquery.Selection) {
		author, _ := s.Attr("content")
		paper.Authors = append(paper.Authors, author)
	})

	paper.Volume, _ = doc.Find(`meta[name="citation_journal_title"]`).Attr("content")

	paper.Venue = aclPrefixToVenue(paper.Id[0:1])

	yearStr, _ := doc.Find(`meta[name="citation_publication_date"]`).Attr("content")
	year, _ := strconv.ParseInt(yearStr, 10, 32)
	paper.Year = int(year)

	paper.Preserver = Aclweb

	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &paper, nil
}

func aclPrefixToVenue(prefix string) string {
	switch prefix {
	case "J":
		return "CL"
	case "Q":
		return "TACL"
	case "P":
		return "ACL"
	case "E":
		return "EACL"
	case "N":
		return "NAACL"
	case "S":
		return "SEMEVAL"
	case "D":
		return "EMNLP"
	case "K":
		return "CONLL"
	default:
		return ""
	}
}

func FromOpenreview(rawurl string) (*Paper, error) {
	return nil, fmt.Errorf("error")
}

type Preserver int

const (
	Arxiv Preserver = iota
	Aclweb
	OpenReview
)

func DetectPreserver(rawurl string) (Preserver, error) {
	parsed, _ := url.Parse(rawurl)
	switch parsed.Hostname() {
	case "arxiv.org":
		return Arxiv, nil
	case "aclweb.org", "aclanthology.info", "aclanthology.coli.uni-saarland.de":
		return Aclweb, nil
	case "openreview.net":
		return OpenReview, nil
	default:
		return -1, fmt.Errorf("given URL is not supported: %s", rawurl)
	}
}

func (p Preserver) ToColor() string {
	switch p {
	case Arxiv:
		return "#b31b1b"
	case Aclweb:
		return "#cc0000"
	case OpenReview:
		return "#8c1b13"
	default:
		return ""
	}
}
