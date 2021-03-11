package options

import (
	"fmt"
	"github.com/google/shlex"
	"strings"
)

type CommandOptions struct {
	opts []string
}

func NewCommandOptions(options string) (*CommandOptions, error) {
	opts, err := shlex.Split(options)
	if err != nil {
		return nil, err
	}
	return &CommandOptions{opts: opts}, err
}

func NewCommandOptionsFromArray(options []string) *CommandOptions {
	return &CommandOptions{opts: options}
}

func (c *CommandOptions) GetAll() []string {
	return c.opts
}

func (c *CommandOptions) GetOptionValue(option string) string {
	idx := c.getOptionIndex(option)

	if idx < 0 {
		return ""
	}

	optionKey := c.opts[idx]

	if isShortOption(option) && len(optionKey) > 2 && optionKey[2] != '=' {
		return optionKey[2:]
	}

	separatedOption := strings.SplitN(optionKey, "=", 2)

	if len(separatedOption) == 2 {
		return separatedOption[1]
	}

	if idx+1 < len(c.opts) {
		value := c.opts[idx+1]

		if value[0] != '-' {
			return value
		}
	}

	return ""
}

func (c *CommandOptions) IncludesOption(option string) bool {
	return c.getOptionIndex(option) >= 0
}

func (c *CommandOptions) IncludesString(substr string) bool {
	return strings.Contains(strings.Join(c.opts, " "), substr)
}

func (c *CommandOptions) AddOption(name, value string) {
	c.opts = append(c.opts, name, value)
}

func (c *CommandOptions) AddOptions(values ...string) {
	c.opts = append(c.opts, values...)
}

func (c *CommandOptions) AddFlag(flag string) {
	c.opts = append(c.opts, flag)
}

func (c *CommandOptions) AddValue(value string) {
	c.opts = append(c.opts, value)
}

func (c *CommandOptions) ToString() string {
	if c == nil {
		return "nil"
	}

	return fmt.Sprintf("[%v]", strings.Join(c.opts, ", "))
}

func (c *CommandOptions) getOptionIndex(option string) int {
	if len(option) < 2 || option[0] != '-' {
		return -1
	}

	isShort := isShortOption(option)

	for idx, opt := range c.opts {
		if isShort {
			if strings.HasPrefix(opt, option) {
				return idx
			}
		} else {
			if opt == option || strings.HasPrefix(opt, option+"=") {
				return idx
			}
		}
	}
	return -1
}

func isShortOption(option string) bool {
	return len(option) == 2 && option[0] == '-' && option[1] != '-'
}
