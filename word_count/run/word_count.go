package run

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type AppEnv struct {
	ByteCount bool
	WordCount bool
	LineCount bool
	CharCount bool
	Filename  string
}

type countResults struct {
	byteCount int
	wordCount int
	lineCount int
	charCount int
}

func (app *AppEnv) Run() error {
	counts, err := app.runCounts()
	if err != nil {
		return err
	}
	var b strings.Builder
	if app.ByteCount {
		fmt.Fprintf(&b, "%d ", counts.byteCount)
	}
	if app.CharCount {
		fmt.Fprintf(&b, "%d ", counts.charCount)
	}
	if app.LineCount {
		fmt.Fprintf(&b, "%d ", counts.lineCount)
	}
	if app.WordCount {
		fmt.Fprintf(&b, "%d ", counts.wordCount)
	}
	b.WriteString(app.Filename)
	fmt.Fprint(os.Stdout, b.String())
	return nil
}

func (app *AppEnv) runCounts() (countResults, error) {
	var word_count, byte_count, line_count, char_count int
	var counts countResults
	var err error
	if app.WordCount {
		word_count, err = getCount(app.Filename, bufio.ScanWords)
		if err != nil {
			return counts, err
		}
		counts.wordCount = word_count
	}
	if app.ByteCount {
		byte_count, err = getCount(app.Filename, bufio.ScanBytes)
		if err != nil {
			return counts, err
		}
		counts.byteCount = byte_count
	}
	if app.LineCount {
		line_count, err = getCount(app.Filename, bufio.ScanLines)
		if err != nil {
			return counts, err
		}
		counts.lineCount = line_count
	}
	if app.CharCount {
		char_count, err = getCount(app.Filename, bufio.ScanRunes)
		if err != nil {
			return counts, err
		}
		counts.charCount = char_count
	}

	return counts, nil
}

func getCount(filename string, split_fun bufio.SplitFunc) (int, error) {
	data, err := os.Open(filename)
	defer data.Close()
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(data)
	scanner.Split(split_fun)
	count := 0
	for scanner.Scan() {
		count += 1
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
