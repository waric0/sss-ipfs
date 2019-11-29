package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/SSSaaS/sssa-golang"
)

func main() {
	flag.Parse()
	filepath := flag.Arg(0)
	stdin := bufio.NewScanner(os.Stdin)
	var publicKeys []string
	var sharesNum, minimumNum, publicKeyNum int
	var err error

	fmt.Print("1番目の公開鍵を入力してください(doneで終了) : ")
	for stdin.Scan() {
		publicKey := stdin.Text()
		if err != nil {
			log.Fatal(err)
		}
		if publicKey == "done" {
			if publicKeyNum == 0 {
				fmt.Println("公開鍵が見つかりません")
			} else {
				fmt.Println("以下の公開鍵を使用します")
				for i := 0; i < publicKeyNum; i++ {
					fmt.Println(publicKeys[i])
				}
				break
			}
		} else {
			publicKeys = append(publicKeys, publicKey)
			publicKeyNum = len(publicKeys)
		}
		index := strconv.Itoa(publicKeyNum + 1)
		fmt.Print(index + "番目の公開鍵を入力してください(doneで終了) : ")
	}

	fmt.Print("シェアの数を入力してください(2以上かつ公開鍵の数以上) : ")
	for stdin.Scan() {
		sharesNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if sharesNum >= publicKeyNum && sharesNum > 1 {
			break
		}
		fmt.Print("正しいシェアの数を入力してください(2以上かつ公開鍵の数以上) : ")
	}

	fmt.Print("閾値を入力してください(2以上かつシェアの数以下) : ")
	for stdin.Scan() {
		minimumNum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if sharesNum >= minimumNum && minimumNum > 1 {
			break
		}
		fmt.Print("正しい閾値を入力してください(2以上かつシェアの数以下) : ")
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// シェアの作成
	created, err := sssa.Create(minimumNum, sharesNum, string(raw))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < sharesNum; i++ {
		content := []byte(created[i])
		index := strconv.Itoa(i + 1)
		err := ioutil.WriteFile("share"+index, content, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

}
