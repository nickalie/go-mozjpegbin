package mozjpegbin

import (
	"bytes"
	"runtime"
	"io/ioutil"
	"github.com/nickalie/go-binwrapper"
	"image"
	"io"
	"image/jpeg"
	"strings"
)

var skipDownload bool
var dest string = "vendor/mozjpeg"

//Detects platforms without prebuilt binaries (alpine and arm).
//For this platforms mozjpeg tools should be built manually.
func DetectUnsupportedPlatforms()  {
	if runtime.GOARCH == "arm" {
		SkipDownload()
	} else if runtime.GOOS == "linux" {
		output, err := ioutil.ReadFile("/etc/issue")

		if err == nil && bytes.Contains(bytes.ToLower(output), []byte("alpine")) {
			SkipDownload()
		}
	}
}

//Skips binary download.
func SkipDownload() {
	skipDownload = true
	dest = ""
}

//Sets directory to download mozjpeg binaries or where to look for them if SkipDownload is used.
func Dest(value string) {
	dest = value
}

func createBinWrapper() *binwrapper.BinWrapper {
	b := binwrapper.NewBinWrapper().AutoExe()

	if !skipDownload {
		b.Src(
			binwrapper.NewSrc().
				Url("https://raw.githubusercontent.com/imagemin/mozjpeg-bin/macos/cjpeg").
				Os("darwin")).
			Src(
			binwrapper.NewSrc().
				Url("https://raw.githubusercontent.com/imagemin/mozjpeg-bin/linux/cjpeg").
				Os("linux")).
			Src(
			binwrapper.NewSrc().
				Url("https://mozjpeg.codelove.de/bin/mozjpeg_3.1_x86.zip").
				Os("win32"))
	}

	return b.Strip(2).Dest(dest)
}

func createReaderFromImage(img image.Image) (io.Reader, error) {
	var buffer bytes.Buffer
	err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality:100})
	return &buffer, err
}

func version(b *binwrapper.BinWrapper) (string, error) {
	b.Reset()
	err := b.Run("-version")

	if err != nil {
		return "", err
	}

	v := string(b.StdErr())
	v = strings.Replace(v, "\n", "", -1)
	v = strings.Replace(v, "\r", "", -1)
	return v, nil
}