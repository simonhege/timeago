// Copyright 2013 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timeago

import (
	"testing"
	"time"
)

// Base time for testing
var tBase = time.Date(2013, 8, 30, 12, 0, 0, 0, time.UTC)

// Test data for TestFormatReference
var formatReferenceTests = []struct {
	t        time.Time // input time
	ref      time.Time // input reference
	cfg      Config    // input cfguage
	expected string    // expected result
}{
	// Lang
	{tBase, tBase, NoMax(EnglishUS), "about a second ago"},
	{tBase, tBase, NoMax(EnglishUK), "about a second ago"},
	{tBase, tBase, NoMax(French), "il y a environ une seconde"},
	{tBase, tBase, NoMax(Chinese), "1 秒前"},
	{tBase, tBase, NoMax(Portuguese), "há menos de um segundo"},
	{tBase, tBase, NoMax(German), "vor einer Sekunde"},
	{tBase, tBase, NoMax(Turkish), "yaklaşık bir saniye önce"},

	// Thresholds
	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(EnglishUS), "about a second ago"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(EnglishUS), "2 seconds ago"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(EnglishUS), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(EnglishUS), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(EnglishUS), "2 minutes ago"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(EnglishUS), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(EnglishUS), "2 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(EnglishUS), "23 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(EnglishUS), "one day ago"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(EnglishUS), "2 days ago"},
	{tBase, tBase.Add(30 * 24 * time.Hour), NoMax(EnglishUS), "one month ago"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(EnglishUS), "11 months ago"},
	{tBase, tBase.Add(Year), NoMax(EnglishUS), "one year ago"},
	{tBase, tBase.Add(547 * Day), NoMax(EnglishUS), "one year ago"},
	{tBase, tBase.Add(548 * Day), NoMax(EnglishUS), "2 years ago"},
	{tBase, tBase.Add(10 * Year), NoMax(EnglishUS), "10 years ago"},

	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(EnglishUK), "about a second ago"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(EnglishUK), "2 seconds ago"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(EnglishUK), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(EnglishUK), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(EnglishUK), "2 minutes ago"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(EnglishUK), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(EnglishUK), "2 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(EnglishUK), "23 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(EnglishUK), "one day ago"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(EnglishUK), "2 days ago"},
	{tBase, tBase.Add(30 * 24 * time.Hour), NoMax(EnglishUK), "one month ago"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(EnglishUK), "11 months ago"},
	{tBase, tBase.Add(Year), NoMax(EnglishUK), "one year ago"},
	{tBase, tBase.Add(547 * Day), NoMax(EnglishUK), "one year ago"},
	{tBase, tBase.Add(548 * Day), NoMax(EnglishUK), "2 years ago"},
	{tBase, tBase.Add(10 * Year), NoMax(EnglishUK), "10 years ago"},

	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(Portuguese), "há uma hora"},
	{tBase, tBase.Add(45 * 24 * time.Hour).Add(-1), NoMax(Portuguese), "há um mês"},
	{tBase, tBase.Add(36 * time.Hour).Add(-1), NoMax(Portuguese), "há um dia"},
	{tBase, tBase.Add(1 * time.Minute).Add(-500000001), NoMax(Portuguese), "há 59 segundos"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second).Add(-1), NoMax(Portuguese), "há 59 minutos"},
	{tBase, tBase.Add(30 * 24 * time.Hour).Add(-12*time.Hour - 1), NoMax(Portuguese), "há 29 dias"},
	{tBase, tBase.Add(45 * 24 * time.Hour), NoMax(Portuguese), "há 2 meses"},
	{tBase, tBase.Add(10 * Year), NoMax(Portuguese), "há 10 anos"},

	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(German), "vor einer Sekunde"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(German), "vor 2 Sekunden"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(German), "vor einer Minute"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(German), "vor einer Minute"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(German), "vor 2 Minuten"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(German), "vor einer Stunde"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(German), "vor 2 Stunden"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(German), "vor 23 Stunden"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(German), "vor einem Tag"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(German), "vor 2 Tagen"},
	{tBase, tBase.Add(30 * 24 * time.Hour), NoMax(German), "vor einem Monat"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(German), "vor 11 Monaten"},
	{tBase, tBase.Add(Year), NoMax(German), "vor einem Jahr"},
	{tBase, tBase.Add(547 * Day), NoMax(German), "vor einem Jahr"},
	{tBase, tBase.Add(548 * Day), NoMax(German), "vor 2 Jahren"},
	{tBase, tBase.Add(10 * Year), NoMax(German), "vor 10 Jahren"},

	// Turkish Thresholds
	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(Turkish), "yaklaşık bir saniye önce"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(Turkish), "2 saniye önce"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(Turkish), "yaklaşık bir dakika önce"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(Turkish), "yaklaşık bir dakika önce"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(Turkish), "2 dakika önce"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(Turkish), "yaklaşık bir saat önce"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(Turkish), "2 saat önce"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(Turkish), "23 saat önce"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(Turkish), "bir gün önce"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(Turkish), "2 gün önce"},
	{tBase, tBase.Add(30 * 24 * time.Hour), NoMax(Turkish), "bir ay önce"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(Turkish), "11 ay önce"},
	{tBase, tBase.Add(Year), NoMax(Turkish), "bir yıl önce"},
	{tBase, tBase.Add(547 * Day), NoMax(Turkish), "bir yıl önce"},
	{tBase, tBase.Add(548 * Day), NoMax(Turkish), "2 yıl önce"},
	{tBase, tBase.Add(10 * Year), NoMax(Turkish), "10 yıl önce"},

	// Max
	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(EnglishUS), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(EnglishUS, 90*time.Minute, ""), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(EnglishUS, 90*time.Minute, "2006-01-02"), "2013-08-30"},

	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(EnglishUK), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(EnglishUK, 90*time.Minute, ""), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(EnglishUK, 90*time.Minute, "02-01-2006"), "30-08-2013"},

	{tBase, tBase.Add(90 * time.Minute), WithMax(Portuguese, 90*time.Minute, "01-02-2006"), "08-30-2013"},

	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(German), "vor einer Stunde"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(German, 90*time.Minute, ""), "vor einer Stunde"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(German, 90*time.Minute, German.DefaultLayout), "30.08.2013"},

	// Turkish Max
	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(Turkish), "yaklaşık bir saat önce"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(Turkish, 90*time.Minute, ""), "yaklaşık bir saat önce"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(Turkish, 90*time.Minute, "02/01/2006"), "30/08/2013"},

	// Future
	{tBase.Add(24 * time.Hour), tBase, NoMax(EnglishUS), "in one day"},
	{tBase.Add(24 * time.Hour), tBase, NoMax(EnglishUK), "in one day"},
	{tBase.Add(24 * time.Hour), tBase, NoMax(Turkish), "bir gün içinde"},

	{tBase.Add(2 * Month), tBase, NoMax(Turkish), "2 ay içinde"},
	{tBase.Add(5 * time.Minute), tBase, NoMax(Turkish), "5 dakika içinde"},
	{tBase.Add(100 * time.Millisecond), tBase, NoMax(Turkish), "yaklaşık bir saniye içinde"},

	{tBase.Add(2 * Month), tBase, NoMax(Portuguese), "daqui a 2 meses"},
	{tBase.Add(24 * time.Hour), tBase, NoMax(Portuguese), "daqui a um dia"},
	{tBase.Add(5 * time.Minute), tBase, NoMax(Portuguese), "daqui a 5 minutos"},
	{tBase.Add(1 * time.Minute), tBase, NoMax(Portuguese), "daqui a um minuto"},
	{tBase.Add(100 * time.Millisecond), tBase, NoMax(Portuguese), "daqui a menos de um segundo"},

	{tBase.Add(2 * Month), tBase, NoMax(German), "in 2 Monaten"},
	{tBase.Add(24 * time.Hour), tBase, NoMax(German), "in einem Tag"},
	{tBase.Add(5 * time.Minute), tBase, NoMax(German), "in 5 Minuten"},
	{tBase.Add(1 * time.Minute), tBase, NoMax(German), "in einer Minute"},
	{tBase.Add(100 * time.Millisecond), tBase, NoMax(German), "in einer Sekunde"},
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
