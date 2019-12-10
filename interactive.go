package main

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// 暗号化で利用する公開鍵のkeyファイルを入力
func (s *uploadSetting) askPubKeys() {

	var pubKeyNum int = len(s.managers)

	stdin := bufio.NewScanner(os.Stdin)

	fmt.Print("1番目の公開鍵ファイルのパスを入力してください(doneで終了) : ")
	for stdin.Scan() {
		filePath := stdin.Text()
		if filePath == "done" {
			if pubKeyNum == 0 {
				fmt.Println("公開鍵が1つも設定されていません")
			} else {
				num := strconv.Itoa(pubKeyNum)
				fmt.Println("以下" + num + "個の公開鍵ファイルを使用します")
				for i := 0; i < pubKeyNum; i++ {
					pubKey := s.managers[i].fileName
					fmt.Println(pubKey)
				}
				break
			}
		} else {
			fileStat, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("公開鍵ファイルの読み込みに失敗しました")
			} else {
				file, err := ioutil.ReadFile(filePath)
				if err != nil {
					log.Fatal(err)
				}
				block, _ := pem.Decode(file)
				pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					log.Fatal(err)
				}
				pubKey, _ := pubKeyInterface.(*rsa.PublicKey)
				manager := keyManager{fileName: fileStat.Name(), publicKey: pubKey, manageShareNum: 0}
				s.managers = append(s.managers, manager)
				pubKeyNum = len(s.managers)
			}
		}
		index := strconv.Itoa(pubKeyNum + 1)
		fmt.Print(index + "番目の公開鍵ファイルのパスを入力してください(doneで終了) : ")
	}
}

// 秘密分散法で利用するシェアの数を入力
func (s *uploadSetting) askShareNum() {

	var (
		pubKeyNum int = len(s.managers)
		err       error
	)
	stdin := bufio.NewScanner(os.Stdin)

	if 2 > pubKeyNum {
		fmt.Print("シェアの数を入力してください(2以上) : ")
	} else {
		fmt.Print("シェアの数を入力してください(公開鍵の数" + strconv.Itoa(pubKeyNum) + "以上) : ")
	}
	for stdin.Scan() {
		s.shareNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if s.shareNum >= pubKeyNum && s.shareNum > 1 {
			break
		}
		if 2 > pubKeyNum {
			fmt.Print("正しいシェアの数を入力してください(2以上) : ")
		} else {
			fmt.Print("正しいシェアの数を入力してください(公開鍵の数" + strconv.Itoa(pubKeyNum) + "以上) : ")
		}
	}
}

// 秘密分散法で利用する閾値の数を入力
func (s *uploadSetting) askMinNum() {

	var err error
	stdin := bufio.NewScanner(os.Stdin)

	fmt.Print("閾値を入力してください(2以上かつシェアの数" + strconv.Itoa(s.shareNum) + "以下) : ")
	for stdin.Scan() {
		s.minNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if s.shareNum >= s.minNum && s.minNum > 1 {
			break
		}
		fmt.Print("正しい閾値を入力してください(2以上かつシェアの数" + strconv.Itoa(s.shareNum) + "以下) : ")
	}
}

// 各公開鍵の担当シェアを入力
func (s *uploadSetting) askShareManagers() {

	var (
		pubKeyNum int = len(s.managers)
		remains   int = s.shareNum - pubKeyNum
	)

	fmt.Println("それぞれの公開鍵が担当するシェアの数を順に入力してください(担当する公開鍵のないシェアの数は閾値未満である必要があります)")
	for i := 0; i < pubKeyNum; i++ {
		stdin := bufio.NewScanner(os.Stdin)
		max := remains + 1
		min := remains + 1 - (s.minNum - 1)
		if i != pubKeyNum-1 || 1 > min {
			min = 1
		}
		fmt.Print(s.managers[i].fileName + "(" + strconv.Itoa(min) + "以上かつ" + strconv.Itoa(max) + "以下) : ")
		for stdin.Scan() {
			input, err := strconv.Atoi(stdin.Text())
			if err != nil {
				log.Fatal(err)
			}
			if input >= min && max >= input && input > 0 {
				remains -= (input - 1)
				s.managers[i].manageShareNum = input
				break
			}
			fmt.Print("正しい数を入力してください(" + strconv.Itoa(min) + "以上かつ" + strconv.Itoa(max) + "以下) : ")
		}
	}
}

// アップロードするファイルを入力
func (s *uploadSetting) askFilePath() {

	stdin := bufio.NewScanner(os.Stdin)

	fmt.Print("アップロードするファイルのパスを入力してください : ")
	for stdin.Scan() {
		s.readFilePath = stdin.Text()
		_, err := os.Stat(s.readFilePath)
		if err == nil {
			break
		}
		fmt.Print("正しいファイルのパスを入力してください : ")
	}
}
