package main

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

var UnitsNA = progress.Units{
	Notation:         "",
	NotationPosition: progress.UnitsNotationPositionAfter,
	Formatter:        func(int64) string { return "n/a" },
}

var UnitsCurrency = progress.Units{
	Notation:         "€ ",
	NotationPosition: progress.UnitsNotationPositionBefore,
	Formatter:        formatCurrency,
}

func formatCurrency(value int64) string {
	return fmt.Sprintf("%.2f", float32(value)/100)
}

func initProgress(trackers int) progress.Writer {

	// instantiate a Progress Writer and set up the options
	pw := progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetTrackerLength(7)
	pw.ShowETA(false)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.ShowTracker(false)
	pw.ShowValue(true)
	pw.SetMessageWidth(15)

	pw.SetNumTrackersExpected(trackers)
	pw.SetSortBy(progress.SortByMessage)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	// pw.Style().Options.PercentFormat = "%4.1f%%"

	return pw
}
