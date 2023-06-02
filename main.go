package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

var args struct {
	Input  string `arg:"-i,--input" help:"Input file path"`
	Output string `arg:"-o,--output" help:"Output file path"`
}

func main() {
	arg.MustParse(&args)
	if args.Input == "" {
		fmt.Println("Input file path is required")
		return
	}
	file, err := os.Open(args.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Println("Transforming " + args.Input)
	transformed := transform(file)
	var outputPath string
	if args.Output == "" {
		outputPath = fmt.Sprintf("%s.export", args.Input)
	} else {
		outputPath = args.Output
	}
	writeToFile(transformed, outputPath)
	fmt.Println("Exporeted to " + outputPath)
}

func transform(file *os.File) string {
	scanner := bufio.NewScanner(file)
	var lines []string
	var interfaceName string
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "class") {
			interfaceName = strings.TrimSpace(strings.Split(strings.Split(scanner.Text(), "class")[1], "{")[0])
			lines = append(lines, "export interface "+interfaceName+" {")
		}
		if !strings.Contains(scanner.Text(), "@") && !strings.Contains(scanner.Text(), "\"") {
			if strings.Contains(scanner.Text(), ":") {
				lines = append(lines, scanner.Text())
			}
		}
	}
	lines = append(lines, "}")
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return strings.Join(lines, "\n")
}

func writeToFile(str string, outputPath string) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	outputFile.WriteString(str)
}
