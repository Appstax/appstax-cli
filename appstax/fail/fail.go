package fail

import (
	"appstax-cli/appstax/log"
	"os"
)

var panicMode = false

func SetPanicMode(m bool) {
	panicMode = m
	if panicMode {
		println("panic on")
	}
}

func Handle(err error) {
	if err != nil {
		if panicMode {
			panic(err)
		} else {
			log.Panicf(err.Error())
			println("\nSomething went wrong! See appstax.log for details.")
			os.Exit(-1)
		}
	}
}
