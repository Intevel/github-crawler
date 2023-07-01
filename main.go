package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

type DataEntry struct {
	Repo     string `json:"repo"`
	Language string `json:"language"`
	Data     string `json:"data"`
	Filename string `json:"filename"`
}

var dataEntrys = make([]DataEntry, 0)

func main() {
	// links := startCrawler()
	// downloadWithGit(links)

	count, _ := readJsFilesInFolder("/home/conner/github-crawler/cloned/")
	fmt.Println("Read " + strconv.Itoa(count) + " JavaScript files.")
	// fmt.Println(dataEntrys)
	jsonData, err := json.MarshalIndent(dataEntrys, "", "  ")
	if err != nil {
		fmt.Println("Fehler beim Konvertieren des Structs in JSON:", err)
		return
	}

	err = ioutil.WriteFile("data.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Fehler beim Schreiben der JSON-Datei:", err)
		return
	}
}

func writeArrayToFile(arr []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("fehler beim Erstellen der Datei: %v", err)
	}
	defer file.Close()

	for _, item := range arr {
		_, err := file.WriteString(item + "\n")
		if err != nil {
			return fmt.Errorf("fehler beim Schreiben der Zeile: %v", err)
		}
	}

	fmt.Printf("Der Array wurde erfolgreich in die Datei '%s' geschrieben.\n", filename)
	return nil
}

func startCrawler() (links []string) {
	linkList := []string{}

	for i := 1; i < 5; i++ {
		c := colly.NewCollector()
		c.OnHTML(`.v-align-middle`, func(e *colly.HTMLElement) {
			e.Request.Visit(e.Attr("href"))
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
			if !strings.HasPrefix(r.URL.String(), "https://github.com/search") {
				linkList = append(linkList, r.URL.String())
			}
		})

		c.Visit("https://github.com/search?q=javascript+language%3AJavaScript+license%3Amit&type=Repositories&ref=advsearch&l=JavaScript&l=&p=" + strconv.Itoa(i))
		continue
	}

	fmt.Println(linkList)
	writeArrayToFile(linkList, "links.txt")
	return linkList
}

func downloadWithGit(links []string) {
	for _, link := range links {
		uuid := uuid.New()
		cmd := exec.Command("git", "clone", link+".git", uuid.String())
		cmd.Dir = "/home/conner/github-crawler/cloned"
		_, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Cloned: " + link + ".git" + " to " + uuid.String())
		cmd.Run()
	}
}

func readJsFilesInFolder(folder string) (int, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return 0, fmt.Errorf("kann den Ordner nicht öffnen: %w", err)
	}

	count := 0

	for _, file := range files {
		if file.IsDir() {
			subFolderPath := filepath.Join(folder, file.Name())
			subCount, err := readJsFilesInFolder(subFolderPath)
			if err != nil {
				fmt.Printf("error reading sub folder %s: %v\n", subFolderPath, err)
			}
			count += subCount
		} else if filepath.Ext(file.Name()) == ".js" {
			count += 1
			filePath := filepath.Join(folder, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Printf("error reading file %s: %v\n", filePath, err)
				continue
			}

			fmt.Printf("content of %s:\n%s\n", filePath, string(content))
			count++
			dataEntrys = append(dataEntrys, DataEntry{Repo: folder, Language: "JavaScript", Data: string(content), Filename: file.Name()})
		}
	}
	return count, nil
}
