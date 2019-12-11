package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *uploadSetting) addToIPFS() {
	var cipherShareNum int = 0
	// 対象シェア
	for mIndex := 0; mIndex < len(s.managers); mIndex++ {
		for sIndex := 0; sIndex < s.managers[mIndex].manageShareNum; sIndex++ {
			index := strconv.Itoa(sIndex + 1)
			name := strings.Replace(s.managers[mIndex].fileName, ".", "_", -1)
			apiRequest(s.writeFilePath + "/" + name + "_share" + index)
			cipherShareNum++
		}
	}
	// 非対象シェア
	for i := cipherShareNum; i < s.shareNum; i++ {
		index := strconv.Itoa(i - cipherShareNum + 1)
		apiRequest(s.writeFilePath + "/un_managed_share" + index)
	}
}

func apiRequest(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", "http://localhost:5001/api/v0/add", body)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
}
