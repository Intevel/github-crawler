package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

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

func main() {
	links := []string{}

	for i := 1; i < 5; i++ {
		c := colly.NewCollector()
		c.OnHTML(`.v-align-middle`, func(e *colly.HTMLElement) {
			e.Request.Visit(e.Attr("href"))
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
			if !strings.HasPrefix(r.URL.String(), "https://github.com/search") {
				links = append(links, r.URL.String())
			}
		})

		c.Visit("https://github.com/search?q=javascript+language%3AJavaScript+license%3Amit&type=Repositories&ref=advsearch&l=JavaScript&l=&p=" + strconv.Itoa(i))
		continue
	}

	fmt.Println(links)
	writeArrayToFile(links, "links.txt")
	downloadWithGit(links)
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
