package router

import (
	"net/url"
	"path"
	"strings"
)

// joinFileUrl 相対パスを絶対パスに修正する
func joinFileUrl(srcUrl *url.URL, fileUrl string) string {
	if strings.HasPrefix(fileUrl, ".") {
		return srcUrl.Scheme + "://" + srcUrl.Host + path.Join(srcUrl.Path, fileUrl)
	}
	return fileUrl
}
