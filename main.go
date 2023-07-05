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

type Settings struct {
	Output int
	Dev    bool
}

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
	settings := new(Settings)
	var window *walk.MainWindow

	MainWindow{
		Title:    "glvdc",
		Layout:   VBox{},
		AssignTo: &window,
		Size:     Size{242, 158},
		MinSize:  Size{242, 158},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 1, Alignment: AlignHCenterVCenter},
				Children: []Widget{
					Composite{
						Layout: Grid{Columns: 1, Alignment: AlignHCenterVCenter},
						Children: []Widget{
							Label{Text: "-----glvdc-----"},
							Label{Text: "-----by: me---"},
						},
					},
					PushButton{
						Text: "upload",
						OnClicked: func() {
							DropFiles(window)
						},
					},
					PushButton{
						Text: "settings",
						OnClicked: func() {
							//walk.App().Exit(0)
							SettingsDialog(window, settings)
						},
					},
				},
			},
		},
	}.Run()
}

type Formats struct {
	Id   int
	Name string
}

func OutFormats() []*Formats {
	return []*Formats{
		{0, "file"},
		{1, "display"},
	}
}

func DropFiles(owner walk.Form) (int, error) {
	var dlg *walk.MainWindow
	var txt *walk.TextEdit

	return MainWindow{
		AssignTo: &dlg,
		Title:    "upload",
		Size:     Size{466, 215},
		MinSize:  Size{466, 215},
		Layout:   VBox{},
		OnDropFiles: func(files []string) {
			//txt.SetText(strings.Join(files, "\r\n"))
			txt.SetText(CheckFiles(files))
		},
		Children: []Widget{
			TextEdit{
				AssignTo: &txt,
				ReadOnly: true,
				Text:     "drop album audio here\r\n(mp3, wav, flac)",
			},
		},
	}.Run()
}

// returns problem files if any extensions are not mp3, wav, or flac
// otherwise displays upload successful
func CheckFiles(files []string) string {
	var e []string
	for _, element := range files {
		fmt.Println("Processing upload... " + element)
		if !strings.HasSuffix(element, ".mp3") && !strings.HasSuffix(element, ".wav") && !strings.HasSuffix(element, ".flac") {
			fmt.Println(">> Illegal extension")
			e = append(e, element) // add problem file to slice
		}
	}
	if len(e) == 0 {
		return "upload successful"
	} else {
		return "error uploading\r\nfiles not allowed:\r\n" + strings.Join(e, "\r\n")
	}
}

func SettingsDialog(owner walk.Form, settings *Settings) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder

	return Dialog{
		AssignTo: &dlg,
		DataBinder: DataBinder{
			AssignTo:   &db,
			Name:       "settings",
			DataSource: settings,
		},
		Title:   "Settings",
		Size:    Size{242, 158},
		MinSize: Size{242, 158},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "output format:",
					},
					ComboBox{
						Value:         Bind("Output"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         OutFormats(),
					},

					Label{
						Text: "dev mode:",
					},
					CheckBox{
						Checked: Bind("Dev"),
					},
				},
			},
		},
	}.Run(owner)
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
