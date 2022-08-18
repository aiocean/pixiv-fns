package internal

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
	"github.com/tidwall/gjson"
)

func FetchWork(id string) (artwork *pixivv1.Artwork, err error) {

	apiEndpoint := "https://www.pixiv.net/ajax/illust/" + id

	client := &http.Client{
		Timeout: time.Minute * 4,
	}

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://www.pixiv.net")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body := gjson.ParseBytes(data).Get("body")

	var ageLimit = "all-age"
	for _, tag := range body.Get("tags.tags.#.tag").Array() {
		if tag.Str == "R-18" {
			ageLimit = "r18"
			break
		}
	}

	artwork = &pixivv1.Artwork{}
	artwork.Id = body.Get("illustId").String()
	artwork.Title = strings.Split(body.Get("alt").String(), "/")[0]

	tagsResult := body.Get("tags.tags.#.translation.en").Array()
	for _, result := range tagsResult {
		if strings.Contains(result.String(), "bookmarks") {
			continue
		}
		artwork.Tags = append(artwork.Tags, result.String())
	}

	urls := body.Get("urls")

	firstUrl := &pixivv1.ImageUrl{
		Thumb:    urls.Get("thumb").String(),
		Mini:     urls.Get("mini").String(),
		Small:    urls.Get("small").String(),
		Regular:  urls.Get("regular").String(),
		Original: urls.Get("original").String(),
		Width:    uint32(body.Get("width").Uint()),
		Height:   uint32(body.Get("height").Uint()),
	}
	artwork.ImageUrls = append(artwork.ImageUrls, firstUrl)
	artwork.AgeLimit = ageLimit
	artwork.UserId = body.Get("userId").String()
	return artwork, err
}
