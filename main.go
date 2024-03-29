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
	keyfileName    string
	publicKey      *rsa.PublicKey
	privateKey     *rsa.PrivateKey
	manageShareNum int
	config         configuration
}

type commonSetting struct {
	readFilePath string
	tempDirPath  string
	writeDirPath string
	shares       []string
}

type uploadSetting struct {
	managers       []keyManager
	shareNum       int
	minNum         int
	cipherShareNum int
	comSet         commonSetting
}

type downloadSetting struct {
	manager keyManager
	comSet  commonSetting
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
	s.comSet.askFilePath("upload")

	// 加工処理
	fmt.Printf("秘密分散法を実行中\n")
	s.comSet.makeTempDir()
	s.sssaCreate()
	fmt.Printf("秘密分散法が完了\n")
	fmt.Printf("暗号化を実行中\n")
	s.encrypt()
	fmt.Printf("\n暗号化が完了\n")

	// アップロード処理
	fmt.Printf("IPFSへのアップロードを実行中\n")
	s.comSet.makeWriteDir()
	s.addToIPFS()
	s.writeConfig()
	fmt.Printf("\nIPFSへのアップロードが完了\n")
}

func download() {

	var s downloadSetting

	// 初期設定処理
	s.askPrivKeys()
	s.comSet.askFilePath("download")

	// ダウンロード処理
	fmt.Printf("IPFSからのダウンロードを実行中\n")
	s.comSet.makeTempDir()
	s.readConfig()
	s.getFromIPFS()
	fmt.Printf("\nIPFSからのダウンロードが完了\n")

	// 復元処理
	fmt.Printf("復号を実行中\n")
	s.decrypt()
	fmt.Printf("\n復号が完了\n")
	fmt.Printf("秘密分散法の復元を実行中\n")
	s.comSet.makeWriteDir()
	s.sssaCombine()
	fmt.Printf("秘密分散法の復元が完了\n")
}
