package app

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strings"
	"testing"
	"testing/iotest"
	"text/tabwriter"

	"github.com/antelman107/parser/top"
	"github.com/stretchr/testify/assert"
)

const testEventsContent = `id,type,actor_id,repo_id
11185376329,PushEvent,8422699,224252202
11185376333,CreateEvent,8422699,231161852
11185376335,PushEvent,2631623,155254893
11185376336,PushEvent,2631623,231065965
11185376338,PushEvent,2631623,225080339
11185376339,PushEvent,5954907,160083795
11185376341,WatchEvent,5954907,221552739
11185376342,PushEvent,5954907,230923653
11185376343,PushEvent,5954907,107471694
11185376344,PushEvent,8422699,223831715
11185376344,PullRequestEvent,8422699,223831715`

const testActorsContent = `id,username
8422699,Apexal
53201765,ArturoCamacho0
2631623,onosendi
52553915,anggi1234
31390726,AdrianWilczynski
5954907,awesomekling
10052381,PercussiveElbow
30060991,m41na
8517910,LombiqBot`

// This is more like integration test for Run.
// Different top parameters are tested in getTopList.
func TestCsvParserAppRun(t *testing.T) {
	b := &bytes.Buffer{}
	err := NewCsvParserApp(3, CsvParserAppConfig{
		EventsEntityColumnIndex:    2,
		EventsEventTypeColumnIndex: 1,
		EntityEntityColumnIndex:    0,
		EntityNameColumnIndex:      1,
		EventTypes:                 []string{"PushEvent", "PullRequestEvent"},
		EventsFileReader:           csv.NewReader(strings.NewReader(testEventsContent)),
		EntityFileReader:           csv.NewReader(strings.NewReader(testActorsContent)),
		Writer:                     tabwriter.NewWriter(b, 1, 1, 1, ' ', 0),
	}).Run()

	wantOutput := `onosendi       3
awesomekling   3
Apexal         3`

	assert.Nil(t, err)
	assert.Equal(t, wantOutput, strings.Trim(b.String(), " \n"))

}

func TestCsvParserAppRun_EmptyEntityFile(t *testing.T) {
	b := &bytes.Buffer{}
	err := NewCsvParserApp(3, CsvParserAppConfig{
		EventsEntityColumnIndex:    2,
		EventsEventTypeColumnIndex: 1,
		EntityEntityColumnIndex:    0,
		EntityNameColumnIndex:      1,
		EventTypes:                 []string{"PushEvent", "PullRequestEvent"},
		EventsFileReader:           csv.NewReader(strings.NewReader(testEventsContent)),
		EntityFileReader:           csv.NewReader(strings.NewReader("")),
		Writer:                     b,
	}).Run()

	assert.Nil(t, err)
}

func TestCsvParserAppRunFailed_EventsReaderError(t *testing.T) {
	b := &bytes.Buffer{}
	wantErr := errors.New("some error")
	err := NewCsvParserApp(3, CsvParserAppConfig{
		EventsFileReader: csv.NewReader(iotest.ErrReader(wantErr)),
		Writer:           b,
	}).Run()

	assert.Equal(t, wantErr, err)
}

func TestCsvParserAppRunFailed_EntityFileReaderError(t *testing.T) {
	b := &bytes.Buffer{}
	wantErr := errors.New("some error")
	err := NewCsvParserApp(3, CsvParserAppConfig{
		EventsEntityColumnIndex:    2,
		EventsEventTypeColumnIndex: 1,
		EntityEntityColumnIndex:    0,
		EntityNameColumnIndex:      1,
		EventTypes:                 []string{"PushEvent", "PullRequestEvent"},
		EventsFileReader:           csv.NewReader(strings.NewReader(testEventsContent)),
		EntityFileReader:           csv.NewReader(iotest.ErrReader(wantErr)),
		Writer:                     b,
	}).Run()

	assert.Equal(t, wantErr, err)
}

type errWriter string

func (e errWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New(string(e))
}

func TestCsvParserAppRunFailed_OutputError(t *testing.T) {
	err := NewCsvParserApp(3, CsvParserAppConfig{
		EventsEntityColumnIndex:    2,
		EventsEventTypeColumnIndex: 1,
		EntityEntityColumnIndex:    0,
		EntityNameColumnIndex:      1,
		EventTypes:                 []string{"PushEvent", "PullRequestEvent"},
		EventsFileReader:           csv.NewReader(strings.NewReader(testEventsContent)),
		EntityFileReader:           csv.NewReader(strings.NewReader(testActorsContent)),
		Writer:                     errWriter("some error"),
	}).Run()

	assert.Equal(t, "some error", err.Error())
}

func TestCsvParserAppGetIDsTop(t *testing.T) {
	for name, tc := range map[string]struct {
		eventsEntityColumnIndex    int
		eventsEventTypeColumnIndex int
		entityEntityColumnIndex    int
		entityNameColumnIndex      int
		eventTypes                 []string
		eventsFileReader           *csv.Reader
		entityFileReader           *csv.Reader
		list                       top.List
	}{
		"user 3 PushEvent": {
			eventsEntityColumnIndex:    2,
			eventsEventTypeColumnIndex: 1,
			entityEntityColumnIndex:    0,
			entityNameColumnIndex:      1,
			eventTypes:                 []string{"PushEvent"},
			eventsFileReader:           csv.NewReader(strings.NewReader(testEventsContent)),
			list: top.List{
				{"2631623", 3},
				{"5954907", 3},
				{"8422699", 2},
			},
		},
		"repo 2 WatchEvent": {
			eventsEntityColumnIndex:    2,
			eventsEventTypeColumnIndex: 1,
			entityEntityColumnIndex:    0,
			entityNameColumnIndex:      1,
			eventTypes:                 []string{"PushEvent"},
			eventsFileReader:           csv.NewReader(strings.NewReader(testEventsContent)),
			list: top.List{
				{"2631623", 3},
				{"5954907", 3},
				{"8422699", 2},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			list, err := NewCsvParserApp(3, CsvParserAppConfig{
				EventsEntityColumnIndex:    tc.eventsEntityColumnIndex,
				EventsEventTypeColumnIndex: tc.eventsEventTypeColumnIndex,
				EntityEntityColumnIndex:    tc.entityEntityColumnIndex,
				EntityNameColumnIndex:      tc.entityNameColumnIndex,
				EventTypes:                 tc.eventTypes,
				EventsFileReader:           tc.eventsFileReader,
				EntityFileReader:           tc.entityFileReader,
			}).getIDsTop()

			assert.Nil(t, err)
			assert.Equal(t, tc.list, list)
		})
	}
}
