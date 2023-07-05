package main

import (
	"bufio"
	"fmt"
	"io"
	_ "math"
	"os"
	"strings"
	"time"

	"github.com/dannav/hhmmss"          // convert hh:mm:ss to duration()
	"github.com/gocolly/colly"          // scraper
	"github.com/hajimehoshi/go-mp3"     // parse mp3
	"github.com/lxn/walk"               // frontend
	. "github.com/lxn/walk/declarative" // dependency
	_ "github.com/lxn/win"              // dependency
)

type Album struct {
	Link   string
	Artist string
	Name   string
	Year   string
	CatNo  string
}

type Song struct {
	Name string
	Time string
}

type Settings struct {
	Output int
	Dev    bool
}

const sampleSize = 4 // from mp3 documentation

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
						Text: "generate",
						OnClicked: func() {
							AlbumInformation(window)
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

func AlbumInformation(owner walk.Form) (int, error) {
	var dlg *walk.MainWindow
	var txt *walk.TextEdit
	var uploadSuccess = walk.NewMutableCondition()
	// TODO: UNCOMMENT UF CALLS
	var uf []string
	album := new(Album)
	var link, ctn *walk.LineEdit

	MustRegisterCondition("uploadSuccess", uploadSuccess)

	return MainWindow{
		AssignTo: &dlg,
		Title:    "upload",
		Size:     Size{450, 280},
		MinSize:  Size{450, 280},
		Layout:   VBox{},
		OnDropFiles: func(files []string) {
			if CheckFiles(files) == "upload successful" {
				uf = files
				uploadSuccess.SetSatisfied(true)
			} else {
				uploadSuccess.SetSatisfied(false)
			}
			txt.SetText(CheckFiles(files))
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "bandcamp link:",
					},
					LineEdit{
						AssignTo: &link,
						OnTextChanged: func() {
							album.Link = link.Text()
						},
						OnEditingFinished: func() {
							fmt.Println("bc link set to: " + album.Link)
						},
					},
				},
			},
			Composite{
				Layout: Grid{
					Columns: 2,
				},
				Children: []Widget{
					Label{
						Text: "catalog no: GEO-",
					},
					LineEdit{
						AssignTo: &ctn,
						MaxSize:  Size{40, 0},
						OnTextChanged: func() {
							album.CatNo = ctn.Text()
						},
						OnEditingFinished: func() {
							fmt.Println("cat no set to: " + album.CatNo)
						},
					},
				},
			},
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					TextEdit{
						ColumnSpan: 4,
						AssignTo:   &txt,
						ReadOnly:   true,
						Text:       "drop album audio here\r\n(must be in format: mp3)",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text:    "submit",
						Enabled: uploadSuccess,
						OnClicked: func() {
							ScrapeBandcamp(album)
							fmt.Println("results-----")
							fmt.Println("name: " + album.Name)
							fmt.Println("artist: " + album.Artist)
							fmt.Println("year: " + album.Year)
							GenerateOutput(album, uf) // generate album description
						},
					},
				},
			},
		},
	}.Run()
}

// collect artist, album, and year fields from link
func ScrapeBandcamp(album *Album) {
	c := colly.NewCollector()
	// callbacks
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("#name-section", func(_ int, el *colly.HTMLElement) {
			album.Name = el.ChildText(".trackTitle") // album name
			album.Artist = el.ChildText("h3 span")   // artist name
		})

		e.ForEach(".tralbumData.tralbum-credits", func(_ int, el *colly.HTMLElement) {
			date, _, _ := strings.Cut(strings.TrimSpace(el.Text), "\n") // full date (released Jan 01, 20xx)
			album.Year = date[len(date)-4:]                             // grab year
		})
	})

	fmt.Println("Visiting: " + album.Link)
	c.Visit(album.Link) // visit, get results
}

// returns problem files if any extensions are not mp3. wav and flac unsupported for now.
// otherwise displays upload successful
func CheckFiles(files []string) string {
	var e []string
	for _, element := range files {
		fmt.Println("Processing upload... " + element)
		//if !strings.HasSuffix(element, ".mp3") && !strings.HasSuffix(element, ".wav") && !strings.HasSuffix(element, ".flac") {
		if !strings.HasSuffix(element, ".mp3") {
			fmt.Println(">> Illegal extension")
			e = append(e, element) // add problem file to slice
		}
	}
	if len(e) == 0 {
		return "upload successful"
	} else {
		return "error uploading\r\nformat not allowed:\r\n" + strings.Join(e, "\r\n")
	}
}

// generates GL video description output based on user input
// currently set to output in file 'output.txt'
// todo: allow user output settings
func GenerateOutput(al *Album, uf []string) {
	fmt.Println("Generating output...")
	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
	}
	w := io.MultiWriter(os.Stdout, f)
	defer f.Close() // close at finish
	fmt.Println("output.txt> File created ")

	// write header
	_, err = fmt.Fprintf(w,
		"Label: Geometric Lullaby\n"+
			"Artist: %s\n"+
			"Album: %s\n"+
			"Year: %s\n"+
			"(GEO - %s)\n",
		al.Artist, al.Name, al.Year, al.CatNo)

	_, err = fmt.Fprint(w,
		"\nDownload for free or purchase at: https://geometriclullaby.bandcamp.com/\n",
	)

	fmt.Println("output.txt> Writing track information")
	for _, element := range uf {
		f, err := os.Open(element) // element in this case is path to file
		if err != nil {
			fmt.Println(err)
		}

		d, err := mp3.NewDecoder(f)
		if err != nil {
			fmt.Println(err)
		}

		samples := d.Length() / sampleSize
		audioLength := samples / int64(d.SampleRate())
		//audiodur := time.Duration(audioLength time.Second)
		fmt.Fprintf(w, "\n"+element+" length: %d", audioLength)

		f.Close()
	}
	fmt.Println("\noutput.txt> Done.")
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

// OLD FUNCTIONS FOR CLI APPLICATION BELOW

// this function writes to the output.txt file and echoes
// its output to the cli
func write(ar string, al string, yr string, cat string, n []Song) {
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
