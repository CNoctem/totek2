package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetConfig(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	properties := make(map[string]string)

	// Start reading from the file using a scanner.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.Index(line, "=")
		if idx != -1 {
			properties[line[:idx]] = line[idx+1:]
		}
	}
	if scanner.Err() != nil {
		fmt.Printf(" > Failed with error %v\n", scanner.Err())
		return nil, scanner.Err()
	}
	return properties, nil
}
