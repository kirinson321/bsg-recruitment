package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kirinson321/bsg-recruitment/pkg/domain"
)

type outputter struct {
}

func NewOutputter() domain.Outputter {
	return &outputter{}
}

var (
	LogFileName       = "log.txt"
	BackupLogFileName = "log.txt.old"
)

func (o *outputter) Output(data domain.StructuredOutput) error {
	// prepare the data for output
	out, err := prepareOutput(data)
	if err != nil {
		return fmt.Errorf("error preparing output: %w", err)
	}
	out += "\n"

	// output the data to stdout and to the log.txt file
	f, err := os.OpenFile(LogFileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening the log file for writing: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(out)
	if err != nil {
		return fmt.Errorf("error writing to the log file: %w", err)
	}

	fmt.Println(out)

	return nil
}

func prepareOutput(data domain.StructuredOutput) (string, error) {
	// prepare the data for output
	out, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling data: %w", err)
	}
	return string(out), nil
}
