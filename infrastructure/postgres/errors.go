package postgres

import "fmt"

func sqlBuilderError(subject string, err error) error {
	return fmt.Errorf("could not build SQL for %s: %v", subject, err)
}

func couldNotRetrieveError(subject string, err error) error {
	return fmt.Errorf("could not retrieve %s: %v", subject, err)
}

func couldNotCreateError(subject string, err error) error {
	return fmt.Errorf("could not create %s: %v", subject, err)
}

func couldNotUpdateError(subject string, err error) error {
	return fmt.Errorf("could not update %s: %v", subject, err)
}

func couldNotDeleteError(subject string, err error) error {
	return fmt.Errorf("could not delete %s: %v", subject, err)
}
