package mozjpegbin

import (
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	generateSources()
}

// generateSources writes deterministic test fixtures locally so the suite stays
// hermetic (no flaky external download). It produces two 512x512 fixtures:
//   - source.jpg: a JPEG, consumed by jpegtran (which transforms JPEGs) and
//     decoded back to an image by the InputImage tests.
//   - source.ppm: a binary PPM, cjpeg's native input format. cjpeg cannot read
//     JPEG, so the file/reader-based cjpeg tests feed it this instead.
func generateSources() {
	img := makeTestImage(512, 512)

	writeFixture("source.jpg", func(f *os.File) error {
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	})

	writeFixture("source.ppm", func(f *os.File) error {
		r, err := createReaderFromImage(img)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, r)
		return err
	})
}

func makeTestImage(width, height int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return img
}

// writeFixture creates path and runs encode against it, unless it already exists.
func writeFixture(path string, encode func(*os.File) error) {
	if _, err := os.Stat(path); err == nil {
		return
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if err := encode(f); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func TestEncodeImage(t *testing.T) {
	c := NewCJpeg()
	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	img, err := jpeg.Decode(f)
	assert.Nil(t, err)
	c.InputImage(img)
	c.OutputFile("target.jpg")
	err = c.Run()
	assert.Nil(t, err)
	validateJpg(t)
}

func TestEncodeImage2(t *testing.T) {

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

	c := NewCJpeg()
	c.InputImage(img)
	c.OutputFile("target.jpg")
	err := c.Run()
	assert.Nil(t, err)
	validateJpgImage(t, img)
}

func TestEncodeReader(t *testing.T) {
	c := NewCJpeg()
	f, err := os.Open("source.ppm")
	assert.Nil(t, err)
	c.Input(f)
	c.OutputFile("target.jpg")
	err = c.Run()
	assert.Nil(t, err)
	validateJpg(t)
}

func TestEncodeFile(t *testing.T) {
	c := NewCJpeg()
	c.Quality(100)
	c.Optimize(true)
	c.InputFile("source.ppm")
	c.OutputFile("target.jpg")
	err := c.Run()
	assert.Nil(t, err)
	validateJpg(t)
}

func TestEncodeWriter(t *testing.T) {
	f, err := os.Create("target.jpg")
	assert.Nil(t, err)
	defer f.Close() //nolint:errcheck // backup close; validateJpg fails if the file is incomplete

	c := NewCJpeg()
	c.InputFile("source.ppm")
	c.Output(f)
	err = c.Run()
	assert.Nil(t, err)
	f.Close() //nolint:errcheck // validateJpg re-reads the file and would fail on a bad write
	validateJpg(t)
}

func TestCJpegVersion(t *testing.T) {
	v, err := NewCJpeg().Version()
	assert.Nil(t, err)
	assert.NotEmpty(t, v)
}

func validateJpg(t *testing.T) {
	//defer os.Remove("target.jpg")
	fSource, err := os.Open("source.jpg")
	assert.Nil(t, err)
	imgSource, err := jpeg.Decode(fSource)
	assert.Nil(t, err)
	validateJpgImage(t, imgSource)
}

func validateJpgImage(t *testing.T, imgSource image.Image) {
	//defer os.Remove("target.jpg")
	fTarget, err := os.Open("target.jpg")
	assert.Nil(t, err)
	defer fTarget.Close() //nolint:errcheck // read-only file, close error irrelevant
	imgTarget, err := jpeg.Decode(fTarget)
	assert.Nil(t, err)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}
