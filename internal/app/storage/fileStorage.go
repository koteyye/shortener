package storage

import (
	"bufio"
	"encoding/json"
	"github.com/koteyye/shortener/internal/app/models"
	"os"
	"path/filepath"
	"strings"
)

type FileStorage struct {
	FileWriter
	FileReader
}

type FileWriter struct {
	filePath string
	file     *os.File
	encoder  *json.Encoder
}

type FileReader struct {
	filePath string
	file     *os.File
	scanner  *bufio.Scanner
}

func (w FileWriter) NewWriter() (*FileWriter, error) {
	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileWriter{file: file, encoder: json.NewEncoder(file)}, err
}

func (w *FileWriter) WriteShortURL(value models.FileString) error {
	data, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	_, err = w.file.Write(data)
	return err
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}

func (w *FileWriter) Mkdir() error {
	path := strings.TrimLeft(filepath.Dir(w.filePath), "/")
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (r *FileReader) NewReader() (*FileReader, error) {
	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &FileReader{file: file, scanner: bufio.NewScanner(file)}, nil
}

func (r *FileReader) ReadShortURL() (*models.FileString, error) {
	var fileString models.FileString

	for r.scanner.Scan() {
		err := json.Unmarshal(r.scanner.Bytes(), &fileString)
		if err != nil {
			return nil, err
		}
	}

	return &fileString, nil
}

func (r *FileReader) Close() error {
	return r.file.Close()
}

func (r *FileReader) FindOriginalURL(shortURL string) (*models.FileString, error) {
	var fileString models.FileString

	for r.scanner.Scan() {
		err := json.Unmarshal(r.scanner.Bytes(), &fileString)
		if err != nil {
			return nil, err
		}
		if fileString.ShortURL == shortURL {
			break
		}
	}
	if fileString.ShortURL != shortURL {
		return nil, ErrNotFound
	}
	return &fileString, nil
}
