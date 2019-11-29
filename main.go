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
	var shares, minimum int
	var err error

	fmt.Print("シェアの数を入力してください(2以上) : ")
	for stdin.Scan() {
		shares, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if shares > 1 {
			break
		}
		fmt.Print("正しいシェアの数を入力してください(2以上) : ")
	}

	fmt.Print("閾値を入力してください(2以上かつシェアの数以下) : ")
	for stdin.Scan() {
		minimum, err = strconv.Atoi(stdin.Text())
		if err != nil {
			log.Fatal(err)
		}
		if shares >= minimum && minimum > 1 {
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
	created, err := sssa.Create(minimum, shares, string(raw))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < shares; i++ {
		content := []byte(created[i])
		index := strconv.Itoa(i + 1)
		err := ioutil.WriteFile("share"+index, content, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

}
