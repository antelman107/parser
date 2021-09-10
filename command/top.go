package command

import (
	"encoding/csv"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/antelman107/parser/app"
	"github.com/spf13/cobra"
)

const (
	paramK                          = "k"
	paramEventTypes                 = "event_types"
	paramEventsEntityColumnIndex    = "events_entity_column_index"
	paramEventsEventTypeColumnIndex = "events_event_type_column_index"
	paramEntityEntityColumnIndex    = "entity_entity_column_index"
	paramEntityNameColumnIndex      = "entity_name_column_index"
	paramEventsFile                 = "events_file"
	paramEntityFile                 = "entity_file"
)

var topCmd = &cobra.Command{
	Use:  "top",
	RunE: topRunE,
}

func init() {
	topCmd.PersistentFlags().Uint32(paramK, 10, "Number of leaderboard participants")
	topCmd.PersistentFlags().StringSlice(paramEventTypes, []string{"PushEvent", "PullRequestEvent"}, "Filtered events names")
	topCmd.PersistentFlags().Int(paramEventsEventTypeColumnIndex, 1, "Index of event file column with event type")
	topCmd.PersistentFlags().Int(paramEventsEntityColumnIndex, 2, "")
	topCmd.PersistentFlags().Int(paramEntityEntityColumnIndex, 0, "")
	topCmd.PersistentFlags().Int(paramEntityNameColumnIndex, 1, "")
	topCmd.PersistentFlags().String(paramEventsFile, "./data/events.csv", "")
	topCmd.PersistentFlags().String(paramEntityFile, "./data/actors.csv", "")
}

// This exists only to replace it in tests.
var csvAppBuilder func(k uint32, cfg app.CsvParserAppConfig) app.CsvApp

func topRunE(cmd *cobra.Command, _ []string) error {
	// Errors are always nil, default value is set
	k, _ := cmd.PersistentFlags().GetUint32(paramK)
	eventTypes, _ := cmd.PersistentFlags().GetStringSlice(paramEventTypes)
	eventsEntityColumnIndex, _ := cmd.PersistentFlags().GetInt(paramEventsEntityColumnIndex)
	eventsEventTypeColumnIndex, _ := cmd.PersistentFlags().GetInt(paramEventsEventTypeColumnIndex)
	entityEntityColumnIndex, _ := cmd.PersistentFlags().GetInt(paramEntityEntityColumnIndex)
	entityNameColumnIndex, _ := cmd.PersistentFlags().GetInt(paramEntityNameColumnIndex)
	eventsFileName, _ := cmd.PersistentFlags().GetString(paramEventsFile)
	entityFileName, _ := cmd.PersistentFlags().GetString(paramEntityFile)

	eventsFile, err := os.Open(eventsFileName)
	if err != nil {
		return fmt.Errorf("failed to open events file: %w", err)
	}

	entityFile, err := os.Open(entityFileName)
	if err != nil {
		return fmt.Errorf("failed to open entity file: %w", err)
	}

	a := csvAppBuilder(
		k,
		app.CsvParserAppConfig{
			EventsEntityColumnIndex:    eventsEntityColumnIndex,
			EventsEventTypeColumnIndex: eventsEventTypeColumnIndex,
			EntityEntityColumnIndex:    entityEntityColumnIndex,
			EntityNameColumnIndex:      entityNameColumnIndex,
			EventTypes:                 eventTypes,
			EventsFileReader:           csv.NewReader(eventsFile),
			EntityFileReader:           csv.NewReader(entityFile),
			Writer:                     tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
		},
	)

	if err := a.Run(); err != nil {
		return err
	}

	return nil
}
