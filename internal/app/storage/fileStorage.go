package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/koteyye/shortener/internal/app/models"
	"log"
	"os"
)

// FileStorage структура файлового хранилища.
type FileStorage struct {
	fileWriter *FileWriter
	fileReader *FileReader
}

// NewFileStorage возвращает новый экземпляр файлового хранилища.
func NewFileStorage(filepath string) *FileStorage {
	return &FileStorage{
		fileWriter: &FileWriter{filePath: filepath},
		fileReader: &FileReader{filePath: filepath},
	}
}

// DeleteURLByUser удаление URL текущего пользователя (не поддерживается).
func (f *FileStorage) DeleteURLByUser(_ context.Context, _ chan string) error {
	return models.ErrMockNotSupported
}

// GetURLByUser получение URL текущего пользователя (не поддерживается).
func (f *FileStorage) GetURLByUser(_ context.Context, _ string) ([]*models.AllURLs, error) {
	return nil, models.ErrFileNotSupported
}

// GetShortURL получение сокращенного URL по оригинальному (не поддерживается).
func (f *FileStorage) GetShortURL(_ context.Context, _ string) (string, error) {
	return "", models.ErrFileNotSupported
}

// Ping проверка подключения к БД (не поддерживается).
func (f *FileStorage) Ping(_ context.Context) error {
	return models.ErrFileNotSupported
}

// AddURL добавление URL в файловое хранилище.
func (f *FileStorage) AddURL(_ context.Context, s string, k string, _ string) error {
	var id int

	reader, err := f.fileReader.NewReader()
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

	writer, err := f.fileWriter.NewWriter()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer writer.Close()

	err = writer.WriteShortURL(models.AllURLs{
		ID:          id,
		ShortURL:    s,
		OriginalURL: k,
	})
	if err != nil {
		return fmt.Errorf("err write shortURL: %w", err)
	}
	return nil
}

// GetURL получение URL из файлового хранилища.
func (f *FileStorage) GetURL(_ context.Context, k string) (*models.URL, error) {
	reader, err := f.fileReader.NewReader()
	if err != nil {
		return &models.URL{}, fmt.Errorf("err reader: %w", err)
	}
	defer reader.Close()

	readFile, err := reader.FindOriginalURL(k)
	if err != nil {
		return &models.URL{}, fmt.Errorf("err read file: %w", err)
	}
	return &models.URL{OriginalURL: readFile.OriginalURL}, nil
}

// FileWriter структура файлового писателя.
type FileWriter struct {
	filePath string
	file     *os.File
	encoder  *json.Encoder
}

// FileReader структура файлового читателя.
type FileReader struct {
	filePath string
	file     *os.File
}

// NewWriter возвращает новый экземпляр файлового писателя.
func (w *FileWriter) NewWriter() (*FileWriter, error) {
	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileWriter{file: file, encoder: json.NewEncoder(file)}, err
}

// WriteShortURL записать сокращенный URL в файл.
func (w *FileWriter) WriteShortURL(value models.AllURLs) error {
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

// NewReader возвращает новый экземпляр файлового читателя.
func (r *FileReader) NewReader() (*FileReader, error) {
	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &FileReader{file: file}, nil
}

// ReadShortURL читает сокращенный URL в файле.
func (r *FileReader) ReadShortURL() (*models.AllURLs, error) {
	var fileString models.AllURLs

	scanner := bufio.NewScanner(r.file)

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &fileString)
		if err != nil {
			return nil, err
		}
	}

	return &fileString, nil
}

// Close закрывает файл
func (r *FileReader) Close() error {
	return r.file.Close()
}

// FindOriginalURL найти оригинальный URL в файле по сокращенному
func (r *FileReader) FindOriginalURL(shortURL string) (*models.AllURLs, error) {
	var fileString models.AllURLs

	scanner := bufio.NewScanner(r.file)

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &fileString)
		if err != nil {
			return nil, err
		}
		if fileString.ShortURL == shortURL {
			break
		}
	}
	if fileString.ShortURL != shortURL {
		return nil, models.ErrNotFound
	}
	return &fileString, nil
}
