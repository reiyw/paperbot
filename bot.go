package main

import (
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/carlescere/scheduler"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"github.com/reiyw/paperbot/queue"
	"github.com/reiyw/paperbot/translate"
	"log"
	"mvdan.cc/xurls"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api := slack.New(os.Getenv("PAPERBOT_SLACK_TOKEN"))
	arxivTrendChannelId := os.Getenv("ARXIV_TREND_CHANNEL_ID")
	botUserId := os.Getenv("BOT_USER_ID")
	botUserName := os.Getenv("BOT_USER_NAME")
	botIconUrl := os.Getenv("BOT_ICON_URL")

	channelQueue := queue.New()

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	requestAndSendTrendingPapers := func() {
		var papers []Paper
		trendingPapers := RequestTrendingPapersOnArxiv()
		for _, tp := range trendingPapers {
			p, err := FromArxivId(tp.Id)
			if err != nil {
				continue
			}
			papers = append(papers, *p)
		}
		for i, p := range papers {
			info := fmt.Sprintf("[%d tweets] %s", trendingPapers[i].TweetCount, formatAsPlainPaperInfo(p))
			rtm.SendMessage(rtm.NewOutgoingMessage(info, arxivTrendChannelId))
			channelQueue.PushBack(arxivTrendChannelId)
		}
	}
	_, err = scheduler.Every().Day().At("12:00").Run(requestAndSendTrendingPapers)
	if err != nil {
		fmt.Printf("Scheduler error: %s", err)
	}

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			urls := xurls.Relaxed().FindAllString(ev.Text, -1)
			var papers []Paper
			for _, url := range urls {
				p, err := Request(url)
				if err != nil {
					continue
				}
				papers = append(papers, *p)
			}
			for _, p := range papers {
				rtm.SendMessage(rtm.NewOutgoingMessage(formatAsPlainPaperInfo(p), ev.Channel))
				channelQueue.PushBack(ev.Channel)
			}

			//if ev.Text == "trend" {
			//	fmt.Println("trend")
			//	requestAndSendTrendingPapers()
			//}

			// if direct message or mention, do translate
			if strings.HasPrefix(ev.Channel, "D") || strings.Contains(ev.Text, botUserId) {
				text := strings.Replace(ev.Text, fmt.Sprintf("<@%s>", botUserId), "", 1)
				lang := whatlanggo.DetectLang(text)
				var langFrom string
				var langTo string
				switch lang {
				case whatlanggo.Jpn:
					langFrom = "ja"
					langTo = "en"
				case whatlanggo.Eng:
					langFrom = "en"
					langTo = "ja"
				default:
					langFrom = "auto"
					langTo = "ja"
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(translate.Google(langFrom, langTo, text), ev.Channel))
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		case *slack.AckMessage:
			fmt.Printf("AckMessage: %v\n", ev)
			urls := xurls.Relaxed().FindAllString(ev.Text, -1)
			fmt.Println("found URLs: ", urls)
			var papers []Paper
			for _, url := range urls {
				p, err := Request(url)
				if err != nil {
					continue
				}
				papers = append(papers, *p)
			}
			for _, p := range papers {
				params := slack.PostMessageParameters{
					Attachments:     []slack.Attachment{formatAsAttachment(p)},
					ThreadTimestamp: ev.Timestamp,
					IconURL:         botIconUrl,
					Username:        botUserName,
				}
				channel := fmt.Sprintf("%s", channelQueue.PopFront())
				_, _, _ = rtm.PostMessage(channel, "", params)
			}

		default:
			fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func formatAsPlainPaperInfo(p Paper) string {
	return fmt.Sprintf("%s. <%s |%s>. %d", concatAuthors(p.Authors), p.AbstUrl, p.Title, p.Year)
}

func formatAsAttachment(p Paper) slack.Attachment {
	abstTextJa := translate.Google("en", "ja", p.AbstText)
	attachment := slack.Attachment{
		Color:      p.Preserver.ToColor(),
		AuthorName: concatAuthors(p.Authors),
		Title:      p.Title,
		TitleLink:  p.AbstUrl,
		Text:       p.Comment,
		Fields: []slack.AttachmentField{
			{
				Title: "Abstract",
				Value: p.AbstText,
			}, {
				Title: "概要",
				Value: abstTextJa,
			},
		},
	}
	return attachment
}

func concatAuthors(authors []string) string {
	var b strings.Builder
	for i, author := range authors {
		if i != 0 {
			_, _ = b.WriteString(", ")
		}
		_, _ = b.WriteString(author)
	}
	return b.String()
}
