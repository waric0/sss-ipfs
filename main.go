package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
)

type configuration struct {
	EncryptedComKey []byte   `json:"encrypted_common_key"`
	ManagedShares   []string `json:"encrypted_share_cid"`
	UnmanagedShares []string `json:"unencrypted_share_cid"`
}

type keyManager struct {
	fileName       string
	publicKey      *rsa.PublicKey
	manageShareNum int
	config         configuration
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
	command := flag.Arg(0)

	switch command {
	case "upload":
		upload()
	case "download":
		download()
	default:
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
	fmt.Printf("秘密分散法を実行中\n")
	s.makeTempDir()
	s.sssaCreate()
	fmt.Printf("秘密分散法が完了\n")
	fmt.Printf("暗号化を実行中\n")
	s.encrypt()
	fmt.Printf("\n暗号化が完了\n")

	// アップロード処理
	fmt.Printf("IPFSへのアップロードを実行中\n")
	s.makeWriteDir()
	s.addToIPFS()
	s.writeConfig()
	fmt.Printf("\nIPFSへのアップロードが完了\n")
}

func download() {

}
