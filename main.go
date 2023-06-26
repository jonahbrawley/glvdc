package main

import(
	"fmt"
	"bufio"
	"os"
	"strings"

	// scraper
	"github.com/gocolly/colly"
)

type song struct {
	Name string
	Time string
}

func main() {
	fmt.Println("-------- glvdc --------")
	fmt.Println("-------- by: me -------")
	fmt.Println("")
	url := prompt("bc url:")
	
	c := colly.NewCollector()
	n := []song{} // create array of struct

	c.OnHTML("div.title", func(e *colly.HTMLElement) {
		name := e.ChildText(".track-title") // track name
		time := e.ChildText(".time.secondaryText") // track time
		
		n = append(n, song{name, time})
	})

	c.Visit(url) // visit specified url

	// print (debug)
	for _, element := range n {
		fmt.Println("Name: ", element.Name)
		fmt.Println("Time: ", element.Time)
		fmt.Println("")
	}
}

// prompt for url
func prompt(label string) string {
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