package wain

import (
	"fmt"
	"gopkg.in/h2non/bimg.v0"
	"regexp"
	"strings"
)

type ResizeOptions struct {
	Width  int
	Height int
	Format string
	Params map[string]string
}

func HandleProcessing(url ConfigUrl, s3 map[string]*S3Connection, r ResizeOptions) ([]byte, error) {

	cacheKey := url.Cache.Key
	originalKey := url.Original.Key
	for name, value := range r.Params {
		cacheKey = strings.Replace(cacheKey, fmt.Sprintf("{%s}", name), value, -1)
		originalKey = strings.Replace(originalKey, fmt.Sprintf("{%s}", name), value, -1)
	}
	unusedVariables := regexp.MustCompile("{.+?}")
	cacheKey = unusedVariables.ReplaceAllString(cacheKey, "")
	originalKey = unusedVariables.ReplaceAllString(originalKey, "")

	fmt.Printf("Checking cached version %s\n", cacheKey)
	imageData, err := s3[url.Cache.Bucket].GetObject(cacheKey)
	if err == nil {
		return imageData, nil
	}

	fmt.Printf("No cached version found\n")
	fmt.Printf("Downloading %s\n", originalKey)

	imageData, err = s3[url.Original.Bucket].GetObject(originalKey)

	fmt.Printf("Resizing...\n")

	options := bimg.Options{
		Height:  r.Height,
		Width:   r.Width,
		Quality: 100,
	}

	imageData, err = bimg.NewImage(imageData).Process(options)

	fmt.Printf("Save cached version\n")
	go s3[url.Cache.Bucket].PutObject(cacheKey, imageData, "image/jpeg")

	return imageData, err
}
