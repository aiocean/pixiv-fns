package handler

import (
	"net/http"

	"github.com/aiocean/cfutil"
	"github.com/aiocean/get-artwork/internal"
	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
	"google.golang.org/protobuf/proto"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfutil.ProtobufHandler(w, r, &pixivv1.GetArtworkRequest{}, getArtwork)
}

func getArtwork(request proto.Message) (proto.Message, error) {
	req := request.(*pixivv1.GetArtworkRequest)

	artwork, err := internal.FetchWork(req.GetArtworkId())
	if err != nil {
		return nil, err
	}

	response := &pixivv1.GetArtworkResponse{
		Artwork: artwork,
	}

	return response, nil
}
