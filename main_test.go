package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
)

func BenchmarkSssaCreate(b *testing.B) {

	var s uploadSetting

	s.initForUpBench()
	s.comSet.makeTempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.sssaCreate()
	}
	b.StopTimer()

	s = uploadSetting{}
}

func BenchmarkEncrypt(b *testing.B) {

	var s uploadSetting

	s.initForUpBench()
	s.comSet.makeTempDir()
	s.sssaCreate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.encrypt()
	}
	b.StopTimer()

	s = uploadSetting{}
}

func BenchmarkAddToIPFS(b *testing.B) {

	var s uploadSetting

	s.initForUpBench()
	s.comSet.makeTempDir()
	s.sssaCreate()
	s.encrypt()
	s.comSet.makeWriteDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.addToIPFS()
	}
	b.StopTimer()

	s = uploadSetting{}
}

func BenchmarkUpload(b *testing.B) {

	var s uploadSetting

	s.initForUpBench()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.comSet.makeTempDir()
		s.sssaCreate()
		s.encrypt()
		s.comSet.makeWriteDir()
		s.addToIPFS()
		s.writeConfig()
	}
	b.StopTimer()

	s = uploadSetting{}
}

func (s *uploadSetting) initForUpBench() {

	var (
		shareNum int = 3
		filePath     = []string{"./test-keys/a-pub.pem"}
		ddArgs       = []string{
			"if=/dev/urandom",
			"of=",
			"bs=100",
			"count=1",
		}
	)

	// askPubKeys()
	for i := 0; i < len(filePath); i++ {
		file, err := ioutil.ReadFile(filePath[i])
		if err != nil {
			log.Fatal(err)
		}
		block, _ := pem.Decode(file)
		pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		pubKey, _ := pubKeyInterface.(*rsa.PublicKey)
		fileStat, err := os.Stat(filePath[i])
		manager := keyManager{keyfileName: fileStat.Name(), publicKey: pubKey, manageShareNum: 0}
		s.managers = append(s.managers, manager)
	}

	// askShareNum()
	s.shareNum = shareNum

	// askMinNum()
	s.minNum = shareNum

	// askShareManagers()
	for i := 0; i < len(filePath); i++ {
		s.managers[i].manageShareNum = shareNum
		s.cipherShareNum += s.managers[i].manageShareNum
	}

	// askFilePath()
	s.comSet.readFilePath = "./test-keys/sample_file"

	// 元データの作成
	ddArgs[1] = ddArgs[1] + s.comSet.readFilePath
	err := exec.Command("dd", ddArgs...).Run()
	if err != nil {
		log.Fatal(err)
	}

	// ディレクトリの削除
	if _, err = os.Stat("temp"); err == nil {
		if err := os.RemoveAll("temp"); err != nil {
			log.Fatal(err)
		}
	}
	if _, err = os.Stat("outputs"); err == nil {
		if err := os.RemoveAll("outputs"); err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkGetFromIPFS(b *testing.B) {

	var s downloadSetting

	s.initForDownBench()
	s.comSet.makeTempDir()
	s.readConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.getFromIPFS()
	}
	b.StopTimer()

	s = downloadSetting{}
}

func BenchmarkDecrypt(b *testing.B) {

	var s downloadSetting

	s.initForDownBench()
	s.comSet.makeTempDir()
	s.readConfig()
	s.getFromIPFS()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.dencrypt()
	}
	b.StopTimer()

	s = downloadSetting{}
}

func BenchmarkSssaCombine(b *testing.B) {

	var s downloadSetting

	s.initForDownBench()
	s.comSet.makeTempDir()
	s.readConfig()
	s.getFromIPFS()
	s.dencrypt()
	s.comSet.makeWriteDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.sssaCombine()
	}
	b.StopTimer()

	s = downloadSetting{}
}

func BenchmarkDownload(b *testing.B) {

	var s downloadSetting

	s.initForDownBench()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.comSet.makeTempDir()
		s.readConfig()
		s.getFromIPFS()
		s.dencrypt()
		s.comSet.makeWriteDir()
		s.sssaCombine()
	}
	b.StopTimer()

	s = downloadSetting{}
}

func (s *downloadSetting) initForDownBench() {

	var filePath string = "./test-keys/a.pem"

	// askPrivKeys()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(file)
	var privKey *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
	default:
		privKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		privKey, _ = privKeyInterface.(*rsa.PrivateKey)
	}
	fileStat, err := os.Stat(filePath)
	s.manager = keyManager{keyfileName: fileStat.Name(), privateKey: privKey}

	// askFilePath()
	s.comSet.readFilePath = "./outputs/a-pub_pem/config.json"

	// ディレクトリの削除
	if err := os.RemoveAll("temp"); err != nil {
		log.Fatal(err)
	}
	if _, err = os.Stat("outputs/content"); err == nil {
		if err := os.Remove("outputs/content"); err != nil {
			log.Fatal(err)
		}
	}
}
