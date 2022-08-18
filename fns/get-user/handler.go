package handler

import (
	"net/http"

	"github.com/aiocean/cfutil"
	"github.com/aiocean/get-user/internal"
	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
	"google.golang.org/protobuf/proto"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfutil.ProtobufHandler(w, r, &pixivv1.GetUserRequest{}, getUser)
}

func getUser(request proto.Message) (proto.Message, error) {
	req := request.(*pixivv1.GetUserRequest)

	artwork, err := internal.FetchUser(req.GetUserId())
	if err != nil {
		return nil, err
	}

	response := &pixivv1.GetUserResponse{
		User: artwork,
	}

	return response, nil
}
