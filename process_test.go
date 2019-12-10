package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/SSSaaS/sssa-golang"
)

var shareNum int = 3
var ddArgs = []string{
	"if=/dev/zero",
	"of=./temp/sample_file",
	"bs=100",
	"count=1",
}

func BenchmarkSssa(b *testing.B) {
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		err = os.Mkdir("temp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := exec.Command("dd", ddArgs...).Run()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open("./temp/sample_file")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sssa.Create(shareNum, shareNum, string(raw))
		if err != nil {
			log.Fatal(err)
		}
	}
	b.StopTimer()
}

func BenchmarkWriteFile(b *testing.B) {
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		err = os.Mkdir("temp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := exec.Command("dd", ddArgs...).Run()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open("./temp/sample_file")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}
	created, err := sssa.Create(shareNum, shareNum, string(raw))
	if err != nil {
		log.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(created); j++ {
			index := strconv.Itoa(j + 1)
			err = ioutil.WriteFile("temp/un_managed_share"+index, []byte(created[j]), 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	b.StopTimer()
}
