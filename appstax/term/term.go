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
)

const indent = "   "

var section = false

func Section() {
	if !section {
		section = true
		fmt.Println("")
	}
}

func Println(text string) {
	fmt.Println(indent + text)
	section = false
}

func Print(text string) {
	fmt.Print(indent + text)
	section = false
}

func Printf(format string, a ...interface{}) {
	fmt.Printf(indent+format, a...)
	section = false
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
