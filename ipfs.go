package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
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

// 出力用ディレクトリ作成
func (s *uploadSetting) makeWriteDir() {
	s.writeDirPath = "outputs"
	if _, err := os.Stat(s.writeDirPath); os.IsNotExist(err) {
		err = os.Mkdir(s.writeDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// IPFSへアップロード
func (s *uploadSetting) addToIPFS() {

	// 対象シェア
	for mIndex := 0; mIndex < len(s.managers); mIndex++ {
		for sIndex := 0; sIndex < s.managers[mIndex].manageShareNum; sIndex++ {
			index := strconv.Itoa(sIndex + 1)
			name := strings.Replace(s.managers[mIndex].fileName, ".", "_", -1)
			apiRequest(s.tempDirPath + "/" + name + "_share" + index)
		}
	}
	// 非対象シェア
	for i := s.cipherShareNum; i < s.shareNum; i++ {
		index := strconv.Itoa(i - s.cipherShareNum + 1)
		apiRequest(s.tempDirPath + "/un_managed_share" + index)
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

// 共有用コンフィグ作成
func (s *uploadSetting) writeConfig() {
	for mIndex := 0; mIndex < len(s.managers); mIndex++ {
		// 各管理者用のディレクトリ作成
		dirName := strings.Replace(s.managers[mIndex].fileName, ".", "_", -1)
		dirPath := s.writeDirPath + "/" + dirName
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err = os.Mkdir(dirPath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
		// PEM形式公開鍵の追加
		file, err := os.Create(dirPath + "/pub-key.pem")
		if err != nil {
			log.Fatal(err)
		}
		bytes, err := x509.MarshalPKIXPublicKey(s.managers[mIndex].publicKey)
		if err != nil {
			log.Fatal(err)
		}
		var block = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: bytes,
		}
		pem.Encode(file, block)

	}

}
