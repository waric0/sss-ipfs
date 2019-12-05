package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// 暗号化で利用する公開鍵の数を入力
func askPubKeys(managers []keyManager) []keyManager {

	var pubKeyNum int = len(managers)

	stdin := bufio.NewScanner(os.Stdin)

	fmt.Print("1番目の公開鍵を入力してください(doneで終了) : ")
	for stdin.Scan() {
		input := stdin.Text()
		if input == "done" {
			if pubKeyNum == 0 {
				fmt.Println("公開鍵が見つかりません")
			} else {
				num := strconv.Itoa(pubKeyNum)
				fmt.Println("以下" + num + "個の公開鍵を使用します(冒頭15文字のみ表示)")
				for i := 0; i < pubKeyNum; i++ {
					index := strconv.Itoa(i + 1)
					pubKey := managers[i].publicKey
					if len(pubKey) > 15 {
						fmt.Println("PublicKey" + index + " : " + pubKey[:15] + "...")
					} else {
						fmt.Println("PublicKey" + index + " : " + pubKey)
					}
				}
				break
			}
		} else {
			manager := keyManager{publicKey: input, manageShareNum: 0}
			managers = append(managers, manager)
			pubKeyNum = len(managers)
		}
		index := strconv.Itoa(pubKeyNum + 1)
		fmt.Print(index + "番目の公開鍵を入力してください(doneで終了) : ")
	}

	return managers
}

// 秘密分散法で利用するシェアの数を入力
func askShareNum(managers []keyManager) int {

	var (
		pubKeyNum int = len(managers)
		shareNum  int
		err       error
	)
	stdin := bufio.NewScanner(os.Stdin)

	if 2 > pubKeyNum {
		fmt.Print("シェアの数を入力してください(2以上) : ")
	} else {
		fmt.Print("シェアの数を入力してください(公開鍵の数" + strconv.Itoa(pubKeyNum) + "以上) : ")
	}
	for stdin.Scan() {
		shareNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if shareNum >= pubKeyNum && shareNum > 1 {
			break
		}
		if 2 > pubKeyNum {
			fmt.Print("正しいシェアの数を入力してください(2以上) : ")
		} else {
			fmt.Print("正しいシェアの数を入力してください(公開鍵の数" + strconv.Itoa(pubKeyNum) + "以上) : ")
		}
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

	fmt.Print("閾値を入力してください(2以上かつシェアの数" + strconv.Itoa(shareNum) + "以下) : ")
	for stdin.Scan() {
		minNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if shareNum >= minNum && minNum > 1 {
			break
		}
		fmt.Print("正しい閾値を入力してください(2以上かつシェアの数" + strconv.Itoa(shareNum) + "以下) : ")
	}

	return minNum
}

// 各公開鍵の担当シェアを入力
func askShareManagers(managers []keyManager, shareNum int, minNum int) []keyManager {

	var (
		pubKeyNum int = len(managers)
		remains   int = shareNum - pubKeyNum
	)

	fmt.Println("それぞれの公開鍵が担当するシェアの数を順に入力してください(担当する公開鍵のないシェアの数は閾値未満である必要があります)")
	for i := 0; i < pubKeyNum; i++ {
		stdin := bufio.NewScanner(os.Stdin)
		max := remains + 1
		min := remains + 1 - (minNum - 1)
		if i != pubKeyNum-1 || 1 > min {
			min = 1
		}
		fmt.Print(managers[i].publicKey + "(" + strconv.Itoa(min) + "以上かつ" + strconv.Itoa(max) + "以下) : ")
		for stdin.Scan() {
			input, err := strconv.Atoi(stdin.Text())
			if err != nil {
				log.Fatal(err)
			}
			if input >= min && max >= input && input > 0 {
				remains -= (input - 1)
				managers[i].manageShareNum = input
				break
			}
			fmt.Print("正しい数を入力してください(" + strconv.Itoa(min) + "以上かつ" + strconv.Itoa(max) + "以下) : ")
		}
	}

	return managers
}
