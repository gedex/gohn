package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nsf/termbox-go"
	"strings"
	"time"
)

const (
	HNURL     = "https://news.ycombinator.com"
	mainTitle = "Hacker News"
	footer    = "Press ESC to exit"
	maxItems  = 10
)

type item struct {
	title, url, points string
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	showLatest()

loop:
	for {
		select {
		case ev := <-event_queue:
			switch ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyEsc {
					break loop
				}
			case termbox.EventResize:
				showLatest()
			}
		case <-time.After(time.Second * 5):
			showLatest()
		}
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func showLatest() {
	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(HNURL); e != nil {
		panic(e.Error())
	}

	items, maxWidth := getItems(doc)
	print(items, maxWidth)
}

func getItems(doc *goquery.Document) (items []item, maxWidth int) {
	doc.Find("td.title a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == maxItems {
			return false
		}

		if s.Text() == "More" {
			return true
		}

		href, _ := s.Attr("href")
		title := s.Text()
		points := s.Parent().Parent().Next().Find("span").Text()
		a, b := len(fmt.Sprintf("%s (%s)", title, points)), len(href)
		maxWidth = max(a, b, maxWidth)

		items = append(items, item{
			title:  title,
			url:    href,
			points: points,
		})

		return true
	})
	return
}

func print(items []item, maxWidth int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Print header
	printLine(mainTitle, 0, maxWidth)

	y := 1

	// Items from HN
	for _, i := range items {
		lenTop := len(fmt.Sprintf("%s(%s)", i.title, i.points))
		strTop := fmt.Sprintf("%s%s(%s)", i.title, strings.Repeat(" ", maxWidth-lenTop), i.points)
		strURL := fmt.Sprintf("%s%s", i.url, strings.Repeat(" ", maxWidth-len(i.url)))

		// Print title and points
		for x := 0; x < maxWidth; x++ {
			termbox.SetCell(x, y, rune(strTop[x]), termbox.ColorBlack, termbox.ColorWhite)
		}
		y++

		// Print URL
		for x := 0; x < maxWidth; x++ {
			termbox.SetCell(x, y, rune(strURL[x]), termbox.ColorBlue|termbox.AttrBold, termbox.ColorWhite)
		}
		y++
	}

	// Print footer
	printLine(footer, y, maxWidth)

	termbox.Flush()
}

func printLine(s string, y, maxWidth int) {
	var ch rune
	l := len(s)
	for x := 0; x < maxWidth; x++ {
		if x < l {
			ch = rune(s[x])
		} else {
			ch = ' '
		}
		termbox.SetCell(x, y, ch, termbox.ColorBlack, termbox.ColorMagenta)
	}
}

func max(a, b, c int) int {
	if a >= b && a >= c {
		return a
	} else if b >= a && b >= c {
		return b
	}
	return c
}
