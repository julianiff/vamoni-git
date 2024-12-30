package main

import (
	"fmt"
	"log"
	"os"
)

// vamoni init -- creates a .vamoni folder that holds the local config

// vamoni detect -- detects changes that were made to the folder and all subfolders
// --> how to detect that a change is done?
// file old (saved in .vamoni?) file new. Compare hash, if the same, break
// if hash not the same, commpare line by line, if the same break
// if line are not same, print new line

// vamoni change "creates a new change of the repository"
// --> cp all files into .vamoni folder (we can later compress)
// change -> a change is the delta of files before and files after
// 0 files -> 3 files -- change is the 3 new files

func vamoniInit() {
	if err := os.MkdirAll(".vamoni/change", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	fmt.Println("vamoni folder created")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "init":
		vamoniInit()
	case "detect":
		fmt.Println("detect")
	case "change":
		fmt.Println("change")
	default:
		fmt.Println("No valid method names provided")
		os.Exit(1)
	}
}
