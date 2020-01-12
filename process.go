package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var dirIn = os.Args[1]
var dirOut = os.Args[2]

var wg sync.WaitGroup

func check(e error) bool {
	if e != nil {
		if e == io.EOF {
			return true
		}
		log.Fatal(e)
	}
	return false
}

func readWriteAsync(name string) {
	defer wg.Done()

	const BufferSize = 20

	file, err := os.Open(dirIn + "/" + name)
	check(err)
	defer file.Close()

	buffer := make([]byte, BufferSize)
	hasher := md5.New()

	for {
		bufferSize, err := file.Read(buffer)

		if check(err) {
			break
		}

		hasher.Write(buffer[:bufferSize])
	}

	file, err = os.Create(dirOut + "/" + name + ".res")
	check(err)
	defer file.Close()

	_, err = file.WriteString(hex.EncodeToString(hasher.Sum(nil)))
	check(err)
}

func main() {
	files, err := ioutil.ReadDir(dirIn)
	check(err)
	if _, err := os.Stat(dirOut); os.IsNotExist(err) {
		os.Mkdir(dirOut, os.ModeDir)
	}

	count := 0
	for _, file := range files {
		wg.Add(1)
		go readWriteAsync(file.Name())
		count++
	}

	wg.Wait()

	fmt.Println("Total number of processed files:", count)
}
