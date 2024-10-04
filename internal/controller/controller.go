package controller

import (
	"fmt"
	"sync"

	"go-flattern-prices/internal/configs"
	"go-flattern-prices/internal/logger"
	"go-flattern-prices/internal/processors"
	"go-flattern-prices/internal/store"
)

type controller struct {
	config *configs.Config
	store  store.Store
	files  map[string]string // {path/to/file:image_extension, } like {~/data/230101_smartbox: jpeg}
}

type Controller interface {
	DoJob() error
	loadFiles() error
}

func New(cfg *configs.Config, s store.Store) Controller {
	return &controller{
		config: cfg,
		store:  s,
	}
}

func (c *controller) DoJob() error {
	if err := c.loadFiles(); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Нашлось %d пар файлов", len(c.files)))

	wg := sync.WaitGroup{}

	for k, v := range c.files {
		wg.Add(1)
		go c.job(k, v, &wg)
	}

	wg.Wait()
	return nil
}

func (c *controller) loadFiles() error {
	p := processors.New(c.config.InputDir)

	files, err := p.GetFilesPairs()
	if err != nil {
		return err
	}

	p.CheckDate(files, c.config.Date)

	c.files = files
	return nil
}

func (c *controller) job(k, v string, wg *sync.WaitGroup) {
	defer wg.Done()

	csvFile := fmt.Sprintf("%s.csv", k)

	tags, err := processors.ReadTags(csvFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка чтения файла %s", csvFile))
		return
	}
	logger.Debug(fmt.Sprintf("Найдено тегов '%d' для файла %s", len(tags), csvFile))

	for _, tag := range tags {
		price, e := c.store.GetPriceByID(tag.Code, c.config.Date)
		if e != nil {
			logger.Error(fmt.Sprintf("Ошибка для кода %d дата %s err: %v", tag.Code, csvFile, e))
			continue
		}

		tag.Price = &price
		logger.Debug(fmt.Sprintf("Присвоена цена %d для кода %d дата %s tag: %v", *tag.Price, tag.Code, csvFile, tag))
	}

	imageFile := fmt.Sprintf("%s.%s", k, v)
	err = processors.DrawTagsOverImage(k, v, tags, c.config.OutputDir, c.config.InputDir)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка при отрисовке цены %s", imageFile), "error", err)
	}
}
