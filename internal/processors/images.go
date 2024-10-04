package processors

import (
	"fmt"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"

	"go-flattern-prices/internal/logger"
)

const (
	fontDir = "fonts"
)

func DrawTagsOverImage(path string, ext string, tags []*Tag, outDir string, inputDir string) error {
	im, err := gg.LoadImage(fmt.Sprintf("%s%s", path, ext))
	if err != nil {
		logger.Error("Ошибка чтения файла изображения: %v", err)
		return err
	}

	W := im.Bounds().Max.X
	H := im.Bounds().Max.Y
	dc := gg.NewContext(W, H)

	dc.DrawImage(im, 0, 0)

	for _, tag := range tags {
		dc.SetColor(tag.Color)
		fontFile := filepath.Join(fontDir, fmt.Sprintf("%s.ttf", tag.Font))

		if e := dc.LoadFontFace(fontFile, float64(tag.Size)); e != nil {
			logger.Error("Не найден файл шрифта: %s, error: %v", fontFile, err)
			return e
		}

		var price string
		switch tag.Price {
		case nil:
			price = "-"
		default:
			price = strconv.Itoa(*tag.Price)
		}

		var ax float64
		switch tag.Align {
		case "right":
			ax = 1.0
		case "center":
			ax = 0.5
		default:
			ax = 0.0
		}
		dc.DrawStringAnchored(price, float64(tag.X), float64(tag.Y), ax, 1.1)
	}

	if ext == pngExt {
		outputPath := fmt.Sprintf("%s%s", path, ext)
		outputPath = strings.Replace(outputPath, inputDir, outDir, 1)

		e := os.MkdirAll(filepath.Dir(outputPath), 0755)
		if e != nil {
			logger.Error("Ошибка создания директории: %v", e)
			return e
		}

		e = dc.SavePNG(outputPath)
		if e != nil {
			logger.Error("Ошибка создания файла изображения: %v", e)
			return e
		}
	} else if ext == jpgExt || ext == jpegExt {
		outputPath := fmt.Sprintf("%s%s", path, ext)
		outputPath = strings.Replace(outputPath, inputDir, outDir, 1)

		e := os.MkdirAll(filepath.Dir(outputPath), 0755)
		if e != nil {
			logger.Error("Ошибка создания директории: %v", e)
			return e
		}

		toimg, e := os.Create(outputPath)
		if e != nil {
			logger.Error("Ошибка создания файла изображения: %v", e)
			return e
		}
		defer toimg.Close()

		var opts jpeg.Options
		opts.Quality = 90
		e = jpeg.Encode(toimg, dc.Image(), &opts)
		if e != nil {
			logger.Error("Ошибка записи файла изображения: %v", e)
			return e
		}
	}

	return nil
}
