package web

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	xdraw "golang.org/x/image/draw"
)

type photoCache struct {
	dir string
	ttl time.Duration
}

func (c *photoCache) load(username string, size int, maxAge *time.Duration) ([]byte, bool, error) {
	if c.ttl == 0 {
		return nil, false, nil
	}

	path := c.path(username, size)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}

		return nil, false, err
	}

	age := time.Since(info.ModTime())
	limit := c.ttl
	if maxAge != nil && *maxAge < limit {
		limit = *maxAge
	}

	if age > limit {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return nil, false, err
		}
		return nil, false, nil
	}

	photo, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return photo, true, nil
}

func (c *photoCache) store(username string, size int, photo []byte) error {
	if c.ttl == 0 {
		return nil
	}

	path := c.path(username, size)
	tmp, err := os.CreateTemp(c.dir, "photo-*.jpg")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	if _, err := tmp.Write(photo); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), path)
}

func (c *photoCache) path(username string, size int) string {
	return filepath.Join(c.dir, username+"_"+strconv.Itoa(size)+".jpg")
}

func resizePhoto(photo []byte, size int) ([]byte, error) {
	if size == 0 {
		return photo, nil
	}

	source, _, err := image.Decode(bytes.NewReader(photo))
	if err != nil {
		return nil, err
	}

	srcBounds := source.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	if srcW == 0 || srcH == 0 {
		return nil, errors.New("invalid photo dimensions")
	}

	scale := math.Max(float64(size)/float64(srcW), float64(size)/float64(srcH))
	dstW := int(math.Ceil(float64(srcW) * scale))
	dstH := int(math.Ceil(float64(srcH) * scale))
	if dstW < size {
		dstW = size
	}
	if dstH < size {
		dstH = size
	}

	scaled := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), source, srcBounds, xdraw.Over, nil)

	offsetX := (dstW - size) / 2
	offsetY := (dstH - size) / 2
	cropped := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(cropped, cropped.Bounds(), scaled, image.Point{X: offsetX, Y: offsetY}, draw.Src)

	var out bytes.Buffer
	if err := jpeg.Encode(&out, cropped, &jpeg.Options{Quality: 85}); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
