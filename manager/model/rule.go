package model

import (
	"regexp"
)

// Match is the requirement for a rule to take action.
type Match struct {
	Filename     string `json:"filename"`      // The name of the uploaded file.
	SaveLocation string `json:"save_location"` // The location to save the file.
}

// Rule is a condition that will be checked to run determined scripts.
type Rule struct {
	Match    Match    `json:"match"`    // Requirements to apply the rule.
	Pipeline []string `json:"pipeline"` // Pipeline of script execution.
}

// Validate performs the validation of the rule, all criterias must be true.
func (m Match) Validate(Filename, SaveLocation string) (matches bool, err error) {
	if match, err := regexp.MatchString(m.Filename, Filename); err != nil || !match {
		return false, err
	}

	if match, err := regexp.MatchString(m.SaveLocation, SaveLocation); err != nil || !match {
		return false, err
	}

	return true, nil
}
