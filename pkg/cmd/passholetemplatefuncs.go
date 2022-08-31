package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/twpayne/chezmoi/v2/pkg/chezmoilog"
)

type passholeCacheKey struct {
	path  string
	field string
}

type passholeConfig struct {
	Command  string
	Args     []string
	Prompt   bool
	cache    map[passholeCacheKey]string
	password string
}

func (c *Config) passholeTemplateFunc(path, field string) string {
	key := passholeCacheKey{
		path:  path,
		field: field,
	}
	if value, ok := c.Passhole.cache[key]; ok {
		return value
	}

	args := append([]string{}, c.Passhole.Args...)
	args = append(args, "show", "--field", field)
	if !c.Passhole.Prompt {
		args = append(args, "--no-password")
	}
	args = append(args, path)
	output, err := c.passholeOutput(c.Passhole.Command, args)
	if err != nil {
		panic(err)
	}

	if c.Passhole.cache == nil {
		c.Passhole.cache = make(map[passholeCacheKey]string)
	}
	c.Passhole.cache[key] = output
	return output
}

func (c *Config) passholeOutput(name string, args []string) (string, error) {
	if c.Passhole.Prompt && c.Passhole.password == "" {
		password, err := c.readPassword("Enter database password: ")
		if err != nil {
			return "", err
		}
		c.Passhole.password = password
	}

	cmd := exec.Command(name, args...)
	// FIXME the following does not work because passhole always uses zenity to
	// read the password if its stdin is not a tty. See:
	// https://github.com/Evidlo/passhole/blob/v1.9.9/passhole/passhole.py#L370-L388
	cmd.Stdin = bytes.NewBufferString(c.Passhole.password + "\n")
	cmd.Stderr = os.Stderr
	output, err := chezmoilog.LogCmdOutput(cmd)
	if err != nil {
		return "", newCmdOutputError(cmd, output, err)
	}
	return string(output), nil
}
