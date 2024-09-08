package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/joho/godotenv"
)

var rollMethod = "random"
var rollBigString string
var rollStrings []string
var rollStringsDelimiter string
var checkFile string
var writePath string
var activationString string
var fixedSuffix string

func handleErrors(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func loadEnv() {
	godotenv.Load()

	rollMethod = os.Getenv("roll_method")
	rollBigString = os.Getenv("roll_strings")
	rollStringsDelimiter = os.Getenv("string_delimiter")
	checkFile = os.Getenv("check_file_path")
	writePath = os.Getenv("write_to_path")
	activationString = os.Getenv("activation_string")
	fixedSuffix = os.Getenv("fixed_suffix")

	fillStringsSlice()
}

func fillStringsSlice() {
	//To be honest, I don't even know where to start on this
	//If I format it like a slice, does it just exist as slice? Surely not, right?
	//If it doesn't, how do we delimit the strings?? I hope commas, backslashes don't work for my use case

	rollStrings = regexp.MustCompile(rollStringsDelimiter).Split(rollBigString, -1)
}

func checkForRotate() {
	t, err := tail.TailFile(checkFile, tail.Config{Follow: true})
	for line := range t.Lines { //Annoying side-effect, it reads all the old ones when you first start it so we might need to flush this file regularly
		//fmt.Println(line.Text)
		if strings.TrimSpace(line.Text) == activationString {
			fmt.Println("Match found, rotate.")
			rotateStrings()
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}

func rotateStrings() {
	currString := rollStrings[rand.Intn(len(rollStrings))]

	writeString := currString + fixedSuffix

	sb := []byte(writeString)
	fmt.Println(writeString)
	errWF := os.WriteFile(writePath, sb, 0644)
	handleErrors(errWF)
}

func main() {
	loadEnv()
	checkForRotate()
}
