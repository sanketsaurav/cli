package dictionaryitem

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
	DictionaryID string
	Key          string
}

// NewDescribeCommand returns a usable command registered under the parent.
func NewDescribeCommand(parent common.Registerer, globals *config.Data) *DescribeCommand {
	var c DescribeCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("describe", "Show detailed information about an item in a dictionary").Alias("get")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("dictionary-id", "The ID of the dictionary containing this item").Required().StringVar(&c.DictionaryID)
	c.CmdClause.Flag("key", "Item key").Required().StringVar(&c.Key)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *DescribeCommand) createInput() (*fastly.GetDictionaryItemInput, error) {
	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input := fastly.GetDictionaryItemInput{
		Service:    serviceID,
		Dictionary: c.DictionaryID,
		ItemKey:    c.Key,
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *DescribeCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	d, err := c.Globals.Client.GetDictionaryItem(input)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Service ID: %s\n", d.ServiceID)
	fmt.Fprintf(out, "Dictionary ID: %s\n", d.DictionaryID)
	fmt.Fprintf(out, "Key: %s\n", d.ItemKey)
	fmt.Fprintf(out, "Value: %v\n", d.ItemValue)

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
