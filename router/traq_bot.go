package router

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hackathon21spring-05/linq-backend/model"
)

// ここでやるの良くない気がするけどめんどい
var entryURL, _ = url.Parse("https://linq.trap.games/entry")

func calcHMACSHA1(message, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func PostTraqMessage(entry *model.Entry) error {
	message := fmt.Sprintf("### [%s](%s)\n### !![Add tags](%s?url=%s)!!",
		entry.Title,
		entry.Url,
		entryURL,
		entry.Url,
	)
	url := fmt.Sprintf("%s/webhooks/%s?embed=1", baseURL, os.Getenv("WEBHOOK_ID"))
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(message))
	if err != nil {
		return fmt.Errorf("failed to create a new webhook: %v", err)
	}
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	req.Header.Set("X-TRAQ-Signature", calcHMACSHA1(message, webhookSecret))
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	httpClient := http.DefaultClient
	_, err = httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post a request: %v", err)
	}

	return nil
}
