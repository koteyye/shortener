package service

import (
	"errors"
	"github.com/koteyye/shortener/internal/app/storage"
	"hash/crc64"
	"strings"
)

const solt = "jifo4j0f9jf09djsa0f92039r32fj09sdjfg09ewjg09ewgj_"

const soltLen = 20

type ShortenerService struct {
	storage storage.URLStorage
}

func NewShortenerService(storage storage.URLStorage) *ShortenerService {
	return &ShortenerService{storage: storage}
}

func (s ShortenerService) LongURL(shortURL string) (string, error) {
	res, err := s.storage.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s ShortenerService) ShortURL(url string) (string, error) {
	var res string
	if res = hashString(url); res == "" {
		return "", errors.New("не удалось сократить")
	}
	ok, _ := s.storage.AddURL(res, url)
	if !ok {
		return "", errors.New("говно")
	}
	urlRes := "http://localhost:8080/" + res
	return urlRes, nil
}

func hashString(s string) string {
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.Checksum([]byte(s), table)

	charArray := make([]uint8, 0, 10)

	for mod := uint64(0); hash != 0 && len(charArray) < 10; {
		mod = hash % soltLen
		hash /= soltLen
		charArray = append(charArray, solt[mod])
	}

	return strings.Repeat(string(solt[0]), 10-len(charArray)) + string(charArray)
}
