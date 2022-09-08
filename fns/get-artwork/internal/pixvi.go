package internal

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
	"github.com/tidwall/gjson"
)

func FetchWork(id string) (artwork *pixivv1.Artwork, err error) {

	apiEndpoint := "https://www.pixiv.net/ajax/illust/" + id + "?lang=en"

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
		if strings.Contains(result.String(), "bookmarks") || strings.Contains(result.String(), "original") {
			continue
		}
		tags := strings.Split(result.String(), ",")
		for _, tagName := range tags {
			tagName = strings.TrimSpace(tagName)
			// check if result is exist in artwork.Tags
			isExist := false
			for _, tag := range artwork.Tags {
				if tag.GetId() == tagName {
					isExist = true
					break
				}
			}
			if isExist {
				continue
			}
			artwork.Tags = append(artwork.Tags, &pixivv1.Tag{
				Name: tagName,
				Id:   tagName,
			})
		}

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

	imageSize, err := getImageSize(firstUrl.Original)
	if err != nil {
		return nil, err
	}

	firstUrl.OriginalSize = imageSize / 1000000

	artwork.ImageUrls = append(artwork.ImageUrls, firstUrl)
	artwork.AgeLimit = ageLimit
	artwork.UserId = body.Get("userId").String()
	artwork.LikeCount = body.Get("likeCount").Int()
	artwork.ViewCount = body.Get("viewCount").Int()
	artwork.CommentCount = body.Get("commentCount").Int()
	artwork.BookmarkCount = body.Get("bookmarkCount").Int()
	return artwork, err
}

func getImageSize(imageUrl string) (float32, error) {
	client := &http.Client{
		Timeout: time.Minute * 4,
	}

	req, err := http.NewRequest("GET", imageUrl, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Referer", "https://www.pixiv.net")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	contentLength := resp.Header.Get("content-length")
	imageSize, err := strconv.ParseFloat(contentLength, 64)
	if err != nil {
		return 0, err
	}

	return float32(imageSize), nil
}
