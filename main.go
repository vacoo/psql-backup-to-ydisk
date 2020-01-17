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
	dumpName := os.Getenv("PSQL_DB") + "-" + time.Now().Format("2006-01-02_15-04") + ".gz"

	// Делаем бэкап в папку backups
	_, err := exec.Command("sh", "./util_dump.sh").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
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

	os.Exit(0)
}
