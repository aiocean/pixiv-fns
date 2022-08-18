package internal

import (
	"io"
	"log"
	"net/http"
	"time"

	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
	"github.com/tidwall/gjson"
)

func FetchUser(id string) (user *pixivv1.User, err error) {

	apiEndpoint := "https://www.pixiv.net/ajax/user/" + id + "?full=1&lang=en"

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

	user = &pixivv1.User{}
	user.Id = body.Get("userId").String()
	user.Name = body.Get("name").String()
	user.AvatarUrl = body.Get("imageBig").String()

	return user, err
}
