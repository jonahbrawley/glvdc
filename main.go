package main

import(
	"fmt"
	"bufio"
	"os"
	_ "io"
	"strings"

	// scraper
	"github.com/gocolly/colly"
)

var(
	year string
	artist string 
	album string
)

type song struct {
	Name string
	Time string
}

func main() {
	fmt.Println("-------- glvdc --------")
	fmt.Println("-------- by: me -------")
	fmt.Println("")
	url := urlprompt("bc url:")
	ctn := catprompt("cat no:")
	
	c := colly.NewCollector()
	n := []song{} // create array of struct

	c.OnHTML("div.title", func(e *colly.HTMLElement) {
		name := e.ChildText(".track-title") // track name
		time := e.ChildText(".time.secondaryText") // track time
		
		n = append(n, song{name, time})
	})

	c.Visit(url) // visit specified url

	fmt.Println("\nTracks:\n")
	// print (debug)
	for _, element := range n {
		fmt.Println("Name: ", element.Name)
		fmt.Println("Time: ", element.Time)
		fmt.Println("")
	}
	fmt.Println("GEO - ", ctn)
}

func write() {
	f, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// write
	_, err = fmt.Fprint(w, 
	`Label: Geometric Lullaby
	Artist: First Kings
	Album: Water Birth
	Year: 2018
	(GEO - )`)
}

// prompt for url
func urlprompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if strings.HasPrefix(s, "https://") {
			break;
		}
	}
	return strings.TrimSpace(s)
}

func catprompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break;
		}
	}
	return strings.TrimSpace(s)
}