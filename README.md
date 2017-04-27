# JPEG Encoder for Golang

[![](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/nickalie/go-mozjpegbin)
[![](https://circleci.com/gh/nickalie/go-mozjpegbin.png?circle-token=bf2a63a9ecd6ca6f4c4d81028d710cb847e58695)](https://circleci.com/gh/nickalie/go-mozjpegbin)

MozJPEG Encoder for Golang based on unofficial mozjpeg distribution

## Install

```go get -u github.com/nickalie/go-mozjpegbin```

## Example of usage

```
package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"github.com/nickalie/go-mozjpegbin"
)

func main() {
	const width, height = 256, 256

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	f, err := os.Create("image.jpg")
	if err != nil {
		log.Fatal(err)
	}

	if err := mozjpegbin.Encode(f, img, nil); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
```

## CJpeg

CJpeg is a wrapper for *cjpeg* command line tool.

Example to convert image.pgm to image.jpg:

```
err := mozjpegbin.NewCJpeg().
		Quality(70).
		InputFile("image.pgm").
		OutputFile("image_min.jpg").
		Run()
```

## JpegTran

JpegTran is a wrapper for *jpegtran* command line tool.

Example to optimize and crop image.jpg to image_opt.jpg:

```
err := mozjpegbin.NewJpegTran().
		InputFile("image.jpg").
		Crop(10, 5, 100, 100)
		OutputFile("image_opt.jpg").
		Run()
```

## mozjpeg distribution

Under the hood library uses *cjpeg* and *jpegtrans* command line tools from mozjpeg. To avoid compatibility issues, it's better to build mozjpeg for your target platform and call ```mozjpegbin.SkipDownload()``` to avoid using of prebuilt binaries 

Snippet to build mozjpeg on alpine:

```
apk add --no-cache --update g++ make autoconf automake libtool nasm wget

wget https://github.com/mozilla/mozjpeg/releases/download/v3.2-pre/mozjpeg-3.2-release-source.tar.gz && \
tar -xvzf mozjpeg-3.2-release-source.tar.gz && \
rm mozjpeg-3.2-release-source.tar.gz && \
cd mozjpeg && \
./configure && \
make install && \
cd / && rm -rf mozjpeg && \
ln -s /opt/mozjpeg/bin/jpegtran /usr/local/bin/jpegtran && \
ln -s /opt/mozjpeg/bin/cjpeg /usr/local/bin/cjpeg
```
