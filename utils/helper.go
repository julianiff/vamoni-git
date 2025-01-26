package utils

import (
	"io"
	"os"
	"slices"
)

func Diff(workingDir []string, storedDir []string) []string {
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

func GetFilesInPath(path string) ([]string, error) {
	stageFiles, err := os.ReadDir(path)
	check(err)
	var allStagedFiles []string
	for _, s := range stageFiles {
		allStagedFiles = append(allStagedFiles, s.Name())
	}

	return allStagedFiles, nil
}

func Copyfile(sourceFile string, destinationFile string) error {

	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destinationFile)
	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
