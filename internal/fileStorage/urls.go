package fileStorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type UrlsFilesRW struct {
	fileToWrite *os.File
	fileToRead  *os.File
	writer      *bufio.Writer
	scanner     *bufio.Scanner
}

func (u UrlsFilesRW) WriteURLInFile(fullURL string, id string) error {
	currentURL := oneURL{
		Id:      id,
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

func (u UrlsFilesRW) ReadAllURLFromFile(id string) (string, error) {
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

func NewUrls() (*UrlsFilesRW, error) {
	//TODO получение переменных окружения
	var filename string
	if filename = os.Getenv("FILE_STORAGE_PATH"); filename == "" {
		filename = "internal/fileStorage/URLStorage.json"
	}
	fileToWrite, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
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
	Id      string
	FullURL string
}
