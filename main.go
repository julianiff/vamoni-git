package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	vamoniDir  = ".vamoni"
	changeDir  = "change"
	stageDir   = "stage"
	permission = 0755
)

type Repository struct {
	basePath string
}

func newRepository() *Repository {
	return &Repository{
		basePath: vamoniDir,
	}
}

func (r *Repository) init() {
	dirs := []string{
		filepath.Join(r.basePath, changeDir),
		filepath.Join(r.basePath, stageDir),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("vamoni folder created")
}

func vamoniInit() {
	if err := os.MkdirAll(".vamoni/change", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(".vamoni/stage", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	fmt.Println("vamoni folder created")
}

func status() {
	detect()
	allStagedFiles := stagedFiles()
	fmt.Println("these files are staged", allStagedFiles)
}

func copyfile(sourceFile string, destinationFile string) {

	source, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer source.Close()

	destination, err := os.Create(destinationFile)
	if err != nil {
		panic(err)
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		panic(err)
	}
}

func diff(workingDir []string, storedDir []string) []string {
	var difference []string
	for _, w1 := range workingDir {
		if !slices.Contains(storedDir, w1) {
			difference = append(difference, w1)
		}
	}

	return difference
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func editedFiles(path string) []string {

	filesStored, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var filesStoredNames []string
	for _, s := range filesStored {
		filesStoredNames = append(filesStoredNames, s.Name())
	}

	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	var changedFiles []string
	for _, s := range files {
		if !strings.HasPrefix(s.Name(), ".") {
			changedFiles = append(changedFiles, s.Name())
		}
	}

	return diff(changedFiles, filesStoredNames)
}

func detect() {
	output := editedFiles(".vamoni/change")
	fmt.Println("these files are changed", output)
}

func stagedFiles() []string {
	stageFiles, err := os.ReadDir(".vamoni/stage")
	if err != nil {
		log.Fatal(err)
	}
	var allStagedFiles []string
	for _, s := range stageFiles {
		allStagedFiles = append(allStagedFiles, s.Name())
	}
	return allStagedFiles
}

func stage(fileToStage string) {
	allEditedFiles := editedFiles(".vamoni/change")
	allStagedFiles := stagedFiles()

	// only files that are edited can be staged
	if slices.Contains(allEditedFiles, fileToStage) && !slices.Contains(allStagedFiles, fileToStage) {
		fmt.Println("file can be staged", fileToStage)
		copyfile(fileToStage, ".vamoni/stage/"+fileToStage)
	}

	newlyStagedFiles := stagedFiles()
	fmt.Println("Currently Staged Files", newlyStagedFiles)
	fmt.Println("files that can be staged", allEditedFiles)
}

func change() {
	stageFiles, err := os.ReadDir(".vamoni/stage")
	if err != nil {
		log.Fatal(err)
	}
	var allStagedFiles []string
	for _, s := range stageFiles {
		allStagedFiles = append(allStagedFiles, s.Name())
	}

	for _, f := range allStagedFiles {
		copyfile(f, ".vamoni/change/"+f)
		os.Remove(".vamoni/stage/" + f)
	}

	fmt.Println("New files added to changeset", allStagedFiles)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}

	repo := newRepository()

	switch os.Args[1] {

	case "init":
		vamoniInit()
	case "status":
		status()
	case "detect":
		detect()
	case "stage":
		args2 := os.Args[2]
		stage(args2)
	case "change":
		change()
	default:
		fmt.Println("No valid method names provided")
		os.Exit(1)
	}
}
