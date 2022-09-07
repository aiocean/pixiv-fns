package handler

import (
	"context"
	"net/http"

	"github.com/aiocean/get-user/internal"
	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	pixivv1.GetUserHandler(w, r, getUser)
}

func getUser(ctx context.Context, request *pixivv1.GetUserRequest) (*pixivv1.GetUserResponse, error) {

	artwork, err := internal.FetchUser(request.GetUserId())
	if err != nil {
		return nil, err
	}

	response := &pixivv1.GetUserResponse{
		User: artwork,
	}

	return response, nil
}
