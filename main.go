package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

func main() {

	inputFile, err := os.Open("files.txt")
	if err != nil {
		fmt.Printf("Failed to open files.txt: %v\n", err)
		return
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)

	var wg sync.WaitGroup

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("Failed to read line in files.txt: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		paths := strings.Split(line, " ")

		wg.Add(1)
		go func() {
			defer wg.Done()
			format(paths[0], paths[1])
		}()
		if err == io.EOF {
			break
		}
	}

	wg.Wait()
}

func format(inputPath, outputPath string) {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Failed to open %q: %v\n", inputPath, err)
		return
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Failed to create %q: %v\n", outputPath, err)
		return
	}
	defer outputFile.Close()

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	inputString, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Printf("Failed to read input file %q: %v\n", inputPath, err)
		return
	}

	inputString = strings.TrimSpace(inputString)

	inputString = fixPunctuation(inputString)

	inputString = transfer(inputString)

	_, err = writer.WriteString(inputString)
	if err != nil {
		fmt.Printf("Failed to write output to %q: %v\n", outputPath, err)
		return
	}
}

// Регулярные выражения
//
//	 "+" - один или более
//	 "*" - ноль или более
//	"\s" пробел
//	"$1" найденый знак
//
// ([,:;—-]) массив подходящих символов
func fixPunctuation(input string) string {

	re := regexp.MustCompile(`\s+([.,:;—-])`)
	input = re.ReplaceAllString(input, "$1")

	re = regexp.MustCompile(`([,:;])\s*`)
	input = re.ReplaceAllString(input, "$1 ")

	re = regexp.MustCompile(`\s*—\s*`)
	input = re.ReplaceAllString(input, " — ")

	re = regexp.MustCompile(`\s*-\s*`)
	input = re.ReplaceAllString(input, "-")

	re = regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")

	return input
}

func transfer(input string) string {

	words := strings.Fields(input)

	lineLength := longestWord(input)
	var result strings.Builder

	currentLine := ""
	currentLength := 0

	for i, word := range words {
		wordLength := len(word)
		if currentLength > 0 {
			wordLength += 1
		}
		if currentLength+wordLength <= lineLength {
			if currentLength > 0 {
				currentLine += " "
			}
			currentLine += word
			currentLength += wordLength
		} else {
			result.WriteString(currentLine + "\n")
			currentLine = word
			currentLength = len(word)
		}
		if i == len(words)-1 {
			result.WriteString(currentLine)
		}
	}

	return result.String()

}

func longestWord(s string) int {

	words := strings.FieldsFunc(s, func(c rune) bool {
		//	НЕ буквы и НЕ цифры как разделители на слова
		return !unicode.IsLetter(c) && !unicode.IsDigit(c)
	})

	maxLength := 0

	for _, word := range words {
		if len(word) > maxLength {
			maxLength = len(word)
		}
	}

	return maxLength * 3
}
