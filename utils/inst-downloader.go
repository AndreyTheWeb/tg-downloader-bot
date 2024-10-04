package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type InstagramResponse struct {
	Links []struct {
		Quality string `json:"quality"`
		Link    string `json:"link"`
	} `json:"links"`
}

func FetchInstagramVideo(url string, apiKey string) (string, error) {
	apiUrl := fmt.Sprintf("https://social-media-video-downloader.p.rapidapi.com/smvd/get/instagram?url=%s", url)

	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "social-media-video-downloader.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("не удалось получить видео, код ошибки: " + resp.Status)
	}

	var result InstagramResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, link := range result.Links {
		fmt.Println(123, link.Quality, link.Link)

		if strings.Contains(strings.ToLower(link.Quality), "video") {
			return link.Link, nil
		}
	}

	return "", errors.New("видео не найдено")
}
