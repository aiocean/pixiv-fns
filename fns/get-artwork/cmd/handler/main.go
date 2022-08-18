package main

import (
	"context"

	pixivv1 "github.com/aiocean/go-sdk/pixiv/v1"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"pkg.aiocean.io/pixiv/get-artwork/internal"
)

func Handler(ctx context.Context, request *pixivv1.GetArtworkRequest) (*pixivv1.GetArtworkResponse, error) {

	artwork, err := internal.FetchWork(request.ArtworkId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to fetch artwork")
	}

	return &pixivv1.GetArtworkResponse{
		Artwork: artwork,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
