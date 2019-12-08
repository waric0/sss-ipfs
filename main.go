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

	var managers []keyManager

	// 初期設定
	managers = askPubKeys(managers)
	shareNum := askShareNum(managers)
	minNum := askMinNum(shareNum)
	managers = askShareManagers(managers, shareNum, minNum)
	filePath := askFilePath()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// 秘密分散
	created, err := sssa.Create(minNum, shareNum, string(raw))
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
	for mIndex := 0; mIndex < len(managers); mIndex++ {
		for sIndex := 0; sIndex < managers[mIndex].manageShareNum; sIndex++ {
			content := []byte(created[cipherShareNum])
			rng := rand.Reader
			cipherContent, err := rsa.EncryptOAEP(sha256.New(), rng, managers[mIndex].publicKey, content, []byte(""))
			if err != nil {
				log.Fatal(err)
			}
			index := strconv.Itoa(sIndex + 1)
			name := strings.Replace(managers[mIndex].fileName, ".", "_", -1)
			err = ioutil.WriteFile("temp/"+name+"_share"+index, cipherContent, 0755)
			if err != nil {
				log.Fatal(err)
			}
			cipherShareNum++
		}
	}
	for i := cipherShareNum; i < shareNum; i++ {
		index := strconv.Itoa(i - cipherShareNum + 1)
		err = ioutil.WriteFile("temp/un_managed_share"+index, []byte(created[i]), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func download() {

}
