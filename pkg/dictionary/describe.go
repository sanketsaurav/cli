package dictionary

import (
	"fmt"
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/go-fastly/fastly"
)

// DescribeCommand calls the Fastly API to describe a dictionary.
type DescribeCommand struct {
	common.Base
	manifest manifest.Data

	// required
	DictionaryName string // Can't shaddow common.Base method Name().
	Version        int
}

// NewDescribeCommand returns a usable command registered under the parent.
func NewDescribeCommand(parent common.Registerer, globals *config.Data) *DescribeCommand {
	var c DescribeCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("describe", "Show detailed information about a dictionary on a Fastly service version").Alias("get")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("version", "Number of service version").Required().IntVar(&c.Version)
	c.CmdClause.Flag("name", "Name of the dictionary").Short('n').Required().StringVar(&c.DictionaryName)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *DescribeCommand) createInput() (*fastly.GetDictionaryInput, error) {
	var input fastly.GetDictionaryInput

	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input.Service = serviceID
	input.Version = c.Version
	input.Name = c.DictionaryName

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *DescribeCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	d, err := c.Globals.Client.GetDictionary(input)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Service ID: %s\n", d.ServiceID)
	fmt.Fprintf(out, "Version: %d\n", d.Version)
	fmt.Fprintf(out, "ID: %s\n", d.ID)
	fmt.Fprintf(out, "Name: %s\n", d.Name)
	fmt.Fprintf(out, "Write only: %v\n", d.WriteOnly)

	if d.CreatedAt != nil {
		fmt.Fprintf(out, "Created (UTC): %s\n", d.CreatedAt.UTC().Format(common.TimeFormat))
	}
	if d.UpdatedAt != nil {
		fmt.Fprintf(out, "Last edited (UTC): %s\n", d.UpdatedAt.UTC().Format(common.TimeFormat))
	}
	if d.DeletedAt != nil {
		fmt.Fprintf(out, "Deleted (UTC): %s\n", d.DeletedAt.UTC().Format(common.TimeFormat))
	}

	return nil
}
