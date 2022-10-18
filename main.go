package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/qeesung/image2ascii/convert"
)

var (
	client  http.Client
	options convert.Options
)

func init() {
	color := flag.Bool("c", false, "colored image")
	width := flag.Int("w", 120, "width of the image")
	height := flag.Int("h", 40, "height of the image")
	fullscreen := flag.Bool("f", false, "add padding if necessary")

	flag.Parse()

	options.Colored = *color
	options.FixedWidth = *width
	options.FixedHeight = *height
	options.FitScreen = *fullscreen
}

func catURL() (string, error) {
	const url = "https://api.thecatapi.com/v1/images/search?mime_types=jpg"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error couldn't be constructed: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error in request: %v", err)
	}

	defer resp.Body.Close()

	var body []struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil || len(body) <= 0 {
		return "", fmt.Errorf("unable to decode response body: %v", err)
	}

	return body[0].URL, nil
}

func catImage(url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build image request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in request: %v", err)
	}

	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error decoding image body: %v", err)
	}

	return img, nil
}

func main() {
	converter := convert.NewImageConverter()
	url, err := catURL()
	if err != nil {
		log.Fatalf("unable to fetch cat url: %v", err)
	}

	img, err := catImage(url)
	if err != nil {
		log.Fatalf("unable to fetch cat image: %v", err)
	}

	fmt.Print(converter.Image2ASCIIString(img, &options))
}
