package externalserv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ImageRequester interface {
	GetImage() ([]byte, error)
}

const (
	imageApiAddress = "https://api.unsplash.com/"
	accessKey       = "cfcffe9f877092173c2397a9eb9e8579569ebb08d91992b83af7f99294fb9453"
	secretKey       = "91eb588d1fb07c5a2276a7819d4c58fae5e581729285c1d12ab8996a6fc6605f"
	randomImageURL  = "/photos/random"
)

var imagesGot = make(map[string]ImageUnsplashResponse)

type imageRequesterUnsplash struct {
	apiAddress string
}

func NewImageRequesterUnsplash() ImageRequester {
	return &imageRequesterUnsplash{}
}

func (i *imageRequesterUnsplash) GetImage() ([]byte, error) {
	var err error
	resp, err := i.getRandomPhoto()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var imgResponse ImageUnsplashResponse
	err = json.NewDecoder(resp.Body).Decode(&imgResponse)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	imagesGot[imgResponse.ID] = imgResponse
	WriteToFileAllJSONImages(imagesGot)

	return nil, nil
}

func (i *imageRequesterUnsplash) setupRequest(urlString string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, imageApiAddress+urlString, nil)
	urlValues := url.Values{}
	urlValues.Set("client_id", accessKey)
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID: %v", accessKey))
	req.Header.Set("Accept-Version", "v1")
	req.URL.RawQuery = urlValues.Encode()
	return req, err
}

func (i *imageRequesterUnsplash) getRandomPhoto() (*http.Response, error) {
	req, err := i.setupRequest(randomImageURL)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func WriteToFileAllJSONImages(jsons map[string]ImageUnsplashResponse) {
	imagesBytes, err := ioutil.ReadFile("images.json")
	if err != nil {
		log.Printf("can't read %v\n", err)
		return
	}

	readJsons := make(map[string]ImageUnsplashResponse)
	err = json.Unmarshal(imagesBytes, &readJsons)
	if err != nil {
		log.Printf("can't unmarshall %v\n", err)
		return
	}

	for key, value := range readJsons {
		if _, ok := jsons[key]; !ok {
			jsons[key] = value
		}
	}

	b, err := json.MarshalIndent(jsons, "", "\t")
	if err != nil {
		log.Printf("can't parse %v\n", err)
		return
	}

	err = ioutil.WriteFile("images.json", b, 0644)
	if err != nil {
		log.Printf("can't write %v\n", err)
		return
	}
}

type ImageUnsplashResponse struct {
	ID     string     `json:"id"`
	Width  int32      `json:"width"`
	Height int32      `json:"height"`
	URLs   ImageURLs  `json:"urls"`
	Links  ImageLinks `json:"links"`
}

type ImageURLs struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
}

type ImageLinks struct {
	Self             string `json:"self"`
	Html             string `json:"html"`
	Download         string `json:"download"`
	DownloadLocation string `json:"download_location"`
}
