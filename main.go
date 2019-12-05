package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/SSSaaS/sssa-golang"
)

type keyManager struct {
	publicKey      string
	manageShareNum int
}

func main() {
	flag.Parse()
	commands := flag.Arg(0)
	filepath := flag.Arg(1)

	if commands == "upload" {
		upload(filepath)
	} else if commands == "download" {
		download(filepath)
	} else {
		fmt.Printf("エラー : 適切なコマンドを入力してください\n")
		fmt.Printf("例 : \n")
		fmt.Printf("  sss-ipfs upload <filepath>\n")
		fmt.Printf("  sss-ipfs download <filepath>\n")
	}
}

func upload(filepath string) {

	var managers []keyManager

	managers = askPubKeys(managers)
	shareNum := askShareNum(managers)
	minNum := askMinNum(shareNum)
	managers = askShareManagers(managers, shareNum, minNum)

	fmt.Println(managers)

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
	created, err := sssa.Create(minNum, shareNum, string(raw))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		err = os.Mkdir("temp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < shareNum; i++ {
		content := []byte(created[i])
		index := strconv.Itoa(i + 1)
		err := ioutil.WriteFile("temp/share"+index, content, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func download(filepath string) {

}
