package main

import(
	"fmt"
	"bufio"
	"os"
	"io"
	"strings"

	// scraper
	"github.com/gocolly/colly"
	// time.Duration is stupid. This helps convert hh:mm:ss to duration obj
	"github.com/dannav/hhmmss"
)

type song struct {
	Name string
	Time string
}

func main() {
	var(
		artist string 
		album string
		year string
		name string
		time string
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
		e.ForEach("div.title", func(_ int, el *colly.HTMLElement) {
			name = el.ChildText(".track-title") // track name
			time = el.ChildText(".time.secondaryText") // track time
			n = append(n, song{name, time})
		})

		e.ForEach("#name-section", func(_ int, el *colly.HTMLElement) {
			album = el.ChildText(".trackTitle") // album name
			artist = el.ChildText("h3 span") // artist name
		})

		e.ForEach(".tralbumData.tralbum-credits", func(_ int, el *colly.HTMLElement) {
			date,_,_ := strings.Cut(strings.TrimSpace(el.Text), "\n") // full date (released Jan 01, 20xx)
			year = date[len(date)-4:] // grab year
		})
	})

	c.Visit(url) // visit, get results

	write(artist, album, year, ctn, n)
}

// this function writes to the output.txt file and echoes
// its output to the cli
func write(ar string, al string, yr string, cat string, n []song) {
	f, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	w := io.MultiWriter(os.Stdout, f)
	defer f.Close() // close at finish

	fmt.Println("\n-----------------------\n")

	// write header
	_, err = fmt.Fprintf(w, 
	"Label: Geometric Lullaby\n" +
	"Artist: %s\n" +
	"Album: %s\n" +
	"Year: %s\n" +
	"(GEO - %s)\n",
	ar, al, yr, cat)

	_, err = fmt.Fprint(w,
	"\nDownload for free or purchase at: https://geometriclullaby.bandcamp.com/\n",
	)

	// save first track info and remove from slice
	var(
		fname string = n[0].Name
		ftime string = n[0].Time
	)
	n = n[1:] // remove first track info

	_, err = fmt.Fprintf(w,
	"\n00:00 - %s\n",
	fname)

	dur, _ := hhmmss.Parse(fc(ftime))
	for _, element := range n {
		// convert float
		s := fmt.Sprintf("%.0f", dur.Seconds())
		m := fmt.Sprintf("%.0f", dur.Minutes())
		h := fmt.Sprintf("%.0f", dur.Hours())
		
		_, err = fmt.Fprintf(w,
		"%s:%s:%s - %s\n",
		h, m, s, element.Name,
		)
		ftime = element.Time
		nextdur, _ := hhmmss.Parse(fc(ftime))
		dur += nextdur
	}

	if err != nil {
		panic(err)
	}
}

// this function is here because sometimes if a track is long,
// bandcamp formats in hh:mm:ss instead of mm:ss
func fc(time string) string { // fc == format check
	x := strings.Count(time, ":")
	if x == 1 {
		return ("00:" + time)
	} else if x == 2 {
		return time
	} else {
		fmt.Println("ERROR examining time format... character ':' count == ", x)
		panic("Execution stopped")
	}
}

// prompts for url and cat no
// todo: these could be 1 function
// ------------------------------------
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