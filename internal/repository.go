package repository

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

func NewRepository() *Repository {
	return &Repository{
		basePath: vamoniDir,
	}
}

func (r *Repository) Init() {
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

func (r *Repository) DetectChangedFiles() ([]string, error) {
	output := r.EditedFiles(filepath.Join(r.basePath, changeDir))
	return output, nil
}

func (r *Repository) Status() error {
	changedFiles, err := r.DetectChangedFiles()
	if err != nil {
		return fmt.Errorf("failed to detect files: %w", err)
	}

	stagedFilePath := filepath.Join(r.basePath, stageDir)
	allStagedFiles, err := utils.GetFilesInPath(stagedFilePath)
	if err != nil {
		return fmt.Errorf("failed to get staged files")
	}
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

func (r *Repository) EditedFiles(path string) []string {

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

	return utils.Diff(changedFiles, filesStoredNames)
}

func (r *Repository) Stage(fileToStage string) {
	allEditedFiles := r.EditedFiles(".vamoni/change")
	stagedFilePath := filepath.Join(r.basePath, stageDir)
	allStagedFiles, err := utils.GetFilesInPath(stagedFilePath)

	if err != nil {
		fmt.Println("err")
	}

	// only files that are edited can be staged
	if slices.Contains(allEditedFiles, fileToStage) && !slices.Contains(allStagedFiles, fileToStage) {
		fmt.Println("file can be staged", fileToStage)
		utils.Copyfile(fileToStage, ".vamoni/stage/"+fileToStage)
	}

	newlyStagedFiles, err := utils.GetFilesInPath(stagedFilePath)
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println("Currently Staged Files", newlyStagedFiles)
	fmt.Println("files that can be staged", allEditedFiles)
}
