package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type UrlsFilesRW struct {
	fileToWrite *os.File
	fileToRead  *os.File
	writer      *bufio.Writer
	scanner     *bufio.Scanner
}

func (u UrlsFilesRW) WriteURL(fullURL string, id string, userID uint64) error {
	currentURL := oneURL{
		ID:      id,
		FullURL: fullURL,
		UserID:  userID,
	}
	data, err := json.Marshal(&currentURL)
	if err != nil {
		return err
	}
	log.Println(string(data))
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

func (u UrlsFilesRW) ReadNextURL() (string, uint64, error) {
	if !u.scanner.Scan() {
		return "", 0, u.scanner.Err()
	}
	// читаем данные из scanner
	data := u.scanner.Bytes()

	event := oneURL{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return "", 0, err
	}

	return event.FullURL, event.UserID, nil
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

func NewUrls(unitTestFlag bool, fileStoragePath string) (*UrlsFilesRW, error) {
	//Получение адреса файла из переменных окружения
	var filename string

	//Получим абсолютный путь к модулю
	_, callerName, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get caller info")
	}
	modulePath := filepath.Dir(callerName)
	filename = fmt.Sprintf("%s%s", modulePath, fileStoragePath)

	// Для инициализации пути в unit-тестах и использования переменных окружения
	if unitTestFlag || fileStoragePath != "/URLStorage.json" {
		filename = fileStoragePath
	}

	// Открыть файл для записи
	fileToWrite, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	// Открыть файл для чтения
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
	UserID  uint64
}
