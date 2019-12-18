package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/SSSaaS/sssa-golang"
)

// 一時ディレクトリ作成
func (s *commonSetting) makeTempDir() {
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

	file, err := os.Open(s.comSet.readFilePath)
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
	var cipherShareNum int = 0

	// 対象シェア
	for mIndex := 0; mIndex < len(s.managers); mIndex++ {

		// 各シェアを共通鍵で暗号化
		comKey := genComKey()
		for sIndex := 0; sIndex < s.managers[mIndex].manageShareNum; sIndex++ {
			fmt.Printf("%d / %d\r", cipherShareNum+1, s.cipherShareNum)
			block, err := aes.NewCipher(comKey)
			if err != nil {
				log.Fatal(err)
			}
			gcm, err := cipher.NewGCM(block)
			if err != nil {
				log.Fatal(err)
			}
			nonce := make([]byte, gcm.NonceSize())
			_, err = rand.Read(nonce)
			if err != nil {
				log.Fatal(err)
			}
			content := []byte(s.created[cipherShareNum])
			cipherContent := gcm.Seal(nil, nonce, content, nil)
			cipherContent = append(nonce, cipherContent...)
			index := strconv.Itoa(sIndex + 1)
			name := strings.Replace(s.managers[mIndex].keyfileName, ".", "_", -1)
			err = ioutil.WriteFile(s.comSet.tempDirPath+"/"+name+"_share"+index, cipherContent, 0755)
			if err != nil {
				log.Fatal(err)
			}
			cipherShareNum++
		}

		// 共通鍵を公開鍵で暗号化
		rng := rand.Reader
		encryptedComKey, err := rsa.EncryptOAEP(sha256.New(), rng, s.managers[mIndex].publicKey, comKey, []byte(""))
		if err != nil {
			log.Fatal(err)
		}
		s.managers[mIndex].config.EncryptedComKey = encryptedComKey
	}

	// 非対象シェア
	for i := s.cipherShareNum; i < s.shareNum; i++ {
		index := strconv.Itoa(i - s.cipherShareNum + 1)
		err := ioutil.WriteFile(s.comSet.tempDirPath+"/un_managed_share"+index, []byte(s.created[i]), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	s.created = nil
}

// 共通鍵の生成
func genComKey() []byte {
	keyList := "abcdefghijklmnopqrstuvwxyzABCDEFHFGHIJKLMNOPQRSTUVWXYZ1234567890"
	size := 32 //256bit

	var key = make([]byte, 0, size)
	for i := 1; i <= size; i++ {
		res, _ := rand.Int(rand.Reader, big.NewInt(64))
		keyGen := keyList[res.Int64()]
		key = append(key, keyGen)
	}
	return key
}
