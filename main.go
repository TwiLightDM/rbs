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

// createWebsitesFromFile Создание файлов с расширением html при существовании сайта в нужную директорию
func createWebsitesFromFile(website string, dirPath string) error {
	site := "http://" + website
	resp, err := http.Get(site)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, website+".html")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fg := flag.NewFlagSet("check", flag.ExitOnError)
	txtPath := flag.String("input", "", "")
	dirPath := flag.String("output", "", "")
	flag.Parse()
	if fg.NFlag() != 2 {
		fmt.Println("Для того, чтобы пользоваться мной, укажите, где находится файл формата txt через -input. \nТакже укажите директорию, куда сохранить файл через -output.")
		return
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
		go func() {

			err := createWebsitesFromFile(line, *dirPath)
			if err != nil {
				fmt.Println("Обработана ошибка", err)
			}
			threads.Done()
		}()
	}

	threads.Wait()
	fmt.Println(time.Since(timer))
}
