// Copyright 2013 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timeago

import (
	"testing"
	"time"
)

//Base time for testing
var tBase = time.Date(2013, 8, 30, 12, 0, 0, 0, time.UTC)

// Test data for TestFormatReference
var formatReferenceTests = []struct {
	t        time.Time // input time
	ref      time.Time // input reference
	cfg      Config    // input cfguage
	expected string    // expected result
}{
	// Lang
	{tBase, tBase, NoMax(English), "about a second ago"},
	{tBase, tBase, NoMax(French), "il y a environ une seconde"},
	{tBase, tBase, NoMax(Chinese), "1 秒前"},
	{tBase, tBase, NoMax(Russian), "около секунды назад"},

	// Thresholds (English)
	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(English), "about a second ago"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(English), "2 seconds ago"},
	{tBase, tBase.Add(1 * time.Minute).Add(-500000001), NoMax(English), "59 seconds ago"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(English), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(English), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(English), "2 minutes ago"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second).Add(-1), NoMax(English), "59 minutes ago"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(English), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(English), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(English), "2 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(English), "23 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(English), "one day ago"},
	{tBase, tBase.Add(36 * time.Hour).Add(-1), NoMax(English), "one day ago"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(English), "2 days ago"},
	{tBase, tBase.Add(30 * Day).Add(-12*time.Hour - 1), NoMax(English), "29 days ago"},
	{tBase, tBase.Add(30 * Day), NoMax(English), "one month ago"},
	{tBase, tBase.Add(45 * Day).Add(-1), NoMax(English), "one month ago"},
	{tBase, tBase.Add(45 * Day), NoMax(English), "2 months ago"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(English), "11 months ago"},
	{tBase, tBase.Add(Year), NoMax(English), "one year ago"},
	{tBase, tBase.Add(547 * Day), NoMax(English), "one year ago"},
	{tBase, tBase.Add(548 * Day), NoMax(English), "2 years ago"},
	{tBase, tBase.Add(10 * Year), NoMax(English), "10 years ago"},

	// Max (English)
	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(English), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(English, 90*time.Minute, ""), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(English, 90*time.Minute, "2006-01-02"), "2013-08-30"},

	// Future (English)
	{tBase.Add(Day), tBase, NoMax(English), "in one day"},

	// Thresholds (Russian)
	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(Russian), "секунду назад"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(Russian), "2 секунды назад"},
	{tBase, tBase.Add(10*time.Second + 500000000), NoMax(Russian), "11 секунд назад"},
	{tBase, tBase.Add(1 * time.Minute).Add(-500000001), NoMax(Russian), "59 секунд назад"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(Russian), "около минуты назад"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(Russian), "около минуты назад"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(Russian), "2 минуты назад"},
	{tBase, tBase.Add(10*time.Minute + 30*time.Second), NoMax(Russian), "11 минут назад"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second).Add(-1), NoMax(Russian), "59 минут назад"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(Russian), "около часа назад"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(Russian), "около часа назад"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(Russian), "2 часа назад"},
	{tBase, tBase.Add(11*time.Hour + 30*time.Minute).Add(-1), NoMax(Russian), "11 часов назад"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(Russian), "23 часа назад"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(Russian), "день назад"},
	{tBase, tBase.Add(36 * time.Hour).Add(-1), NoMax(Russian), "день назад"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(Russian), "2 дня назад"},
	{tBase, tBase.Add(12 * Day).Add(-12*time.Hour - 1), NoMax(Russian), "11 дней назад"},
	{tBase, tBase.Add(30 * Day).Add(-12*time.Hour - 1), NoMax(Russian), "29 дней назад"},
	{tBase, tBase.Add(30 * Day), NoMax(Russian), "месяц назад"},
	{tBase, tBase.Add(45 * Day).Add(-1), NoMax(Russian), "месяц назад"},
	{tBase, tBase.Add(45 * Day), NoMax(Russian), "2 месяца назад"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(Russian), "11 месяцев назад"},
	{tBase, tBase.Add(Year), NoMax(Russian), "год назад"},
	{tBase, tBase.Add(547 * Day), NoMax(Russian), "год назад"},
	{tBase, tBase.Add(548 * Day), NoMax(Russian), "2 года назад"},
	{tBase, tBase.Add(10 * Year), NoMax(Russian), "10 лет назад"},

	// Max (Russian)
	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(Russian), "около часа назад"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(Russian, 90*time.Minute, ""), "около часа назад"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(Russian, 90*time.Minute, "2006-01-02"), "2013-08-30"},

	// Future (Russian)
	{tBase.Add(Day), tBase, NoMax(Russian), "через день"},
	{tBase.Add(2 * Day), tBase, NoMax(Russian), "через 2 дня"},
}

// Test the FormatReference method
func TestFormatReference(t *testing.T) {
	for i, tt := range formatReferenceTests {
		actual := tt.cfg.FormatReference(tt.t, tt.ref)
		if actual != tt.expected {
			t.Errorf("%d) FormatReference(%s,%s): expected '%s', actual '%s'", i+1, tt.t, tt.ref, tt.expected, actual)
		}
	}
}
