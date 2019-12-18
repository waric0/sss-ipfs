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

	fmt.Printf("1番目の公開鍵ファイルのパスを入力してください(doneで終了) : ")
	for stdin.Scan() {
		filePath := stdin.Text()
		if filePath == "done" {
			if pubKeyNum == 0 {
				fmt.Printf("公開鍵が1つも設定されていません\n")
			} else {
				num := strconv.Itoa(pubKeyNum)
				fmt.Printf("以下" + num + "個の公開鍵ファイルを使用します\n")
				for i := 0; i < pubKeyNum; i++ {
					pubKey := s.managers[i].keyfileName
					fmt.Printf("%s\n", pubKey)
				}
				break
			}
		} else {
			fileStat, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("公開鍵ファイルの読み込みに失敗しました\n")
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
				manager := keyManager{keyfileName: fileStat.Name(), publicKey: pubKey, manageShareNum: 0}
				s.managers = append(s.managers, manager)
				pubKeyNum = len(s.managers)
			}
		}
		fmt.Printf("%d番目の公開鍵ファイルのパスを入力してください(doneで終了) : ", pubKeyNum+1)
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
		fmt.Printf("シェアの数を入力してください(2以上) : ")
	} else {
		fmt.Printf("シェアの数を入力してください(公開鍵の数%d以上) : ", pubKeyNum)
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
			fmt.Printf("正しいシェアの数を入力してください(2以上) : ")
		} else {
			fmt.Printf("正しいシェアの数を入力してください(公開鍵の数%d以上) : ", pubKeyNum)
		}
	}
}

// 秘密分散法で利用する閾値の数を入力
func (s *uploadSetting) askMinNum() {

	var err error
	stdin := bufio.NewScanner(os.Stdin)

	fmt.Printf("閾値を入力してください(2以上かつシェアの数%d以下) : ", s.shareNum)
	for stdin.Scan() {
		s.minNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if s.shareNum >= s.minNum && s.minNum > 1 {
			break
		}
		fmt.Printf("正しい閾値を入力してください(2以上かつシェアの数%d以下) : ", s.shareNum)
	}
}

// 各公開鍵の担当シェアを入力
func (s *uploadSetting) askShareManagers() {

	var (
		pubKeyNum int = len(s.managers)
		remains   int = s.shareNum - pubKeyNum
	)

	fmt.Printf("それぞれの公開鍵が担当するシェアの数を順に入力してください(担当する公開鍵のないシェアの数は閾値未満である必要があります)\n")
	for i := 0; i < pubKeyNum; i++ {
		stdin := bufio.NewScanner(os.Stdin)
		max := remains + 1
		min := remains + 1 - (s.minNum - 1)
		if i != pubKeyNum-1 || 1 > min {
			min = 1
		}
		fmt.Printf("%s(%d以上かつ%d以下) : ", s.managers[i].keyfileName, min, max)
		for stdin.Scan() {
			input, err := strconv.Atoi(stdin.Text())
			if err != nil {
				log.Fatal(err)
			}
			if input >= min && max >= input && input > 0 {
				remains -= (input - 1)
				s.managers[i].manageShareNum = input
				s.cipherShareNum += input
				break
			}
			fmt.Printf("正しい数を入力してください(%d以上かつ%d以下) : ", min, max)
		}
	}
}

// アップロードするファイルを入力
func (s *uploadSetting) askFilePath() {

	stdin := bufio.NewScanner(os.Stdin)

	fmt.Printf("アップロードするファイルのパスを入力してください : ")
	for stdin.Scan() {
		s.comSet.readFilePath = stdin.Text()
		_, err := os.Stat(s.comSet.readFilePath)
		if err == nil {
			break
		}
		fmt.Printf("正しいファイルのパスを入力してください : ")
	}
}
