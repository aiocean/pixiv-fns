# Pixiv Functions

The main purpose of this library is to provide a sample implementation of the Pixiv API for other developers to use as a reference. It is not intended to be a complete library for interacting with the Pixiv API.

## Features

- [x] Support both JSON and Protobuf binary request/response
- [x] Reduce boilerplate code, Now you only need to define the business function, no need to marshal/unmarshal the request/response

```go
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
```

## Configuration

Create new file prop.auto.tfvars and add the following:

```hcl
gcp_project_id="wibuzone"
gcp_region="asia-southeast1"
```

## Installation

```shell
terraform init
```

## Deployment

```shell
terraform apply
```

