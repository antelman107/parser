package command

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/antelman107/parser/app"
	"github.com/antelman107/parser/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTopRunE(t *testing.T) {
	cmd := topCmd
	_ = cmd.PersistentFlags().Set("k", "1")
	_ = cmd.PersistentFlags().Set("event_types", "PushEvent,PullRequestEvent")
	_ = cmd.PersistentFlags().Set("events_entity_column_index", "2")
	_ = cmd.PersistentFlags().Set("events_event_type_column_index", "1")
	_ = cmd.PersistentFlags().Set("entity_entity_column_index", "0")
	_ = cmd.PersistentFlags().Set("entity_name_column_index", "0")
	eventsFile, err := ioutil.TempFile("/tmp", "events")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove(eventsFile.Name())
	}()
	actorsFile, err := ioutil.TempFile("/tmp", "actors")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove(actorsFile.Name())
	}()
	_ = cmd.PersistentFlags().Set("events_file", eventsFile.Name())
	_ = cmd.PersistentFlags().Set("entity_file", actorsFile.Name())

	appMock := mocks.NewCsvAppMock()
	appMock.On("Run").Return(nil)

	builderMock := mocks.NewCsvAppBuilderMock()
	builderMock.On("Build", uint32(1), mock.MatchedBy(func(c app.CsvParserAppConfig) bool {
		assert.Equal(t, c.EventsEntityColumnIndex, 2)
		assert.Equal(t, c.EventsEventTypeColumnIndex, 1)
		assert.Equal(t, c.EntityEntityColumnIndex, 0)
		assert.Equal(t, c.EntityNameColumnIndex, 0)
		return true
	})).Return(appMock)

	csvAppBuilder = builderMock.Build
	err = topRunE(cmd, nil)

	assert.NoError(t, err)
	builderMock.AssertExpectations(t)
	appMock.AssertExpectations(t)
}

func TestTopRunEFailed_Run(t *testing.T) {
	cmd := topCmd
	_ = cmd.PersistentFlags().Set("k", "1")
	_ = cmd.PersistentFlags().Set("event_types", "PushEvent,PullRequestEvent")
	_ = cmd.PersistentFlags().Set("events_entity_column_index", "2")
	_ = cmd.PersistentFlags().Set("events_event_type_column_index", "1")
	_ = cmd.PersistentFlags().Set("entity_entity_column_index", "0")
	_ = cmd.PersistentFlags().Set("entity_name_column_index", "0")
	eventsFile, err := ioutil.TempFile("/tmp", "events")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove(eventsFile.Name())
	}()
	actorsFile, err := ioutil.TempFile("/tmp", "actors")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove(actorsFile.Name())
	}()
	_ = cmd.PersistentFlags().Set("events_file", eventsFile.Name())
	_ = cmd.PersistentFlags().Set("entity_file", actorsFile.Name())

	wantErr := errors.New("some error")
	appMock := mocks.NewCsvAppMock()
	appMock.On("Run").Return(wantErr)

	builderMock := mocks.NewCsvAppBuilderMock()
	builderMock.On("Build", uint32(1), mock.MatchedBy(func(c app.CsvParserAppConfig) bool {
		assert.Equal(t, c.EventsEntityColumnIndex, 2)
		assert.Equal(t, c.EventsEventTypeColumnIndex, 1)
		assert.Equal(t, c.EntityEntityColumnIndex, 0)
		assert.Equal(t, c.EntityNameColumnIndex, 0)
		return true
	})).Return(appMock)

	csvAppBuilder = builderMock.Build
	err = topRunE(cmd, nil)

	assert.Equal(t, wantErr, err)
	builderMock.AssertExpectations(t)
	appMock.AssertExpectations(t)
}

func TestTopRunEFailed_NoEventsFile(t *testing.T) {
	cmd := topCmd
	_ = cmd.PersistentFlags().Set("events_file", "i_think_it_doesnt_exist")

	err := topRunE(cmd, nil)

	assert.Contains(t, err.Error(), "failed to open events file:")
}

func TestTopRunEFailed_NoEntityFile(t *testing.T) {
	cmd := topCmd
	eventsFile, err := ioutil.TempFile("/tmp", "events")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove(eventsFile.Name())
	}()
	_ = cmd.PersistentFlags().Set("events_file", eventsFile.Name())
	_ = cmd.PersistentFlags().Set("entity_file", "i_think_it_doesnt_exist")

	err = topRunE(cmd, nil)

	assert.Contains(t, err.Error(), "failed to open entity file:")
}
