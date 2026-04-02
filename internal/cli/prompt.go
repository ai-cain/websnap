package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type linePrompter struct {
	reader *bufio.Reader
	writer io.Writer
}

func newLinePrompter(input io.Reader, output io.Writer) linePrompter {
	if input == nil {
		input = strings.NewReader("")
	}

	if output == nil {
		output = io.Discard
	}

	return linePrompter{
		reader: bufio.NewReader(input),
		writer: output,
	}
}

func (p linePrompter) Ask(prompt string) (string, error) {
	if _, err := fmt.Fprint(p.writer, prompt); err != nil {
		return "", err
	}

	value, err := p.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.TrimRight(value, "\r\n"), nil
}
