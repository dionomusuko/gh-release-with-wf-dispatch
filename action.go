package main

import (
	"fmt"
	"os"
)

func setOutput(name, value string) error {
	outputFilePath := os.Getenv("GITHUB_OUTPUT")
	env := fmt.Sprintf("%s=%s\n", name, value)

	file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("failed to close")
		}
	}(file)

	_, err = file.WriteString(env)
	if err != nil {
		fmt.Println("failed to write env to file")
		return err
	}

	return nil
}
