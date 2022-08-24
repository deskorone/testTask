package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	s := time.Now()
	path := flag.String("path", "/", "find path")
	dest := flag.String("dest", "/home/deskor/", "find name")
	flag.Parse()

	err := parceUrl(*dest, *path)

	if err != nil {
		return
	}

	e := time.Since(s)
	println(e.Seconds())

}

// берем из файла URL и перносим содержимое старницы в файлы
func parceUrl(dest, path string) error {

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	a := 0
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var mutex sync.Mutex

	for scanner.Scan() {
		s := scanner.Text()
		a++
		if s != "" {
			wg.Add(1)
			go createFile(dest, s, &mutex, a)
		}
	}
	wg.Wait()
	return nil
}

// Создание файла на основе URL
func createFile(dest, s string, m *sync.Mutex, a int) error {
	defer wg.Done()
	data, err := getDocFromUrl(s)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	newFile, err := os.Create(dest + "NewFile" + strconv.Itoa(a) + ".txt")

	if err != nil {
		fmt.Println(err.Error())
		newFile.Close()
		return err
	}
	newFile.Write(data)
	newFile.Close()
	return nil
}

// Получение документа по URL
func getDocFromUrl(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	return b, err
}
