package term

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"appstax-cli/appstax/log"
	"github.com/cheggaaa/pb"
	"github.com/howeyc/gopass"
	"github.com/keronsen/tablewriter"
)

var layout = true
var indent = "   "
var section = false

func Section() {
	section = true
}

func PrintSection() {
	section = true
	printSection()
}

func Layout(flag bool) {
	layout = flag
}

func printSection() {
	if layout && section {
		fmt.Println("")
		section = false
	}
}

func printIndent() {
	if layout {
		fmt.Print(indent)
	}
}

func Dump(value interface{}) {
	fmt.Printf("%v", value)
}

func Println(text string) {
	printSection()
	printIndent()
	fmt.Println(text)
}

func Print(text string) {
	printSection()
	printIndent()
	fmt.Print(text)
}

func Printf(format string, a ...interface{}) {
	printSection()
	printIndent()
	fmt.Printf(format, a...)
}

func PrintTable(headers []string, rows [][]string) {
	printSection()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}

func GetString(prompt string) string {
	Print(prompt + ": ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func GetInt(prompt string) int {
	for {
		str := GetString(prompt)
		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			Println("Not a number, try again.")
		} else {
			return int(i)
		}
	}
}

func GetPassword(prompt string) string {
	Print(prompt + ": ")
	password := gopass.GetPasswdMasked()
	return strings.TrimSpace(string(password))
}

func ShowProgressBar(totalBytes int64) *pb.ProgressBar {
	log.Debugf("Creating progress bar for %d bytes", totalBytes)
	bar := pb.New64(totalBytes)
	bar.SetMaxWidth(70)
	bar.ShowCounters = false
	bar.ShowSpeed = true
	bar.SetUnits(pb.U_BYTES)
	bar.Format("[##-]")
	bar.Prefix("   ")
	bar.Start()
	return bar
}
