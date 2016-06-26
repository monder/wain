package wain

import (
	"fmt"
	. "github.com/tj/go-debug"
	"regexp"
	"strings"
)

var debug = Debug("handler")

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

	imageData, err := s3[url.Cache.Bucket].GetObject(cacheKey)
	if err == nil {
		return imageData, nil
	}

	debug("No cached version found %s\n", cacheKey)
	debug("Downloading %s\n", originalKey)

	imageData, err = s3[url.Original.Bucket].GetObject(originalKey)
	if err != nil {
		return nil, err
	}

	debug("Resizing...\n")

	imageData, err = VipsResize(imageData, r)

	if err == nil {
		debug("Save cached version\n")
		go func() {
			e := s3[url.Cache.Bucket].PutObject(cacheKey, imageData, "image/jpeg")
			if e != nil {
				fmt.Println(e)
			}
		}()
	}

	return imageData, err
}
