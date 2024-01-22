package sqlite

import "fmt"

func couldNotRetrieveAllError(subject string, expectedNum, actualNum int) error {
	return fmt.Errorf("could not retrieve all %s; expected %d, actual %d", subject, expectedNum, actualNum)
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
