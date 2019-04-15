package main

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
	"log"
	"strconv"
	"strings"
)

type TrendingPaper struct {
	Id         string
	TweetCount int
}

func RequestTrendingPapersOnArxiv() []TrendingPaper {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true)))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var res string
	err = c.Run(ctxt, requestBody(&res))
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}

	var ids []string
	doc.Find(".apaper").Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		ids = append(ids, id)
	})

	var tweetcounts []int
	doc.Find(".tweetcount").Each(func(_ int, s *goquery.Selection) {
		countStr := strings.Split(s.Text(), " ")[0]
		count, _ := strconv.ParseInt(countStr, 10, 32)
		tweetcounts = append(tweetcounts, int(count))
	})

	var min int
	if len(ids) < len(tweetcounts) {
		min = len(ids)
	} else {
		min = len(tweetcounts)
	}

	var papers []TrendingPaper
	for i := 0; i < min; i++ {
		papers = append(papers, TrendingPaper{ids[i], tweetcounts[i]})
	}

	return papers
}

func requestBody(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("http://www.arxiv-sanity.com/toptwtr"),
		chromedp.InnerHTML("//body", res),
	}
}
