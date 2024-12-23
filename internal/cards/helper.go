package cards

import (
	"fmt"
	"regexp"
)

// validateEmail checks if the given email has a valid format.
func validateEmail(email string) error {
	if len(email) == 0 {
		return nil // Allow empty email, assuming it's optional
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return fmt.Errorf("error validating email: %v", err)
	}
	if !match {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
