# Wain

[![license](https://img.shields.io/github/license/monder/wain.svg?maxAge=2592000&style=flat-square)]()
[![GitHub tag](https://img.shields.io/github/tag/monder/wain.svg?style=flat-square)]()

Wain is a http proxy that allows to dynamically resize and serve images from an s3 bucket.

## Installation

### Installing libvips on MacOS
```
brew install homebrew/science/vips --with-mozjpeg
```

## Docker image
```
docker run -v /home/user/config.yaml:/config.yaml -p 3000:3000 monder/wain /config.yaml
```

## Configuration

Wain allows to configure routes in `yaml` format.

```yaml
# vim: ts=2:sw=2
port: 3000
buckets:
- name: wain-demo-content
  region: eu-west-1
  accessKey: AKIAJMK6JWZF4F4JBSXQ
  accessSecret: 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
urls:
- pattern: /images/{id:IMG_\d+}/{width:[0-9]+}{height:(x[0-9]+)?}.jpg
  original:
    bucket: wain-demo-content
    key: original/{id}.JPG
  cache:
    bucket: wain-demo-content
    key: cache/{id}/{width}{height}.jpg
```

The configuration above will process the urls only matching the pattern `/images/{id:IMG_\d+}/{width:[0-9]+}{height:(x[0-9]+)?}.jpg` capturing the variables.

Note thate the `height` parameter is optional (`?`) and the pattern will match both:
`/images/IMG_2530/100.jpg`
and
`/images/IMG_2530/100x150.jpg`

The file will be downloaded from the bucket specified in `original`. `bucket` name is looked up in `buckets` section. `key` should be the file name inside the bucket. All variables captured in `pattern` could be used here to generate the key. All unknown variables are replaced with an empty string (`""`).

Requesting the file `/images/IMG_2530/100.jpg` will download and resize file `wain-demo-content/original/IMG_2530.JPG` to width `100px` and some height keeping the aspect ratio of the original.

`100`:
![100](https://cloud.githubusercontent.com/assets/232147/16442096/52763c16-3dd8-11e6-8a27-6e491ac7cd91.jpg)

Specifying the `height` as well as `width` will fit the original, adding a blurred background:

`100x150`:
![100x150](https://cloud.githubusercontent.com/assets/232147/16442132/9ebb4b2a-3dd8-11e6-8e0a-0fadda7c8c0a.jpg)

`150x100`:
![150x100](https://cloud.githubusercontent.com/assets/232147/16442133/a0502690-3dd8-11e6-835a-80526f8f71c9.jpg)
