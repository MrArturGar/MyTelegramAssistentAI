package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetValue(key string) string {
	value := ""
	readFile, err := os.Open("setting.env")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), key) {
			value = strings.Replace(fileScanner.Text(), key+"=", "", 1)
		}
	}

	readFile.Close()
	return value
}
