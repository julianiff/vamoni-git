package main

import (
	"fmt"
	repository "julianiff/vamoni-git/internal"
	"log"
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
		err := repo.Status()
		if err != nil {
			fmt.Printf("could not call status %v", err)
		}
	case "detect":
		_, err := repo.DetectChangedFiles()
		if err != nil {
			fmt.Printf("Could not detect %v", err)
		}
	case "stage":
		args2 := os.Args[2]
		repo.Stage(args2)
	case "commit":
		if err := repo.CommitStagedFiles(); err != nil {
			log.Fatalf("Failed to commit")
		}
	default:
		fmt.Println("No valid method names provided")
		os.Exit(1)
	}
}
