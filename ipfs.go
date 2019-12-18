package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
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

type responseJSON struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

// 出力用ディレクトリ作成
func (s *commonSetting) makeWriteDir() {
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
		fmt.Printf("%d / %d\r", mIndex+1, s.shareNum)
		for sIndex := 0; sIndex < s.managers[mIndex].manageShareNum; sIndex++ {
			index := strconv.Itoa(sIndex + 1)
			name := strings.Replace(s.managers[mIndex].keyfileName, ".", "_", -1)
			hash := apiAddRequest(s.comSet.tempDirPath + "/" + name + "_share" + index)
			s.managers[mIndex].config.ManagedShares = append(s.managers[mIndex].config.ManagedShares, hash)
		}
	}
	// 非対象シェア
	for i := s.cipherShareNum; i < s.shareNum; i++ {
		fmt.Printf("\r%d / %d", i+1, s.shareNum)
		index := strconv.Itoa(i - s.cipherShareNum + 1)
		hash := apiAddRequest(s.comSet.tempDirPath + "/un_managed_share" + index)
		for mIndex := 0; mIndex < len(s.managers); mIndex++ {
			s.managers[mIndex].config.UnmanagedShares = append(s.managers[mIndex].config.UnmanagedShares, hash)
		}
	}
}

func apiAddRequest(path string) string {
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
	resJSON := new(responseJSON)
	if err := json.Unmarshal(content, resJSON); err != nil {
		log.Fatal(err)
	}

	return resJSON.Hash
}

// 共有用コンフィグ作成
func (s *uploadSetting) writeConfig() {
	for mIndex := 0; mIndex < len(s.managers); mIndex++ {
		// 各管理者用のディレクトリ作成
		dirName := strings.Replace(s.managers[mIndex].keyfileName, ".", "_", -1)
		dirPath := s.comSet.writeDirPath + "/" + dirName
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
		pubKeyBytes, err := x509.MarshalPKIXPublicKey(s.managers[mIndex].publicKey)
		if err != nil {
			log.Fatal(err)
		}
		var block = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubKeyBytes,
		}
		pem.Encode(file, block)

		// JSON形式コンフィグファイルの追加
		for sIndex := 0; sIndex < s.managers[mIndex].manageShareNum; sIndex++ {
			jsonBytes, err := json.Marshal(s.managers[mIndex].config)
			if err != nil {
				log.Fatal(err)
			}
			// インデントの整形
			out := &bytes.Buffer{}
			json.Indent(out, jsonBytes, "", "    ")
			writeJSON, err := ioutil.ReadAll(out)
			if err != nil {
				log.Fatal(err)
			}
			err = ioutil.WriteFile(dirPath+"/config.json", writeJSON, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// コンフィグ読み取り
func (s *downloadSetting) readConfig() {

	file, err := os.Open(s.comSet.readFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	configJSON := new(configuration)
	if err := json.Unmarshal(raw, configJSON); err != nil {
		log.Fatal(err)
	}
	s.manager.config = *configJSON
}

// IPFSからダウンロード
func (s *downloadSetting) getFromIPFS() {
	for i := 0; i < len(s.manager.config.ManagedShares); i++ {
		fmt.Printf("\r%d / %d", i+1, len(s.manager.config.ManagedShares))
		apiCatRequest(s.manager.config.ManagedShares[i])
	}
	fmt.Printf("\n")
	for i := 0; i < len(s.manager.config.UnmanagedShares); i++ {
		fmt.Printf("\r%d / %d", i+1, len(s.manager.config.UnmanagedShares))
		apiCatRequest(s.manager.config.UnmanagedShares[i])
	}
}

func apiCatRequest(cid string) {

	request, err := http.NewRequest("GET", "http://localhost:5001/api/v0/cat", nil)
	if err != nil {
		log.Fatal(err)
	}
	params := request.URL.Query()
	params.Add("arg", cid)
	request.URL.RawQuery = params.Encode()

	request.Header.Add("Content-Type", "text/plain")
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
	err = ioutil.WriteFile("./temp/"+cid, content, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
