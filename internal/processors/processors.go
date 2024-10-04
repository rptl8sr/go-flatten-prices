package processors

import (
	"errors"
	"fmt"
	"go-flatten-prices/internal/logger"
	"os"
	"path/filepath"
	"strings"
)

const (
	jpegExt = ".jpeg"
	jpgExt  = ".jpg"
	pngExt  = ".png"
	csvExt  = ".csv"
)

var ErrNoImageFiles = errors.New("файлы изображений не найдены")

type ErrNotSameLength struct {
	csv int
	img int
}

func (e *ErrNotSameLength) Error() string {
	return fmt.Errorf("разное количество файлов csv (%d) и изображений (%d)", e.csv, e.img).Error()
}

type ErrNoSameFile struct {
	filename string
	ext      string
}

func (e *ErrNoSameFile) Error() string {
	return fmt.Errorf("нет соответствующего csv-файла %s%s", e.filename, e.ext).Error()
}

type ErrNoNeededDate struct {
	filename string
	date     string
}

func (e *ErrNoNeededDate) Error() string {
	return fmt.Errorf("дата в файле %s не соответствует заданной %s, он будет пропущен", e.filename, e.date).Error()
}

type processor struct {
	path string
}

type Processor interface {
	GetFilesPairs() (map[string]string, error)
	CheckDate(files map[string]string, date string)
}

func New(path string) Processor {
	return &processor{path: path}
}

func (p *processor) GetFilesPairs() (map[string]string, error) {
	images := make(map[string]string)
	prices := make(map[string]struct{})
	res := make(map[string]string)

	// get files tree
	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case jpegExt, jpgExt, pngExt:
				withoutExt := strings.TrimSuffix(path, ext)
				images[withoutExt] = ext
			case csvExt:
				withoutExt := strings.TrimSuffix(path, ext)
				prices[withoutExt] = struct{}{}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, ErrNoImageFiles
	}

	logger.Info(fmt.Sprintf("Изображения: %d, Цены: %d", len(images), len(prices)))
	if len(images) != len(prices) {
		e := &ErrNotSameLength{csv: len(prices), img: len(images)}
		logger.Warn(e.Error())
	}

	// check matches between images and prices
	for k, v := range images {
		if _, ok := prices[k]; !ok {
			e := &ErrNoSameFile{filename: k}
			logger.Warn(fmt.Sprintf("нет соответствующего csv-файла %s%s", e.filename, v))
		}
	}

	// check matches between prices and images
	for k := range prices {
		if _, ok := images[k]; !ok {
			e := &ErrNoSameFile{filename: k, ext: "jpg/png"}
			logger.Error(e.Error())
			continue
		}

		// save matched pairs
		res[k] = images[k]
		logger.Debug(fmt.Sprintf("Пара: %s", k))
	}

	return res, nil
}

func (p *processor) CheckDate(files map[string]string, date string) {
	for k := range files {
		fileName := filepath.Base(k)
		if !strings.HasPrefix(fileName, date) {
			err := &ErrNoNeededDate{filename: fileName, date: date}
			logger.Warn(err.Error())
			delete(files, k)
		}
	}
}
