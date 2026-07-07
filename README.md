# JPEG Encoder for Golang

[![Go Reference](https://pkg.go.dev/badge/github.com/nickalie/go-mozjpegbin.svg)](https://pkg.go.dev/github.com/nickalie/go-mozjpegbin)
[![CI](https://github.com/nickalie/go-mozjpegbin/actions/workflows/ci.yml/badge.svg)](https://github.com/nickalie/go-mozjpegbin/actions/workflows/ci.yml)

MozJPEG Encoder for Golang. Wraps the `cjpeg` and `jpegtran` command line tools.

## Install

```go get github.com/nickalie/go-mozjpegbin```

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

Under the hood the library shells out to the *cjpeg* and *jpegtran* command line
tools. They must be available on `PATH` (or in the directory set via
`mozjpegbin.Dest(...)`). Install them from mozjpeg or a compatible
libjpeg-turbo build (e.g. `apt-get install libjpeg-turbo-progs`).

> Note: `mozjpegbin.SkipDownload()` is now a deprecated no-op — automatic binary
> download is no longer supported; install the tools on your system.

Build mozjpeg for your target platform so you get its better compression
(trellis quantization and other optimizations that stock libjpeg-turbo lacks).
Since mozjpeg 4.x the build uses CMake:

```
apk add --no-cache build-base cmake nasm wget   # Debian/Ubuntu: cmake nasm g++ make wget

wget -O mozjpeg.tar.gz https://github.com/mozilla/mozjpeg/archive/refs/tags/v4.1.1.tar.gz && \
tar -xzf mozjpeg.tar.gz && rm mozjpeg.tar.gz && \
cd mozjpeg-4.1.1 && \
cmake -G"Unix Makefiles" -DCMAKE_POLICY_VERSION_MINIMUM=3.5 -DPNG_SUPPORTED=0 -DCMAKE_INSTALL_PREFIX=/opt/mozjpeg . && \
make -j"$(nproc)" && make install && \
cd / && rm -rf mozjpeg-4.1.1 && \
ln -s /opt/mozjpeg/bin/cjpeg /usr/local/bin/cjpeg && \
ln -s /opt/mozjpeg/bin/jpegtran /usr/local/bin/jpegtran
```

Ready-to-use [`docker/`](docker) images build mozjpeg this way. If you only need
the tools working (e.g. for tests) and don't care about the extra compression,
the system libjpeg-turbo package also provides compatible `cjpeg`/`jpegtran`
(`apt-get install libjpeg-turbo-progs` / `apk add libjpeg-turbo-utils`).
