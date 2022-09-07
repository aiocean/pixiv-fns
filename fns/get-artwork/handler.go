package handler

import (
	"context"
	"net/http"

	"github.com/aiocean/get-artwork/internal"
	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	pixivv1.GetArtworkHandler(w, r, getArtwork)
}

func getArtwork(ctx context.Context, request *pixivv1.GetArtworkRequest) (*pixivv1.GetArtworkResponse, error) {

	artwork, err := internal.FetchWork(request.GetArtworkId())
	if err != nil {
		return nil, err
	}

	response := &pixivv1.GetArtworkResponse{
		Artwork: artwork,
	}

	return response, nil
}
