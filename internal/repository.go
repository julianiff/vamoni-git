package repository

import (
	"encoding/base64"
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

	if err := os.WriteFile(filepath.Join(r.basePath, "index"), []byte("000000000"+"\n"), permission); err != nil {
		log.Fatal((err))
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
	fmt.Printf("Files with changes: %v\n", changedFiles)
	fmt.Printf("Staged files: %v\n", allStagedFiles)

	return nil
}

func encodeFilePathsContent(filePaths []string) (string, error) {
	var fileBytes []byte
	for _, s := range filePaths {
		bytes, _ := os.ReadFile(s)
		dst := base64.StdEncoding.EncodeToString(bytes)
		fileBytes = append(fileBytes, dst...)
	}

	combined := make([]byte, base64.StdEncoding.EncodedLen(len(fileBytes)))
	base64.StdEncoding.Encode(combined, fileBytes)

	return string(combined), nil
}

func (r *Repository) CommitStagedFiles(message string) error {

	stagePath := filepath.Join(r.basePath, stageDir)
	stageFiles, err := os.ReadDir(stagePath)
	if err != nil {
		return fmt.Errorf("filed to read stagedFiles")
	}

	var allStagedFiles []string
	for _, s := range stageFiles {
		allStagedFiles = append(allStagedFiles, s.Name())
	}

	// here we have the files that we want to commit
	newCommitHash, _ := encodeFilePathsContent(allStagedFiles)
	file, err := os.OpenFile(filepath.Join(r.basePath, "index"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open example.txt")
		return err
	}

	defer file.Close()

	if _, err2 := file.WriteString(newCommitHash[:24] + " parenthash " + message + "\n"); err2 != nil {
		fmt.Println(err2)
		return err2
	}

	changedPath := filepath.Join(r.basePath, changeDir)
	for _, f := range allStagedFiles {
		err := utils.Copyfile(f, changedPath+"/"+newCommitHash[:24]+f)
		if err != nil {
			return fmt.Errorf("error while copying")
		}
		os.Remove(stagePath + f)
	}

	fmt.Println("Files commited", allStagedFiles)

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
