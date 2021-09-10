package app

import (
	"encoding/csv"
	"io"

	"github.com/antelman107/parser/top"
)

type CsvParserApp struct {
	k   uint32
	cfg CsvParserAppConfig
}

// CsvParserAppConfig represents config values required to process csv file and entity file.
// There are 3 steps of this app, described in Run.
type CsvParserAppConfig struct {
	// Column index of entity(user/repo) ID from events.csv file.
	EventsEntityColumnIndex int

	// Column index of event type from events.csv file.
	EventsEventTypeColumnIndex int

	// Column index of entity(user/repo) ID from entity file (actors.csv, repos.csv).
	EntityEntityColumnIndex int

	// Column index of entity(user/repo) name from entity file (actors.csv, repos.csv).
	EntityNameColumnIndex int

	// Event types to filter.
	EventTypes []string

	// Reader of events.csv
	EventsFileReader *csv.Reader

	// Reader of entity file (actors.csv, repos.csv).
	EntityFileReader *csv.Reader

	// Writer to output the results.
	Writer io.Writer
}

// NewCsvParserApp returns new CsvParserApp by k and config.
func NewCsvParserApp(k uint32, cfg CsvParserAppConfig) *CsvParserApp {
	return &CsvParserApp{
		k:   k,
		cfg: cfg,
	}
}

// Run executes CsvParserApp.
func (s *CsvParserApp) Run() error {
	// 1. Get IDs top list.
	list, err := s.getIDsTop()
	if err != nil {
		return err
	}

	// 2. Get named top list.
	namedList, err := s.getNamedTop(list)
	if err != nil {
		return err
	}

	// 3. Output results.
	return namedList.WriteResults(s.cfg.Writer)
}

// getIDsTop returns top.List by using heavykeeper.TopK againts events csv file.
func (s *CsvParserApp) getIDsTop() (top.List, error) {
	hk := top.GetHK(s.k, top.HkConfigDefault)

	// Iterate through csv lines.
	for i := 0; ; i++ {
		parts, err := s.cfg.EventsFileReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Check if this line matched by event types
		for j := range s.cfg.EventTypes {
			if s.cfg.EventTypes[j] == parts[s.cfg.EventsEventTypeColumnIndex] {
				// Increase counter of corresponding ID
				hk.Add(parts[s.cfg.EventsEntityColumnIndex])
				break
			}
		}
	}

	hk.Wait()

	return top.GetListFromHK(s.k, hk), nil
}

// getNamedTop returns same list as in, but the IDs of entities are replaced by
func (s *CsvParserApp) getNamedTop(in top.List) (top.List, error) {
	namedRecordsCount := uint32(0)

	// Read events csv lines
	for {
		// Check if all records are named already to avoid useless iterations.
		if namedRecordsCount == s.k {
			break
		}

		parts, err := s.cfg.EntityFileReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// For each value from list - check if this line matches entity ID from list.
		for j := range in {
			if in[j].Name == parts[s.cfg.EntityEntityColumnIndex] {
				in[j].Name = parts[s.cfg.EntityNameColumnIndex]
				namedRecordsCount++
				break
			}
		}
	}

	return in, nil
}
