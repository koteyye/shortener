package storage

import (
	"errors"
	"fmt"
	"github.com/koteyye/shortener/internal/app/models"
	"log"
	"sync"
)

var ErrNotFound = errors.New("не найдено такого значения")

type URLStorage interface {
	AddURL(string, string) error
	GetURL(string) (string, error)
}

type URLHandler struct {
	URLStorage
}

type URLMap struct {
	storage     sync.Map
	fileStorage FileStorage
}

func NewURLHandle(filePath string) *URLHandler {
	if filePath != "" {
		return &URLHandler{
			URLStorage: &URLMap{
				fileStorage: FileStorage{
					FileWriter: FileWriter{
						filePath: filePath,
					},
					FileReader: FileReader{
						filePath: filePath,
					},
				},
			},
		}
	}

	return &URLHandler{
		URLStorage: &URLMap{
			storage: sync.Map{},
		},
	}
}

func (u *URLMap) AddURL(k, s string) error {
	b := u.fileStorage.FileWriter.filePath
	if b != "" {
		var id int
		//err := u.fileStorage.FileWriter.Mkdir()
		//if err != nil {
		//	return err
		//}

		reader, err := u.fileStorage.FileReader.NewReader()
		if err != nil {
			return fmt.Errorf("err reader: %w", err)
		}
		defer reader.Close()

		readFile, err := reader.ReadShortURL()
		if err != nil {
			return fmt.Errorf("err read file: %w", err)
		}
		if readFile == nil {
			id = 1
		} else {
			id = readFile.ID + 1
		}

		writer, err := u.fileStorage.FileWriter.NewWriter()
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer writer.Close()

		err = writer.WriteShortURL(models.FileString{
			ID:          id,
			ShortURL:    k,
			OriginalURL: s,
		})
		if err != nil {
			return fmt.Errorf("err write shortURL: %w", err)
		}
		return nil
	} else {
		u.storage.Store(k, s)
		return nil
	}
}

func (u *URLMap) GetURL(k string) (string, error) {
	b := u.fileStorage.FileReader.filePath
	if b != "" {
		//err := u.fileStorage.FileWriter.Mkdir()
		//if err != nil {
		//	return "", err
		//}

		reader, err := u.fileStorage.FileReader.NewReader()
		if err != nil {
			return "", fmt.Errorf("err reader: %w", err)
		}
		defer reader.Close()

		readFile, err := reader.FindOriginalURL(k)
		if err != nil {
			return "", fmt.Errorf("err read file: %w", err)
		}
		return readFile.OriginalURL, nil

	} else {
		url, ok := u.storage.Load(k)
		if !ok {
			return "", ErrNotFound
		}
		return url.(string), nil
	}
}
