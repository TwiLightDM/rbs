package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var errors = []error{}

// createWebsitesFromFile Создание файлов с расширением html при существовании сайта в нужную директорию
func createWebsitesFromFile(website string, dirPath string, threads *sync.WaitGroup) {
	defer threads.Done()
	site := "http://" + website
	resp, err := http.Get(site)
	if err != nil {
		errors = append(errors, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errors = append(errors, err)
		return
	}

	filePath := filepath.Join(dirPath, website+".html")
	file, err := os.Create(filePath)
	if err != nil {
		errors = append(errors, err)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		errors = append(errors, err)
		return
	}
}

func main() {
	txtPath := flag.String("input", "", "")
	dirPath := flag.String("output", "", "")
	flag.Parse()
	if *txtPath == "" || *dirPath == "" {
		fmt.Println("Для того, чтобы пользоваться мной, укажите, где находится файл формата txt через -input. \nТакже укажите директорию, куда сохранить файл через -output.")

	}

	timer := time.Now()

	txt, err := os.Open(*txtPath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла ", err)
		return
	}
	defer txt.Close()

	if err = os.MkdirAll(*dirPath, os.ModePerm); err != nil {
		fmt.Println("Ошибка при создании директории:", err)
		return
	}

	var threads sync.WaitGroup

	scanner := bufio.NewScanner(txt)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка чтения файла ", err)
		}
		threads.Add(1)
		go createWebsitesFromFile(line, *dirPath, &threads)

		fmt.Println(line)
	}

	threads.Wait()

	if len(errors) > 0 {
		fmt.Println("Обработаны следующие ошибки:")
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Обработка завершена без ошибок.")
	}

	fmt.Println(time.Since(timer))
}
