package main

import (
	"fmt"
	"os"
)


func log(msg string, args ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintf(msg, args...) + "\n")
}
