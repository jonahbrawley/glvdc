package main

import(
	"fmt"
	"bufio"
	"os"
	"strings"
)

func main() {
	fmt.Println("-------- glvdc --------")
	fmt.Println("-------- by: me -------")
	fmt.Println("")
	url := prompt("bc url:")
	fmt.Println(url)
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