package output

import (
	"encoding/json"
	"fmt"

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

	// output the data to stdout and to the log.txt file

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
