package dto

import (
	"encoding/json"
	"errors"
	"fmt"
)

type MediaResponse struct {
	Albums Albums `json:"albums"`
}

func (m *MediaResponse) GetSmallestImage() (Image, error) {
	var album Album
	if len(m.Albums.Items) > 0 {
		album = m.Albums.Items[0]
		image, err := album.GetSmallestImage()
		if err != nil {
			return image, err
		}
		return image, nil
	}
	return Image{}, errors.New("no image found")
}

func (s *MediaResponse) Hydrate(response []byte) {
	err := json.Unmarshal(response, &s)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
}

type Albums struct {
	Items []Album `json:"items"`
}

type Album struct {
	AlbumType string   `json:"album_type"`
	Artists   []Artist `json:"artists"`
	Images    []Image  `json:"images"`
	Name      string   `json:"name"`
}

func (a *Album) GetSmallestImage() (Image, error) {
	if len(a.Images) == 0 {
		return Image{}, errors.New("no images available")
	}
	smallestImage := a.Images[0]
	for _, image := range a.Images {
		if image.Height < smallestImage.Height {
			smallestImage = image
		}
	}
	return smallestImage, nil
}

type Artist struct {
	Name string `json:"name"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
