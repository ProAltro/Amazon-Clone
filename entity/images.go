package entity

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func GetImages(images []string) ([]string, error) {
	image_urls := []string{}
	for _, image := range images {
		image_url, err := getImage(image)
		if err != nil {
			return nil, err
		}
		image_urls = append(image_urls, image_url)
	}
	return image_urls, nil
}

func getImage(image string) (string, error) {
	filepath := "./images/" + image + ".jpg"
	_, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return filepath, nil
}

func CreateImages(images []string, uri string) ([]string, error) {
	image_urls := []string{}
	for i, image := range images {
		image_url := uri + "_" + strconv.Itoa(i)
		err := createImage(image, uri+"_"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		image_urls = append(image_urls, image_url)
	}
	return image_urls, nil
}

func createImage(image string, uri string) error {

	decodedImg, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return err
	}

	filepath := "./images/" + uri + ".jpg"
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()
	fmt.Println("File created")
	_, err = file.Write(decodedImg)
	if err != nil {
		return err
	}
	return nil
}

func DeleteImages(images []string) error {
	for _, image := range images {
		filepath := "./images/" + image + ".jpg"
		err := os.Remove(filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

func JSON_To_Image(images []byte) ([]string, error) {
	var prods []string
	err := json.Unmarshal(images, &prods)
	if err != nil {
		return nil, err
	}
	return prods, nil
}

func Images_To_JSON(images []string) ([]byte, error) {
	return json.Marshal(images)
}
