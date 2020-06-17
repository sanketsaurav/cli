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

// UpdateCommand calls the Fastly API to update dictionaries.
type UpdateCommand struct {
	common.Base
	manifest manifest.Data

	// required
	DictionaryID string
	Key          string
	Value        string
}

// NewUpdateCommand returns a usable command registered under the parent.
func NewUpdateCommand(parent common.Registerer, globals *config.Data) *UpdateCommand {
	var c UpdateCommand
	c.Globals = globals

	c.CmdClause = parent.Command("update", "Update an item in an dictionary")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("dictionary-id", "The ID of the dictionary containing this item").Required().StringVar(&c.DictionaryID)
	c.CmdClause.Flag("key", "Item key").Required().StringVar(&c.Key)
	c.CmdClause.Flag("value", "Item value").Required().StringVar(&c.Value)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *UpdateCommand) createInput() (*fastly.UpdateDictionaryItemInput, error) {
	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input := fastly.UpdateDictionaryItemInput{
		Service:    serviceID,
		Dictionary: c.DictionaryID,
		ItemKey:    c.Key,
		ItemValue:  c.Value,
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *UpdateCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	i, err := c.Globals.Client.UpdateDictionaryItem(input)
	if err != nil {
		return err
	}

	text.Success(out, "Updated dictionary item %s (service %s dictionary %s)", i.ItemKey, i.ServiceID, i.DictionaryID)
	return nil
}
