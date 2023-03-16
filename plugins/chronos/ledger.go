package chronos

import (
	"strings"
	"time"

	"github.com/dominikbraun/timetrace/core"
)

func getModuleKey(record *core.Record) string {
	tokens := strings.Split(record.Project.Key, "@")
	if len(tokens) > 1 {
		return tokens[0]
	} else {
		return ""
	}
}

type ChronosProjectLedgerItem struct {
	TotalTime  time.Duration
	Comments   []string
	ProjectKey string
}

func BuildChronosLedger(reporter *core.Reporter) []ChronosProjectLedgerItem {
	ledger := make([]ChronosProjectLedgerItem, 0)

	for projectKey, records := range reporter.GetReport() {
		totalDuration := time.Duration(0)
		comments := make([]string, 0)

		for _, record := range records {
			totalDuration += record.Duration()

			if record.Project.IsModule() {
				comments = append(comments, getModuleKey(record))
			}

			if len(record.Tags) > 0 {
				comments = append(comments, record.Tags...)
			}
		}

		ledger = append(ledger, ChronosProjectLedgerItem{
			TotalTime:  reporter.GetTotals()[projectKey],
			Comments:   comments,
			ProjectKey: projectKey,
		})
	}

	return ledger
}
