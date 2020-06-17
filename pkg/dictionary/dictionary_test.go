package dictionary

import (
	"testing"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/mock"
	"github.com/fastly/cli/pkg/testutil"
	"github.com/fastly/go-fastly/fastly"
)

func TestCreateDictionaryInput(t *testing.T) {
	for _, testcase := range []struct {
		name      string
		cmd       *CreateCommand
		want      *fastly.CreateDictionaryInput
		wantError string
	}{
		{
			name: "required values set flag serviceID",
			cmd:  createCommandRequired(),
			want: &fastly.CreateDictionaryInput{
				Service: "123",
				Version: 2,
				Name:    "dictionary",
			},
		},
		{
			name: "all values set flag serviceID",
			cmd:  createCommandOK(),
			want: &fastly.CreateDictionaryInput{
				Service:   "123",
				Version:   2,
				Name:      "dictionary",
				WriteOnly: fastly.CBool(true),
			},
		},
		{
			name:      "error missing serviceID",
			cmd:       createCommandMissingServiceID(),
			want:      nil,
			wantError: errors.ErrNoServiceID.Error(),
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			have, err := testcase.cmd.createInput()
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertEqual(t, testcase.want, have)
		})
	}
}

func TestUpdateDictionaryInput(t *testing.T) {
	for _, testcase := range []struct {
		name      string
		cmd       *UpdateCommand
		api       mock.API
		want      *fastly.UpdateDictionaryInput
		wantError string
	}{
		{
			name: "no updates",
			cmd:  updateCommandNoUpdates(),
			api:  mock.API{GetDictionaryFn: getDictionaryOK},
			want: &fastly.UpdateDictionaryInput{
				Service: "123",
				Version: 2,
				Name:    "dictionary",
			},
		},
		{
			name: "all values set flag serviceID",
			cmd:  updateCommandAll(),
			api:  mock.API{GetDictionaryFn: getDictionaryOK},
			want: &fastly.UpdateDictionaryInput{
				Service: "123",
				Version: 2,
				Name:    "dictionary",
				NewName: "new-dictionary",
			},
		},
		{
			name:      "error missing serviceID",
			cmd:       updateCommandMissingServiceID(),
			want:      nil,
			wantError: errors.ErrNoServiceID.Error(),
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.cmd.Base.Globals.Client = testcase.api

			have, err := testcase.cmd.createInput()
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertEqual(t, testcase.want, have)
		})
	}
}

func createCommandOK() *CreateCommand {
	return &CreateCommand{
		manifest:       manifest.Data{Flag: manifest.Flag{ServiceID: "123"}},
		Version:        2,
		DictionaryName: "dictionary",
		WriteOnly:      common.OptionalBool{Optional: common.Optional{Valid: true}, Value: true},
	}
}

func createCommandRequired() *CreateCommand {
	return &CreateCommand{
		manifest:       manifest.Data{Flag: manifest.Flag{ServiceID: "123"}},
		Version:        2,
		DictionaryName: "dictionary",
	}
}

func createCommandMissingServiceID() *CreateCommand {
	res := createCommandOK()
	res.manifest = manifest.Data{}
	return res
}

func updateCommandNoUpdates() *UpdateCommand {
	return &UpdateCommand{
		Base:           common.Base{Globals: &config.Data{Client: nil}},
		manifest:       manifest.Data{Flag: manifest.Flag{ServiceID: "123"}},
		Version:        2,
		DictionaryName: "dictionary",
	}
}

func updateCommandAll() *UpdateCommand {
	return &UpdateCommand{
		Base:           common.Base{Globals: &config.Data{Client: nil}},
		manifest:       manifest.Data{Flag: manifest.Flag{ServiceID: "123"}},
		Version:        2,
		DictionaryName: "dictionary",
		NewName:        common.OptionalString{Optional: common.Optional{Valid: true}, Value: "new-dictionary"},
	}
}

func updateCommandMissingServiceID() *UpdateCommand {
	res := updateCommandAll()
	res.manifest = manifest.Data{}
	return res
}

func getDictionaryOK(i *fastly.GetDictionaryInput) (*fastly.Dictionary, error) {
	return &fastly.Dictionary{
		ServiceID: i.Service,
		Version:   i.Version,
		ID:        "123",
		Name:      "dictionary",
		WriteOnly: false,
	}, nil
}
