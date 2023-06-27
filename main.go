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

type song struct {
	Name string
	Time string
}

func main() {
	var(
		//year string
		artist string 
		album string
	)

	fmt.Println("-------- glvdc --------")
	fmt.Println("-------- by: me -------")
	fmt.Println("")
	url := urlprompt("bc url:")
	ctn := catprompt("cat no:")
	
	c := colly.NewCollector()
	n := []song{} // create array of struct

	// callbacks
	c.OnHTML("body", func(e *colly.HTMLElement) {
		var(
			name string 
			time string 
		)
		e.ForEach("div.title", func(_ int, el *colly.HTMLElement) {
			name = el.ChildText(".track-title") // track name
			time = el.ChildText(".time.secondaryText") // track time
			n = append(n, song{name, time})
		})

		e.ForEach("#name-section", func(_ int, el *colly.HTMLElement) {
			album = el.ChildText(".trackTitle") // album name
			artist = el.ChildText("h3 span") // artist name
		})
	})

	/*c.OnHTML("h2.trackTitle", func(e *colly.HTMLElement) {
		album = e.ChildText("h2.trackTitle") // album name
		//artist = e.ChildText("a href") // artist name
	})

	c.OnHTML("div.tralbumData.tralbum-credits", func(e *colly.HTMLElement) {
		year = e.ChildText("#text") // full date
		//year = date[len(date)-4:] // release year
	})*/

	c.Visit(url) // visit bc url

	// print (debug)
	fmt.Println("\nAlbum: ", album)
	fmt.Println("Artist: ", artist)
	fmt.Println("\nTracks:")
	for _, element := range n {
		fmt.Println("Name: ", element.Name)
		fmt.Println("Time: ", element.Time)
		fmt.Println("")
	}
	fmt.Println("GEO - ", ctn)
	//fmt.Println("Year: ", year)
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