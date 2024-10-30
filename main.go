package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Создание файлов с расширением html при существовании сайта в нужную директорию
func createWebsitesFromFile(website string, dirPath string) {
	site := "http://" + website
	resp, err := http.Get(site)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return
	}

	filePath := filepath.Join(dirPath, website+".html")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
}

func main() {
	txtPath := flag.String("input", "", "")
	dirPath := flag.String("output", "", "")
	flag.Parse()

	txt, err := os.Open(*txtPath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла ", err)
		return
	}
	defer txt.Close()

	scanner := bufio.NewScanner(txt)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка чтения файла ", err)
		}
		createWebsitesFromFile(line, *dirPath)
		fmt.Println(line)
	}
}
