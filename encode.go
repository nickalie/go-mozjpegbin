package mozjpegbin

import (
	"image"
	"io"
)

type Options struct{
	Quality uint
	Optimize bool
}


func Encode(w io.Writer, m image.Image, o *Options) error {
	cjpeg := NewCJpeg()

	if o != nil {
		cjpeg.Quality(o.Quality)
		cjpeg.Optimize(o.Optimize)
	}

	return cjpeg.InputImage(m).Output(w).Run()
}
