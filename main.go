package main

import (
	"fmt"
	"julianiff/vamoni-git/utils"
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

	fmt.Println("Initialized the vamoni repository")
}

func (r *Repository) detectChangedFiles() ([]string, error) {
	output := editedFiles(filepath.Join(r.basePath, changeDir))
	return output, nil
}

func (r *Repository) Status() error {
	changedFiles, err := r.detectChangedFiles()
	if err != nil {
		return fmt.Errorf("failed to detect files: %w", err)
	}

	allStagedFiles := stagedFiles()
	fmt.Printf("Changed files: %v\n", changedFiles)
	fmt.Printf("Staged files: %v\n", allStagedFiles)

	return nil
}

func (r *Repository) CommitStagedFiles() error {

	stagePath := filepath.Join(r.basePath, stageDir)
	stageFiles, err := os.ReadDir(stagePath)
	if err != nil {
		return fmt.Errorf("filed to read stagedFiles")
	}
	var allStagedFiles []string
	for _, s := range stageFiles {
		allStagedFiles = append(allStagedFiles, s.Name())
	}

	changedPath := filepath.Join(r.basePath, changeDir)
	for _, f := range allStagedFiles {
		err := utils.Copyfile(f, changedPath+f)
		if err != nil {
			return fmt.Errorf("error while copying")
		}
		os.Remove(stagePath + f)
	}

	fmt.Println("New files added to changeset", allStagedFiles)

	return nil
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

func stagedFiles() []string {
	stageFiles, err := os.ReadDir(".vamoni/stage")
	check(err)
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
		utils.Copyfile(fileToStage, ".vamoni/stage/"+fileToStage)
	}

	newlyStagedFiles := stagedFiles()
	fmt.Println("Currently Staged Files", newlyStagedFiles)
	fmt.Println("files that can be staged", allEditedFiles)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}

	repo := newRepository()

	switch os.Args[1] {

	case "init":
		repo.init()
	case "status":
		repo.Status()
	case "detect":
		repo.detectChangedFiles()
	case "stage":
		args2 := os.Args[2]
		stage(args2)
	case "commit":
		repo.CommitStagedFiles()
	default:
		fmt.Println("No valid method names provided")
		os.Exit(1)
	}
}
