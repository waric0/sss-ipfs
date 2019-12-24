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

func BenchmarkUpload(b *testing.B) {

	var s uploadSetting

	s.initForBench()

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

func BenchmarkSssaCreate(b *testing.B) {

	var s uploadSetting

	s.initForBench()
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

	s.initForBench()
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

	s.initForBench()
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

func (s *uploadSetting) initForBench() {

	var (
		filePath = []string{
			"./test-keys/a-pub.pem",
			"./test-keys/b-pub.pem",
			"./test-keys/c-pub.pem"}
		ddArgs = []string{
			"if=/dev/zero",
			"of=",
			"bs=10",
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
	s.shareNum = len(filePath)

	// askMinNum()
	s.minNum = len(filePath)

	// askShareManagers()
	for i := 0; i < len(filePath); i++ {
		s.managers[i].manageShareNum = 1
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
	if _, err := os.Stat("temp"); os.IsExist(err) {
		if err := os.RemoveAll("temp"); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat("outputs"); os.IsExist(err) {
		if err := os.RemoveAll("outputs"); err != nil {
			log.Fatal(err)
		}
	}
}
