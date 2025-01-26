package main

import (
	"fmt"
	repository "julianiff/vamoni-git/internal"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}

	repo := repository.NewRepository()

	switch os.Args[1] {

	case "init":
		repo.Init()
	case "status":
		repo.Status()
	case "detect":
		repo.DetectChangedFiles()
	case "stage":
		args2 := os.Args[2]
		repo.Stage(args2)
	case "commit":
		repo.CommitStagedFiles()
	default:
		fmt.Println("No valid method names provided")
		os.Exit(1)
	}
}
