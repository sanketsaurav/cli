package dictionary

import (
	"fmt"
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/fastly"
)

// ListCommand calls the Fastly API to list domains.
type ListCommand struct {
	common.Base
	manifest manifest.Data

	// required
	Version int
}

// NewListCommand returns a usable command registered under the parent.
func NewListCommand(parent common.Registerer, globals *config.Data) *ListCommand {
	var c ListCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("list", "List dictionaries on a Fastly service version")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("version", "Number of service version").Required().IntVar(&c.Version)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *ListCommand) createInput() (*fastly.ListDictionariesInput, error) {
	var input fastly.ListDictionariesInput

	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input.Service = serviceID
	input.Version = c.Version

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *ListCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	dictionaries, err := c.Globals.Client.ListDictionaries(input)
	if err != nil {
		return err
	}

	if !c.Globals.Verbose() {
		tw := text.NewTable(out)
		tw.AddHeader("SERVICE", "VERSION", "ID", "NAME", "WRITE ONLY")
		for _, d := range dictionaries {
			tw.AddLine(d.ServiceID, d.Version, d.ID, d.Name, d.WriteOnly)
		}
		tw.Print()
		return nil
	}

	fmt.Fprintf(out, "Service ID: %s\n", input.Service)
	fmt.Fprintf(out, "Version: %d\n", input.Version)
	for i, d := range dictionaries {
		fmt.Fprintf(out, "\tDictionary %d/%d\n", i+1, len(dictionaries))
		fmt.Fprintf(out, "\t\tID: %s\n", d.ID)
		fmt.Fprintf(out, "\t\tName: %s\n", d.Name)
		fmt.Fprintf(out, "\t\tWrite only: %v\n", d.WriteOnly)

		if d.CreatedAt != nil {
			fmt.Fprintf(out, "\t\tCreated (UTC): %s\n", d.CreatedAt.UTC().Format(common.TimeFormat))
		}
		if d.UpdatedAt != nil {
			fmt.Fprintf(out, "\t\tLast edited (UTC): %s\n", d.UpdatedAt.UTC().Format(common.TimeFormat))
		}
		if d.DeletedAt != nil {
			fmt.Fprintf(out, "\t\tDeleted (UTC): %s\n", d.DeletedAt.UTC().Format(common.TimeFormat))
		}
	}
	fmt.Fprintln(out)

	return nil
}
