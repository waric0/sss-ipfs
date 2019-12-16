package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
)

type keyManager struct {
	fileName       string
	publicKey      *rsa.PublicKey
	manageShareNum int
}

type uploadSetting struct {
	managers       []keyManager
	shareNum       int
	minNum         int
	cipherShareNum int
	readFilePath   string
	tempDirPath    string
	writeDirPath   string
	created        []string
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

	var s uploadSetting

	// 初期設定処理
	s.askPubKeys()
	s.askShareNum()
	s.askMinNum()
	s.askShareManagers()
	s.askFilePath()

	// 加工処理
	s.makeTempDir()
	s.sssaCreate()
	s.encrypt()

	// アップロード処理
	s.makeWriteDir()
	s.addToIPFS()
	s.writeConfig()
}

func download() {

}
