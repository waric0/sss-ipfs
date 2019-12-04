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

func main() {
	flag.Parse()
	filepath := flag.Arg(0)

	pubKeys := askPubKeys()
	pubKeyNum := len(pubKeys)
	shareNum := askShareNum(pubKeyNum)
	minNum := askMinNum(shareNum)
	manageShareNums := askShareManagers(pubKeys, shareNum, minNum)

	fmt.Println(manageShareNums)

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
