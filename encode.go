package mozjpegbin

import (
	"image"
	"io"
)

// Options to use with encoder
type Options struct {
	Quality  uint
	Optimize bool
}

// Encode encodes image.Image into jpeg using cjpeg.
func Encode(w io.Writer, m image.Image, o *Options) error {
	cjpeg := NewCJpeg()

	if o != nil {
		cjpeg.Quality(o.Quality)
		cjpeg.Optimize(o.Optimize)
	}

	return cjpeg.InputImage(m).Output(w).Run()
}
