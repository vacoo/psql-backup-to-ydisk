package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
)

// SendError Отправка ошибки на почту
func SendError(logError error) {
	from := mail.Address{Name: os.Getenv("PROJECT"), Address: os.Getenv("MAIL_SMPT_USER")}
	to := mail.Address{Name: "mail", Address: os.Getenv("MAIL_TO")}
	subj := "Ошибка резервного копирования базы данных "
	body := "ERR_MSG: " + logError.Error()

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := fmt.Sprintf("%s:%s", os.Getenv("MAIL_SMPT_HOST"), os.Getenv("MAIL_SMPT_PORT"))
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", os.Getenv("MAIL_SMPT_USER"), os.Getenv("MAIL_SMPT_PASS"), host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = c.Auth(auth); err != nil {
		fmt.Println(err)
		return
	}
	if err = c.Mail(from.Address); err != nil {
		fmt.Println(err)
		return
	}
	if err = c.Rcpt(to.Address); err != nil {
		fmt.Println(err)
		return
	}
	w, err := c.Data()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = w.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Quit()
}

// RequestParams ...
type RequestParams struct {
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
	Query  string      `json:"query"`
	URL    string      `json:"path"`
}

// Request ...
func Request(method, url string) ([]byte, int, error) {
	var bodyByte []byte

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyByte))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", os.Getenv("YANDEX_DISK_ACCESS_TOKEN")))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody, resp.StatusCode, err
}

// UploadFile Загрузка файла на сервер
func UploadFile(url string, filename string) ([]byte, int, error) {
	file, err := os.Open(filename)

	if err != nil {
		return []byte{}, 0, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, 0, err
	}
	body.Write(b)

	req, err := http.NewRequest("PUT", url, body)

	if err != nil {
		return []byte{}, 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody, resp.StatusCode, err
}
