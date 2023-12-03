package restapi

import "fmt"

func couldNotRetrieveError(subject string) error {
	return fmt.Errorf("could not retrieve %s", subject)
}
