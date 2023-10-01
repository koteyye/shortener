package storage

import (
	"errors"
	"github.com/koteyye/shortener/internal/app/models"
	"log"
	"sync"
)

var ErrNotFound = errors.New("не найдено такого значения")

type URLStorage interface {
	AddURL(string, string)
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

func (u *URLMap) AddURL(k, s string) {
	b := u.fileStorage.FileWriter.filePath
	if b != "" {
		var id int
		err := u.fileStorage.FileWriter.Mkdir()
		if err != nil {
			log.Fatal(err)
			return
		}

		reader, err := u.fileStorage.FileReader.NewReader()
		if err != nil {
			log.Fatal(err)
			return
		}
		defer reader.Close()

		readFile, err := reader.ReadShortURL()
		if readFile == nil {
			id = 1
		} else {
			id = readFile.Id + 1
		}

		writer, err := u.fileStorage.FileWriter.NewWriter()
		if err != nil {
			log.Fatal(err)
			return
		}
		defer writer.Close()

		err = writer.WriteShortURL(models.FileString{
			Id:          id,
			ShortURL:    k,
			OriginalURL: s,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		u.storage.Store(k, s)
	}
}

func (u *URLMap) GetURL(k string) (string, error) {
	b := u.fileStorage.FileReader.filePath
	if b != "" {
		err := u.fileStorage.FileWriter.Mkdir()
		if err != nil {
			return "", err
		}

		reader, err := u.fileStorage.FileReader.NewReader()
		if err != nil {
			return "", err
		}
		defer reader.Close()

		readFile, err := reader.FindOriginalURL(k)
		if err != nil {
			return "", err
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
