package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/SSSaaS/sssa-golang"
)

// 一時ディレクトリ作成
func (s *uploadSetting) makeTempDir() {
	s.tempDirPath = "temp"
	if _, err := os.Stat(s.tempDirPath); os.IsNotExist(err) {
		err = os.Mkdir(s.tempDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// 秘密分散法
func (s *uploadSetting) sssaCreate() {

	file, err := os.Open(s.readFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	s.created, err = sssa.Create(s.minNum, s.shareNum, string(raw))
	if err != nil {
		log.Fatal(err)
	}
}

// 暗号化
func (s *uploadSetting) encrypt() {

	// 対象シェア
	for mIndex := 0; mIndex < len(s.managers); mIndex++ {
		for sIndex := 0; sIndex < s.managers[mIndex].manageShareNum; sIndex++ {
			content := []byte(s.created[s.cipherShareNum])
			rng := rand.Reader
			cipherContent, err := rsa.EncryptOAEP(sha256.New(), rng, s.managers[mIndex].publicKey, content, []byte(""))
			if err != nil {
				log.Fatal(err)
			}
			index := strconv.Itoa(sIndex + 1)
			name := strings.Replace(s.managers[mIndex].fileName, ".", "_", -1)
			err = ioutil.WriteFile(s.tempDirPath+"/"+name+"_share"+index, cipherContent, 0755)
			if err != nil {
				log.Fatal(err)
			}
			s.cipherShareNum++
		}
	}
	// 非対象シェア
	for i := s.cipherShareNum; i < s.shareNum; i++ {
		index := strconv.Itoa(i - s.cipherShareNum + 1)
		err := ioutil.WriteFile(s.tempDirPath+"/un_managed_share"+index, []byte(s.created[i]), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
