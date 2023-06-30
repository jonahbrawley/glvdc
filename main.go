package main

import (
	"bufio"
	"fmt"
	"io"
	_ "math"
	"os"
	_ "strconv"
	"strings"
	"time"

	"github.com/dannav/hhmmss"          // convert hh:mm:ss to duration()
	_ "github.com/gocolly/colly"        // scraper
	"github.com/lxn/walk"               // frontend
	. "github.com/lxn/walk/declarative" // dependency
	_ "github.com/lxn/win"              // dependency
)

type song struct {
	Name string
	Time string
}

var (
	window *walk.MainWindow
	width  int = 320
	height int = 240
	// x      int = 200
	// y      int = 200
)

// CLI MAIN
//
/*func main() {
	var (
		artist string
		album  string
		year   string
		name   string
		time   string
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
			name = el.ChildText(".track-title")        // track name
			time = el.ChildText(".time.secondaryText") // track time
			n = append(n, song{name, time})
		})

		e.ForEach("#name-section", func(_ int, el *colly.HTMLElement) {
			album = el.ChildText(".trackTitle") // album name
			artist = el.ChildText("h3 span")    // artist name
		})

		e.ForEach(".tralbumData.tralbum-credits", func(_ int, el *colly.HTMLElement) {
			date, _, _ := strings.Cut(strings.TrimSpace(el.Text), "\n") // full date (released Jan 01, 20xx)
			year = date[len(date)-4:]                                   // grab year
		})
	})

	c.Visit(url) // visit, get results

	write(artist, album, year, ctn, n)
}*/

func main() {
	MainWindow{
		Title:  "glvdc",
		Layout: VBox{},
		//AssignTo: &window,
		Size: Size{320, 240},
		Children: []Widget{
			Label{Text: "----glvdc----"},
			PushButton{
				Text: "Upload",
				OnClicked: func() {
					walk.App().Exit(0)
				},
			},
			PushButton{
				Text: "Quit",
				OnClicked: func() {
					walk.App().Exit(0)
				},
			},
		},
	}.Run()
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
		"Label: Geometric Lullaby\n"+
			"Artist: %s\n"+
			"Album: %s\n"+
			"Year: %s\n"+
			"(GEO - %s)\n",
		ar, al, yr, cat)

	_, err = fmt.Fprint(w,
		"\nDownload for free or purchase at: https://geometriclullaby.bandcamp.com/\n",
	)

	// save first track info and remove from slice
	var (
		fname string = n[0].Name
		ftime string = n[0].Time
	)
	n = n[1:] // remove first track info

	_, err = fmt.Fprintf(w,
		"\n00:00:00 - %s\n",
		fname)

	ftparsed, _ := hhmmss.Parse(fc(ftime))
	added_dur := time.Duration(int64(ftparsed/time.Second)) * time.Second

	for _, element := range n {
		_, err = fmt.Fprintf(w,
			fmtDuration(added_dur)+" - %s\n",
			element.Name,
		)
		if err != nil {
			panic(err)
		}
		cdur, _ := hhmmss.Parse(fc(element.Time)) // GRAB NEXT TRACK DURATION
		// ADD TO added_dur
		added_dur = time.Duration(int64(added_dur/time.Second)+int64(cdur/time.Second)) * time.Second
		//fmt.Println(added_dur) // DEBUG
	}
}

func fmtDuration(d time.Duration) string {
	//d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute // new
	s := d / time.Second // new
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// sometimes if a track is long,
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
			break
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
			break
		}
	}
	return strings.TrimSpace(s)
}
