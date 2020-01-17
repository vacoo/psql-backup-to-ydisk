package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	mode := os.Args[1]

	if mode == "" {
		fmt.Println("Пожалуйста выберите режим: backup | restore")
		os.Exit(1)
	}

	if mode == "backup" {
		Backup()
	}
	if mode == "restore" {
		Restore()
	}

	os.Exit(0)
}

// Restore Восстановление
func Restore() {
	if len(os.Args) == 2 {
		fmt.Println("Введите отностительный путь к дампу на яндекс диска: 2020-01/database-2020-01-17_10-08.sql.gz")
		os.Exit(1)
	}
	if len(os.Args) == 3 {
		fmt.Println("Укажите базу данных для восстановления: database1")
		os.Exit(1)
	}

	// Берем указанный путь к дампу
	filePath := os.Args[2]
	targetDatabase := os.Args[3]

	// Ссылка на скачивание дампа
	folderPath := fmt.Sprintf("disk:/Приложения/%s/%s", os.Getenv("YANDEX_DISK_APP_FOLDER"), filePath)
	url := fmt.Sprintf("https://cloud-api.yandex.net:443/v1/disk/resources/download?path=%s", folderPath)
	resp, status, err := Request("GET", url)
	if err != nil {
		SendError(err)
		os.Exit(0)
	}
	if status != http.StatusOK {
		SendError(errors.New(string(resp)))
		os.Exit(0)
	}

	var link map[string]interface{}
	err = json.Unmarshal(resp, &link)
	if err != nil {
		SendError(err)
		os.Exit(0)
	}

	// Извлекаем путь скачивания файла
	href := link["href"].(string)
	fileName := filepath.Base(filePath)
	p := "backups/" + fileName

	// Скачиваем дамп
	if err := DownloadFile(p, href); err != nil {
		SendError(err)
		os.Exit(0)
	}

	fmt.Println("Дамп успешно загружен в " + p)

	fmt.Println("Восстановление в базу данных " + targetDatabase + "...")

	// Восстановление
	_, err = exec.Command("sh", "./util_restore.sh", p, targetDatabase).CombinedOutput()
	if err != nil {
		SendError(err)
		os.Exit(0)
	}

	fmt.Println("Успешное восстановление")

	os.Exit(0)
}

// Backup Бэкап
func Backup() {
	dumpName := os.Getenv("PSQL_DB") + "-" + time.Now().Format("2006-01-02_15-04") + ".sql.gz"

	// Делаем бэкап в папку backups
	_, err := exec.Command("sh", "./util_dump.sh").CombinedOutput()
	if err != nil {
		SendError(err)
		os.Exit(0)
	}

	file := "backups/" + dumpName
	folder := time.Now().Format("2006-01")
	fileName := filepath.Base(file)
	folderPath := fmt.Sprintf("disk:/Приложения/%s/%s", os.Getenv("YANDEX_DISK_APP_FOLDER"), folder)

	// Создаем папку в хранилище
	url := fmt.Sprintf("https://cloud-api.yandex.net:443/v1/disk/resources?path=%s", folderPath)
	resp, status, err := Request("PUT", url)
	if err != nil {
		SendError(err)
		os.Exit(0)
	}
	if !(status == http.StatusCreated || status == http.StatusConflict) {
		SendError(errors.New(string(resp)))
		os.Exit(0)
	}

	fmt.Println("OK папка создана")

	// Получаем ссылку для загрузки файла
	url = fmt.Sprintf("https://cloud-api.yandex.net:443/v1/disk/resources/upload?path=%s/%s", folderPath, fileName)
	resp, status, err = Request("GET", url)
	if err != nil {
		SendError(err)
		os.Exit(0)
	}
	if status != http.StatusOK {
		SendError(errors.New(string(resp)))
		os.Exit(0)
	}

	var link map[string]interface{}
	err = json.Unmarshal(resp, &link)
	if err != nil {
		SendError(err)
		os.Exit(0)
	}

	href := link["href"].(string)
	fmt.Println("OK Получена ссылка для загрузки:", href)

	fmt.Println("Загрузка файла на yandex disk...")
	// Загружаем файл
	resp, status, err = UploadFile(href, file)
	if err != nil {
		SendError(err)
		os.Exit(0)
	}
	if status != http.StatusCreated {
		SendError(errors.New(string(resp)))
		os.Exit(0)
	}

	fmt.Println("Файл загружен ", folderPath+"/"+fileName)

	// Удаляем дамп из локального диска
	err = os.Remove(file)
	if err != nil {
		SendError(err)
		os.Exit(0)
	}

	os.Exit(0)
}
