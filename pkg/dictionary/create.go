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

// CreateCommand calls the Fastly API to create dictionaries.
type CreateCommand struct {
	common.Base
	manifest manifest.Data

	// required
	DictionaryName string // Can't shaddow common.Base method Name().
	Version        int

	// optional
	WriteOnly common.OptionalBool
}

// NewCreateCommand returns a usable command registered under the parent.
func NewCreateCommand(parent common.Registerer, globals *config.Data) *CreateCommand {
	var c CreateCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("create", "Create a dictionary on a Fastly service version").Alias("add")

	c.CmdClause.Flag("name", "Dictionary name").Short('n').Required().StringVar(&c.DictionaryName)
	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("version", "Number of service version").Required().IntVar(&c.Version)
	c.CmdClause.Flag("write-only", "Determines if items in the dictionary are readable or not.").Action(c.WriteOnly.Set).BoolVar(&c.WriteOnly.Value)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *CreateCommand) createInput() (*fastly.CreateDictionaryInput, error) {
	var input fastly.CreateDictionaryInput

	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input.Service = serviceID
	input.Version = c.Version
	input.Name = c.DictionaryName

	if c.WriteOnly.Valid {
		input.WriteOnly = fastly.CBool(c.WriteOnly.Value)
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *CreateCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	d, err := c.Globals.Client.CreateDictionary(input)
	if err != nil {
		return err
	}

	text.Success(out, "Created dictionary %s (service %s version %d)", d.Name, d.ServiceID, d.Version)
	return nil
}
