package postgres

import "fmt"

func sqlBuilderError(subject string) error {
	return fmt.Errorf("could not build SQL for %s", subject)
}

func couldNotRetrieveError(subject string) error {
	return fmt.Errorf("could not retrieve %s", subject)
}
