package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// 暗号化で利用する公開鍵の数を入力
func askPubKeys() []string {

	var (
		pubKeys   []string
		pubKeyNum int
	)
	stdin := bufio.NewScanner(os.Stdin)

	fmt.Print("1番目の公開鍵を入力してください(doneで終了) : ")
	for stdin.Scan() {
		pubKey := stdin.Text()
		if pubKey == "done" {
			if pubKeyNum == 0 {
				fmt.Println("公開鍵が見つかりません")
			} else {
				fmt.Println("以下の公開鍵を使用します")
				for i := 0; i < pubKeyNum; i++ {
					fmt.Println(pubKeys[i])
				}
				break
			}
		} else {
			pubKeys = append(pubKeys, pubKey)
			pubKeyNum = len(pubKeys)
		}
		index := strconv.Itoa(pubKeyNum + 1)
		fmt.Print(index + "番目の公開鍵を入力してください(doneで終了) : ")
	}

	return pubKeys
}

// 秘密分散法で利用するシェアの数を入力
func askShareNum(pubKeyNum int) int {

	var (
		shareNum int
		err      error
	)
	stdin := bufio.NewScanner(os.Stdin)

	fmt.Print("シェアの数を入力してください(2以上かつ公開鍵の数以上) : ")
	for stdin.Scan() {
		shareNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if shareNum >= pubKeyNum && shareNum > 1 {
			break
		}
		fmt.Print("正しいシェアの数を入力してください(2以上かつ公開鍵の数以上) : ")
	}

	return shareNum
}

// 秘密分散法で利用する閾値の数を入力
func askMinNum(shareNum int) int {

	var (
		minNum int
		err    error
	)
	stdin := bufio.NewScanner(os.Stdin)

	fmt.Print("閾値を入力してください(2以上かつシェアの数以下) : ")
	for stdin.Scan() {
		minNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if shareNum >= minNum && minNum > 1 {
			break
		}
		fmt.Print("正しい閾値を入力してください(2以上かつシェアの数以下) : ")
	}

	return minNum
}
