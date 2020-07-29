package parse

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type SpaceSeparatedListFlag struct {
	name   string
	values []string
}

func NewSpaceSeparatedListFlag(name string) *SpaceSeparatedListFlag {
	return &SpaceSeparatedListFlag{name: name, values: []string{}}
}

// Implement the flag.Value interface
func (s *SpaceSeparatedListFlag) String() string {
	return ""
}

func (s *SpaceSeparatedListFlag) Type() string {
	return "stringList"
}

// empty implementation
func (s *SpaceSeparatedListFlag) Set(_ string) error {
	return nil
}

// use this manually to actually initialize
// split string to a map according to special rules due to task params resolution limitations
func (s *SpaceSeparatedListFlag) SetReal() error {
	if s == nil || s.name == "" {
		return errors.New("must be initialized with a name")
	}

	fullFlagName := "--" + s.name
	fullFlagNameWithEquals := fullFlagName + "="

	// locate the next occurrence of this flag
	foundFlag := false
	saveValues := false

	// use all args instead since flag/pflag does not support this feature
	for _, argRaw := range os.Args[1:] {
		arg := strings.TrimSpace(argRaw)
		if arg == "" {
			continue
		}
		switch {
		case arg == fullFlagName:
			if foundFlag {
				return fmt.Errorf("duplicates of %v option are not allowed", fullFlagName)
			}
			foundFlag = true
			saveValues = true
		case strings.HasPrefix(arg, fullFlagNameWithEquals):
			foundFlag = true
			saveValues = true

			if val := arg[len(fullFlagNameWithEquals):]; val != "" {
				// save value behind equals
				s.values = append(s.values, val)
			}
		case arg[0] == '-':
			if foundFlag {
				// new flag found - stop processing this one
				saveValues = false
			}
		default:
			if saveValues {
				// save all consequent values of our flag
				s.values = append(s.values, arg)
			}
		}
	}
	return nil
}

// Example: --param NAME a DESC "description with spaces"
//   becomes ["NAME", "a", "DESC", "description with spaces"]
func (s *SpaceSeparatedListFlag) GetValues() []string {
	return s.values
}

// split string to a map according to special rules due to task params resolution limitations
//
// Example: --param NAME a DESC "description with spaces"
//   becomes { "NAME": "a", "DESC": "description with spaces"}
func (s *SpaceSeparatedListFlag) GetMapValues() map[string]string {
	result := make(map[string]string)

	key := ""
	for idx, val := range s.values {
		if idx%2 == 0 {
			key = val
		} else if key != "" {
			result[key] = val
		}

	}
	return result
}
