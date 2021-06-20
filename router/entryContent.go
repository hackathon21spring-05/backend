package router

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

var (
	httpClient = http.DefaultClient
)

func getEntryContent(entryUrl string) (model.Entry, error) {
	u, err := url.Parse(entryUrl)
	if err != nil {
		return model.Entry{}, fmt.Errorf("failed to parse url: %w", err)
	}

	var entry model.Entry
	switch u.Hostname() {
	case "q.trap.jp":
		// traQ message
		entry, err = getTraqMessage(u)
	case "wiki.trap.jp":
		// traP wiki message
		entry, err = getTrapWikiContent(u)
	default:
		// else
		entry, err = getWebContent(u)
	}
	if err != nil {
		return model.Entry{}, fmt.Errorf("failed to get entryData from web: %w", err)
	}
	return entry, nil
}

// url上のデータを取ってくる
func getWebContent(url *url.URL) (model.Entry, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return model.Entry{}, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return model.Entry{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.Entry{}, fmt.Errorf("fail to get message: (Status: %d %s)", res.StatusCode, res.Status)
	}

	// 読み取り
	buf, _ := ioutil.ReadAll(res.Body)

	// 文字コード判定
	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)
	// => EUC-JP

	// 文字コード変換
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)

	// HTMLパース
	doc, _ := goquery.NewDocumentFromReader(reader)

	// どんな情報を取ってくるのか，優先度など要相談
	title := doc.Find("title").Text()
	thumbnail := ""
	caption := ""
	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		attr := s.AttrOr("property", "")
		if attr == "og:image" {
			thumbnail = s.AttrOr("content", "")
		}
		if attr == "og:description" {
			caption = s.AttrOr("content", "")
		}
	})

	if title == "" {
		return model.Entry{}, fmt.Errorf("fail to get title from url")
	}
	return model.Entry{
		Url:       url.String(),
		Title:     title,
		Caption:   caption,
		Thumbnail: thumbnail,
	}, nil
}

// TODO
func getTraqMessage(url *url.URL) (model.Entry, error) {
	return model.Entry{}, nil
}

func getTrapWikiContent(url *url.URL) (model.Entry, error) {
	return model.Entry{}, nil
}
