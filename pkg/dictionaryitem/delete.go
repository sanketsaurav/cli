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

// DeleteCommand calls the Fastly API to delete dictionaries.
type DeleteCommand struct {
	common.Base
	manifest manifest.Data

	// required
	DictionaryID string
	Key          string
}

// NewDeleteCommand returns a usable command registered under the parent.
func NewDeleteCommand(parent common.Registerer, globals *config.Data) *DeleteCommand {
	var c DeleteCommand
	c.Globals = globals

	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("delete", "Delete an item from an dictionary").Alias("remove")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("dictionary-id", "The ID of the dictionary containing this item").Required().StringVar(&c.DictionaryID)
	c.CmdClause.Flag("key", "Item key").Required().StringVar(&c.Key)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *DeleteCommand) createInput() (*fastly.DeleteDictionaryItemInput, error) {
	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input := fastly.DeleteDictionaryItemInput{
		Service:    serviceID,
		Dictionary: c.DictionaryID,
		ItemKey:    c.Key,
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *DeleteCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	if err := c.Globals.Client.DeleteDictionaryItem(input); err != nil {
		return err
	}

	text.Success(out, "Deleted dictionary item %s (service %s dictionary %s)", input.ItemKey, input.Service, input.Dictionary)
	return nil
}
