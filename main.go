package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/SSSaaS/sssa-golang"
)

type keyManager struct {
	fileName       string
	publicKey      *rsa.PublicKey
	manageShareNum int
}

type uploadSetting struct {
	managers []keyManager
	shareNum int
	minNum   int
	filePath string
	created  []string
}

func main() {
	flag.Parse()
	commands := flag.Arg(0)

	if commands == "upload" {
		upload()
	} else if commands == "download" {
		download()
	} else {
		fmt.Printf("エラー : 適切なコマンドを入力してください\n")
		fmt.Printf("例 : \n")
		fmt.Printf("  sss-ipfs upload\n")
		fmt.Printf("  sss-ipfs download\n")
	}
}

func upload() {

	var setting uploadSetting

	// 初期設定
	setting.askPubKeys()
	setting.askShareNum()
	setting.askMinNum()
	setting.askShareManagers()
	setting.askFilePath()

	file, err := os.Open(setting.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// 秘密分散
	created, err := sssa.Create(setting.minNum, setting.shareNum, string(raw))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		err = os.Mkdir("temp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	var cipherShareNum int = 0

	// 公開鍵暗号
	// データが小さいファイルのみ
	for mIndex := 0; mIndex < len(setting.managers); mIndex++ {
		for sIndex := 0; sIndex < setting.managers[mIndex].manageShareNum; sIndex++ {
			content := []byte(created[cipherShareNum])
			rng := rand.Reader
			cipherContent, err := rsa.EncryptOAEP(sha256.New(), rng, setting.managers[mIndex].publicKey, content, []byte(""))
			if err != nil {
				log.Fatal(err)
			}
			index := strconv.Itoa(sIndex + 1)
			name := strings.Replace(setting.managers[mIndex].fileName, ".", "_", -1)
			err = ioutil.WriteFile("temp/"+name+"_share"+index, cipherContent, 0755)
			if err != nil {
				log.Fatal(err)
			}
			cipherShareNum++
		}
	}
	for i := cipherShareNum; i < setting.shareNum; i++ {
		index := strconv.Itoa(i - cipherShareNum + 1)
		err = ioutil.WriteFile("temp/un_managed_share"+index, []byte(created[i]), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func download() {

}
