package dictionaryitem_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/fastly/cli/pkg/app"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/mock"
	"github.com/fastly/cli/pkg/testutil"
	"github.com/fastly/cli/pkg/update"
	"github.com/fastly/go-fastly/fastly"
)

func TestDictionaryItemCreate(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary-item", "create", "--service-id", "123", "--dictionary-id", "1", "--key", "foo"},
			wantError: "error parsing arguments: required flag --value not provided",
		},
		{
			args:       []string{"dictionary-item", "create", "--service-id", "123", "--dictionary-id", "1", "--key", "foo", "--value", "bar"},
			api:        mock.API{CreateDictionaryItemFn: createDictionaryItemOK},
			wantOutput: "Created dictionary item foo (service 123 dictionary 1)",
		},
		{
			args:      []string{"dictionary-item", "create", "--service-id", "123", "--dictionary-id", "1", "--key", "foo", "--value", "bar"},
			api:       mock.API{CreateDictionaryItemFn: createDictionaryItemError},
			wantError: errTest.Error(),
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertStringContains(t, out.String(), testcase.wantOutput)
		})
	}
}

func TestDictionaryItemsList(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:       []string{"dictionary-item", "list", "--service-id", "123", "--dictionary-id", "1"},
			api:        mock.API{ListDictionaryItemsFn: listDictionaryItemsOK},
			wantOutput: listDictionaryItemsShortOutput,
		},
		{
			args:       []string{"dictionary-item", "list", "--service-id", "123", "--dictionary-id", "1", "--verbose"},
			api:        mock.API{ListDictionaryItemsFn: listDictionaryItemsOK},
			wantOutput: listDictionaryItemsVerboseOutput,
		},
		{
			args:       []string{"dictionary-item", "list", "--service-id", "123", "--dictionary-id", "1", "-v"},
			api:        mock.API{ListDictionaryItemsFn: listDictionaryItemsOK},
			wantOutput: listDictionaryItemsVerboseOutput,
		},
		{
			args:       []string{"dictionary-item", "--verbose", "list", "--service-id", "123", "--dictionary-id", "1"},
			api:        mock.API{ListDictionaryItemsFn: listDictionaryItemsOK},
			wantOutput: listDictionaryItemsVerboseOutput,
		},
		{
			args:      []string{"dictionary-item", "list", "--service-id", "123", "--dictionary-id", "1"},
			api:       mock.API{ListDictionaryItemsFn: listDictionaryItemsError},
			wantError: errTest.Error(),
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertString(t, testcase.wantOutput, out.String())
		})
	}
}

func TestDictionaryItemDescribe(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary-item", "describe", "--service-id", "123", "--dictionary-id", "1"},
			wantError: "error parsing arguments: required flag --key not provided",
		},
		{
			args:      []string{"dictionary-item", "describe", "--service-id", "123", "--dictionary-id", "1", "--key", "foo"},
			api:       mock.API{GetDictionaryItemFn: getDictionaryItemError},
			wantError: errTest.Error(),
		},
		{
			args:       []string{"dictionary-item", "describe", "--service-id", "123", "--dictionary-id", "1", "--key", "foo"},
			api:        mock.API{GetDictionaryItemFn: getDictionaryItemOK},
			wantOutput: describeDictionaryItemOutput,
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertString(t, testcase.wantOutput, out.String())
		})
	}
}

func TestDictionaryItemUpdate(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary-item", "update", "--service-id", "123", "--dictionary-id", "1", "--key", "foo"},
			wantError: "error parsing arguments: required flag --value not provided",
		},
		{
			args: []string{"dictionary-item", "update", "--service-id", "123", "--dictionary-id", "1", "--key", "foo", "--value", "baz"},
			api: mock.API{
				GetDictionaryItemFn:    getDictionaryItemOK,
				UpdateDictionaryItemFn: updateDictionaryItemError,
			},
			wantError: errTest.Error(),
		},
		{
			args: []string{"dictionary-item", "update", "--service-id", "123", "--dictionary-id", "1", "--key", "foo", "--value", "bar"},
			api: mock.API{
				GetDictionaryItemFn:    getDictionaryItemOK,
				UpdateDictionaryItemFn: updateDictionaryItemOK,
			},
			wantOutput: "Updated dictionary item foo (service 123 dictionary 1)",
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertStringContains(t, out.String(), testcase.wantOutput)
		})
	}
}

