package mozjpegbin

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/nickalie/go-binwrapper"
)

var dest = ""

// SkipDownload is deprecated and has no effect.
//
// Deprecated: binary download is no longer supported by go-binwrapper. Install
// cjpeg/jpegtran on the system (see README) so they are resolved from PATH or
// from the directory set via Dest.
func SkipDownload() {}

// Dest sets the directory to look for mozjpeg binaries in.
// By default binaries are resolved from PATH.
func Dest(value string) {
	dest = value
}

func createBinWrapper() *binwrapper.BinWrapper {
	return binwrapper.NewBinWrapper().AutoExe().Dest(dest)
}

// createReaderFromImage encodes img as a binary PPM (P6) stream. PPM is cjpeg's
// native input format, so this avoids an extra lossy JPEG pass and works with
// any cjpeg build (unlike PNG/JPEG input, which older builds don't accept).
func createReaderFromImage(img image.Image) (io.Reader, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var buffer bytes.Buffer
	if _, err := fmt.Fprintf(&buffer, "P6\n%d %d\n255\n", width, height); err != nil {
		return nil, err
	}

	row := make([]byte, width*3)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		i := 0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			row[i] = byte(r >> 8)
			row[i+1] = byte(g >> 8)
			row[i+2] = byte(b >> 8)
			i += 3
		}
		if _, err := buffer.Write(row); err != nil {
			return nil, err
		}
	}

	return &buffer, nil
}

func version(b *binwrapper.BinWrapper) (string, error) {
	b.Reset()
	err := b.Run("-version")

	if err != nil {
		return "", err
	}

	v := string(b.StdErr())
	v = strings.ReplaceAll(v, "\n", "")
	v = strings.ReplaceAll(v, "\r", "")
	return v, nil
}
