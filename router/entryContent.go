package router

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type selector struct {
	name    string
	process string
}

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
func getWebContent(u *url.URL) (model.Entry, error) {
	doc, err := requestDocument(u.String())
	if err != nil {
		return model.Entry{}, err
	}
	// どんな情報を取ってくるのか，優先度など要相談
	// タイトル取得
	titleSelector := []selector{
		{name: "title", process: "text"},
		{name: "meta[property='og:title']", process: "content"},
		{name: "h1", process: "text"},
	}
	title := getSelectorData(doc, titleSelector)
	title = strings.TrimSpace(title) // 前後の空白除去

	// サムネイル取得
	thumbnailSelector := []selector{
		{name: "meta[property='og:image']", process: "content"},
		{name: "link[rel='icon']", process: "href"},
		{name: "link[rel='apple-touch-icon']", process: "href"},
		{name: "link[rel='shortcut icon']", process: "href"},
	}
	thumbnail := getSelectorData(doc, thumbnailSelector)
	// urlが相対パスで指定されている場合，絶対パスに治す
	thumbnail = joinFileUrl(u, thumbnail)

	// キャプション取得
	captionSelector := []selector{
		{name: "meta[property='og:description']", process: "content"},
		{name: "meta[name='description']", process: "content"},
		{name: "meta[name='twitter:description']", process: "content"},
		{name: "h2", process: "text"}, // ここ，諸説あり
	}
	caption := getSelectorData(doc, captionSelector)

	if title == "" {
		return model.Entry{}, fmt.Errorf("fail to get title from url")
	}
	return model.Entry{
		Url:       u.String(),
		Title:     title,
		Caption:   caption,
		Thumbnail: thumbnail,
	}, nil
}

// getSelectorData指定されたセレクタの要素からメッセージを抜き取る
func getSelectorData(doc *goquery.Document, selectors []selector) string {
	message := ""
	for _, s := range selectors {
		switch s.process {
		case "text":
			message = doc.Find(s.name).Text()
		case "content":
			message = doc.Find(s.name).AttrOr("content", "")
		case "href":
			message = doc.Find(s.name).AttrOr("href", "")
		}
		if message != "" {
			break
		}
	}
	return message
}

// requestDocument 指定されたURLのDOMを取得する
func requestDocument(src string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("fail to get message: (Status: %d %s)", res.StatusCode, res.Status)
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
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// TODO
func getTraqMessage(url *url.URL) (model.Entry, error) {
	return model.Entry{}, nil
}

func getTrapWikiContent(url *url.URL) (model.Entry, error) {
	return model.Entry{}, nil
}
