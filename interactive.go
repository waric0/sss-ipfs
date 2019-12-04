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
				num := strconv.Itoa(pubKeyNum)
				fmt.Println("以下" + num + "個の公開鍵を使用します")
				for i := 0; i < pubKeyNum; i++ {
					index := strconv.Itoa(i + 1)
					fmt.Println(index + " : " + pubKeys[i])
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
func askShareManagers(pubKeys []string, shareNum int, minNum int) []int {

	pubKeyNum := len(pubKeys)
	remains := shareNum - pubKeyNum
	var manageShareNums []int

	fmt.Println("それぞれの公開鍵が担当するシェアの数を順に入力してください(担当する公開鍵のないシェアの数は閾値未満である必要があります)")
	for i := 0; i < pubKeyNum; i++ {
		stdin := bufio.NewScanner(os.Stdin)
		max := remains + 1
		min := remains + 1 - (minNum - 1)
		if i != pubKeyNum-1 || 1 > min {
			min = 1
		}
		fmt.Print(pubKeys[i] + "(" + strconv.Itoa(min) + "以上かつ" + strconv.Itoa(max) + "以下) : ")
		for stdin.Scan() {
			manageShareNum, err := strconv.Atoi(stdin.Text())
			if err != nil {
				log.Fatal(err)
			}
			if manageShareNum >= min && max >= manageShareNum && manageShareNum > 0 {
				remains -= (manageShareNum - 1)
				manageShareNums = append(manageShareNums, manageShareNum)
				break
			}
			fmt.Print("正しい数を入力してください(" + strconv.Itoa(min) + "以上かつ" + strconv.Itoa(max) + "以下) : ")
		}
	}

	return manageShareNums
}