func TestDictionaryItemDelete(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary-item", "delete", "--service-id", "123", "--dictionary-id", "1"},
			wantError: "error parsing arguments: required flag --key not provided",
		},
		{
			args:      []string{"dictionary-item", "delete", "--service-id", "123", "--dictionary-id", "1", "--key", "foo"},
			api:       mock.API{DeleteDictionaryItemFn: deleteDictionaryItemError},
			wantError: errTest.Error(),
		},
		{
			args:       []string{"dictionary-item", "delete", "--service-id", "123", "--dictionary-id", "1", "--key", "foo"},
			api:        mock.API{DeleteDictionaryItemFn: deleteDictionaryItemOK},
			wantOutput: "Deleted dictionary item foo (service 123 dictionary 1)",
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertStringContains(t, out.String(), testcase.wantOutput)
		})
	}
}

var errTest = errors.New("fixture error")

func createDictionaryItemOK(i *fastly.CreateDictionaryItemInput) (*fastly.DictionaryItem, error) {
	s := fastly.DictionaryItem{
		ServiceID:    i.Service,
		DictionaryID: i.Dictionary,
		ItemKey:      i.ItemKey,
	}

	return &s, nil
}

func createDictionaryItemError(i *fastly.CreateDictionaryItemInput) (*fastly.DictionaryItem, error) {
	return nil, errTest
}

func listDictionaryItemsOK(i *fastly.ListDictionaryItemsInput) ([]*fastly.DictionaryItem, error) {
	return []*fastly.DictionaryItem{
		{
			ServiceID:    i.Service,
			DictionaryID: i.Dictionary,
			ItemKey:      "foo",
			ItemValue:    "bar",
		},
		{
			ServiceID:    i.Service,
			DictionaryID: i.Dictionary,
			ItemKey:      "fiz",
			ItemValue:    "buz",
		},
	}, nil
}

func listDictionaryItemsError(i *fastly.ListDictionaryItemsInput) ([]*fastly.DictionaryItem, error) {
	return nil, errTest
}

var listDictionaryItemsShortOutput = strings.TrimSpace(`
SERVICE  DICTIONARY ID  KEY  VALUE
123      1              foo  bar
123      1              fiz  buz
`) + "\n"

var listDictionaryItemsVerboseOutput = strings.TrimSpace(`
Fastly API token not provided
Fastly API endpoint: https://api.fastly.com
Service ID: 123
Dictionary ID: 1
	Item 1/2
		Key: foo
		Value: bar
	Item 2/2
		Key: fiz
		Value: buz
`) + "\n\n"

func getDictionaryItemOK(i *fastly.GetDictionaryItemInput) (*fastly.DictionaryItem, error) {
	return &fastly.DictionaryItem{
		ServiceID:    i.Service,
		DictionaryID: "1",
		ItemKey:      "foo",
		ItemValue:    "bar",
	}, nil
}

func getDictionaryItemError(i *fastly.GetDictionaryItemInput) (*fastly.DictionaryItem, error) {
	return nil, errTest
}

var describeDictionaryItemOutput = strings.TrimSpace(`
Service ID: 123
Dictionary ID: 1
Key: foo
Value: bar
`) + "\n"

func updateDictionaryItemOK(i *fastly.UpdateDictionaryItemInput) (*fastly.DictionaryItem, error) {
	return &fastly.DictionaryItem{
		ServiceID:    i.Service,
		DictionaryID: i.Dictionary,
		ItemKey:      "foo",
		ItemValue:    "bar",
	}, nil
}

func updateDictionaryItemError(i *fastly.UpdateDictionaryItemInput) (*fastly.DictionaryItem, error) {
	return nil, errTest
}

func deleteDictionaryItemOK(i *fastly.DeleteDictionaryItemInput) error {
	return nil
}

func deleteDictionaryItemError(i *fastly.DeleteDictionaryItemInput) error {
	return errTest
}
