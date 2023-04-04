package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"yaGoShortURL/internal/static"
)

type UrlsFilesRW struct {
	fileToWrite *os.File
	fileToRead  *os.File
	writer      *bufio.Writer
	scanner     *bufio.Scanner
}

func (u UrlsFilesRW) WriteURLInFile(fullURL string, id string) error {
	currentURL := oneURL{
		ID:      id,
		FullURL: fullURL,
	}
	fmt.Println(currentURL)
	data, err := json.Marshal(&currentURL)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	// записываем событие в буфер
	if _, err := u.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := u.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return u.writer.Flush()
}

func (u UrlsFilesRW) ReadNextURLFromFile() (string, error) {
	if !u.scanner.Scan() {
		return "", u.scanner.Err()
	}
	// читаем данные из scanner
	data := u.scanner.Bytes()

	event := oneURL{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return "", err
	}

	return event.FullURL, nil
}

func (u UrlsFilesRW) CloseFile() error {
	err := u.fileToWrite.Close()
	if err != nil {
		return err
	}
	err = u.fileToRead.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewUrls(cfg static.ConfigInit) (*UrlsFilesRW, error) {
	//Получение адреса файла из переменных окружения
	var filename string

	//Получим абсолютный путь к модулю
	_, callerName, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get caller info")
	}

	modulePath := filepath.Dir(callerName)
	filename = fmt.Sprintf("%s%v", modulePath, cfg.FileStoragePath)

	// Для инициализации пути в unit-тестах
	if cfg.UnitTestFlag {
		filename = cfg.FileStoragePath
	}

	// Открыть файл для записи
	fileToWrite, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	// Открыть файлл для чтения
	fileToRead, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &UrlsFilesRW{
		//Для записи
		fileToWrite: fileToWrite,
		writer:      bufio.NewWriter(fileToWrite),
		//Для чтения
		fileToRead: fileToRead,
		scanner:    bufio.NewScanner(fileToRead),
	}, nil
}

// Структура для записи в json-формате
type oneURL struct {
	ID      string
	FullURL string
}