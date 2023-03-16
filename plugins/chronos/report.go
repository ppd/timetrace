package chronos

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dominikbraun/timetrace/core"
)

type ChronosBooking struct {
	ProjectName string `json:"project_name"`
	AccountName string `json:"account_name"`
	WorkingDay  int    `json:"working_day"`
	Duration    string `json:"duration"`
	Comment     string `json:"comment"`
}

type ChronosWorkingHours struct {
	WorkingDay int    `json:"working_day"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Breaks     string `json:"breaks"`
}

type ChronosImport struct {
	Bookings     []ChronosBooking      `json:"bookings"`
	WorkingHours []ChronosWorkingHours `json:"working_hours"`
}

func toChronosTimeFormat(duration time.Duration) string {
	totalInMinutes := int(duration.Minutes())
	wholeHours := totalInMinutes / 60
	minutes := totalInMinutes - wholeHours*60
	return fmt.Sprintf("%d:%d", wholeHours, minutes)
}

func roundUpDurationTo15Minutes(duration time.Duration) time.Duration {
	fifteenMinutes, _ := time.ParseDuration("15m")
	sevenAndAHalfMinutes, _ := time.ParseDuration("7.5m")
	roundedUpDuration := (duration + sevenAndAHalfMinutes).Round(fifteenMinutes)
	if roundedUpDuration.Minutes() < 15 {
		return fifteenMinutes
	} else {
		return roundedUpDuration
	}
}

func formatComments(comments []string) string {
	deduplicatedComments := make([]string, 0)
	seen := make(map[string]bool)

	for _, comment := range comments {
		_, haveComment := seen[comment]
		if !haveComment {
			deduplicatedComments = append(deduplicatedComments, comment)
			seen[comment] = true
		}
	}

	return strings.Join(deduplicatedComments, ", ")
}

func ChronosJson(t *core.Timetrace, reporter *core.Reporter, workingDay int) ([]byte, error) {
	ledger := BuildChronosLedger(reporter)
	bookings := make([]ChronosBooking, 0)
	totalDuration := time.Duration(0)

	for _, ledgerItem := range ledger {
		roundedDuration := roundUpDurationTo15Minutes(ledgerItem.TotalTime)
		totalDuration += roundedDuration
		project, err := t.LoadProject(ledgerItem.ProjectKey)
		if err != nil {
			return []byte{}, err
		}
		if project.ChronosProject == "" {
			return []byte{}, fmt.Errorf("no Chronos project configured for project %s", project.Key)
		}
		accountName := ledgerItem.ProjectKey
		if project.ChronosProject != "" {
			accountName = project.ChronosProject
		}
		booking := ChronosBooking{
			ProjectName: project.ChronosProject,
			AccountName: accountName,
			WorkingDay:  workingDay,
			Duration:    toChronosTimeFormat(roundedDuration),
			Comment:     formatComments(ledgerItem.Comments),
		}
		bookings = append(bookings, booking)
	}

	startWorkAt, _ := time.Parse("15:04", "8:45")
	breaks, _ := time.ParseDuration("1h")
	endWorkAt := startWorkAt.Add(totalDuration).Add(breaks)
	workingHours := ChronosWorkingHours{
		WorkingDay: workingDay,
		Start:      "8:45",
		End:        endWorkAt.Format("15:04"),
		Breaks:     toChronosTimeFormat(breaks),
	}

	importData := ChronosImport{
		Bookings:     bookings,
		WorkingHours: []ChronosWorkingHours{workingHours},
	}

	jsonData, err := json.MarshalIndent(importData, "", "  ")

	if err != nil {
		return nil, fmt.Errorf("could not serialize Chronos report")
	}

	return jsonData, nil
}
