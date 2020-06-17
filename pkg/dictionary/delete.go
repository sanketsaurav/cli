package dictionary

import (
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/fastly"
)

// DeleteCommand calls the Fastly API to delete dictionaries.
type DeleteCommand struct {
	common.Base
	manifest manifest.Data

	// required
	DictionaryName string // Can't shaddow common.Base method Name().
	Version        int
}

// NewDeleteCommand returns a usable command registered under the parent.
func NewDeleteCommand(parent common.Registerer, globals *config.Data) *DeleteCommand {
	var c DeleteCommand
	c.Globals = globals

	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("delete", "Delete a dictionary on a Fastly service version").Alias("remove")

	c.CmdClause.Flag("name", "Domain name").Short('n').Required().StringVar(&c.DictionaryName)
	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("version", "Number of service version").Required().IntVar(&c.Version)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *DeleteCommand) createInput() (*fastly.DeleteDictionaryInput, error) {
	var input fastly.DeleteDictionaryInput

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
func (c *DeleteCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	if err := c.Globals.Client.DeleteDictionary(input); err != nil {
		return err
	}

	text.Success(out, "Deleted dictionary %s (service %s version %d)", input.Name, input.Service, input.Version)
	return nil
}
