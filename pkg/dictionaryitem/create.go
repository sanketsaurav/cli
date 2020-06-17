package dictionaryitem

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
	DictionaryID string
	Key          string
	Value        string
}

// NewCreateCommand returns a usable command registered under the parent.
func NewCreateCommand(parent common.Registerer, globals *config.Data) *CreateCommand {
	var c CreateCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("create", "Create an entry in an edge dictionary").Alias("add")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("dictionary-id", "The ID of the dictionary containing this item").Required().StringVar(&c.DictionaryID)
	c.CmdClause.Flag("key", "Item key, maximum 256 characters").Required().StringVar(&c.Key)
	c.CmdClause.Flag("value", "Item value, maximum 8000 characters").Required().StringVar(&c.Value)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *CreateCommand) createInput() (*fastly.CreateDictionaryItemInput, error) {

	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input := fastly.CreateDictionaryItemInput{
		Service:    serviceID,
		Dictionary: c.DictionaryID,
		ItemKey:    c.Key,
		ItemValue:  c.Value,
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *CreateCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	d, err := c.Globals.Client.CreateDictionaryItem(input)
	if err != nil {
		return err
	}

	text.Success(out, "Created dictionary item %s (service %s dictionary %s)", d.ItemKey, d.ServiceID, d.DictionaryID)
	return nil
}
