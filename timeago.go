// Copyright 2013 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package timeago allows the formatting of time in terms of fuzzy timestamps.
//For example:
//	one minute ago
//	3 years ago
//	in 2 minutes
package timeago

import (
	"fmt"
	"strings"
	"time"
)

const (
	Day   time.Duration = time.Hour * 24
	Month time.Duration = Day * 30
	Year  time.Duration = Day * 365
)

type FormatPeriod struct {
	D    time.Duration
	One  string
	Many string
}

//Config allows the customization of timeago.
//You may configure string items (language, plurals, ...) and
//maximum allowed duration value for fuzzy formatting.
type Config struct {
	PastPrefix   string
	PastSuffix   string
	FuturePrefix string
	FutureSuffix string

	Periods []FormatPeriod

	Zero string
	Max  time.Duration //Maximum duration for using the special formatting.
	//DefaultLayout is used if delta is greater than the minimum of last period
	//in Periods and Max. It is the desired representation of the date 2nd of
	// January 2006.
	DefaultLayout string
}

//Predefined english configuration
var English = Config{
	PastPrefix:   "",
	PastSuffix:   " ago",
	FuturePrefix: "in ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		{D: time.Second, One: "about a second", Many: "a second%d seconds"},
		{D: time.Minute, One: "about a minute", Many: "%d minutes"},
		{D: time.Hour, One: "about an hour", Many: "%d hours"},
		{D: Day, One: "one day", Many: "%d days"},
		{D: Month, One: "one month", Many: "%d months"},
		{D: Year, One: "one year", Many: "%d years"},
	},

	Zero: "about a second",

	Max:           73 * time.Hour,
	DefaultLayout: "2006-01-02",
}

var Portuguese = Config{
	PastPrefix:   "há ",
	PastSuffix:   "",
	FuturePrefix: "daqui a ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		{D: time.Second, One: "um segundo", Many: "%d segundos"},
		{D: time.Minute, One: "um minuto", Many: "%d minutos"},
		{D: time.Hour, One: "uma hora", Many: "%d horas"},
		{D: Day, One: "um dia", Many: "%d dias"},
		{D: Month, One: "um mês", Many: "%d meses"},
		{D: Year, One: "um ano", Many: "%d anos"},
	},

	Zero: "menos de um segundo",

	Max:           73 * time.Hour,
	DefaultLayout: "02-01-2006",
}

var Chinese = Config{
	PastPrefix:   "",
	PastSuffix:   "前",
	FuturePrefix: "于 ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		{D: time.Second, One: "1 秒", Many: "%d 秒"},
		{D: time.Minute, One: "1 分钟", Many: "%d 分钟"},
		{D: time.Hour, One: "1 小时", Many: "%d 小时"},
		{D: Day, One: "1 天", Many: "%d 天"},
		{D: Month, One: "1 月", Many: "%d 月"},
		{D: Year, One: "1 年", Many: "%d 年"},
	},

	Zero: "1 秒",

	Max:           73 * time.Hour,
	DefaultLayout: "2006-01-02",
}

//Predefined french configuration
var French = Config{
	PastPrefix:   "il y a ",
	PastSuffix:   "",
	FuturePrefix: "dans ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		{D: time.Second, One: "environ une seconde", Many: "moins d'une minute"},
		{D: time.Minute, One: "environ une minute", Many: "%d minutes"},
		{D: time.Hour, One: "environ une heure", Many: "%d heures"},
		{D: Day, One: "un jour", Many: "%d jours"},
		{D: Month, One: "un mois", Many: "%d mois"},
		{D: Year, One: "un an", Many: "%d ans"},
	},

	Zero: "environ une seconde",

	Max:           73 * time.Hour,
	DefaultLayout: "02/01/2006",
}

//Predefined german configuration
var German = Config{
	PastPrefix:   "vor ",
	PastSuffix:   "",
	FuturePrefix: "in ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		{D: time.Second, One: "einer Sekunde", Many: "%d Sekunden"},
		{D: time.Minute, One: "einer Minute", Many: "%d Minuten"},
		{D: time.Hour, One: "einer Stunde", Many: "%d Stunden"},
		{D: Day, One: "einem Tag", Many: "%d Tagen"},
		{D: Month, One: "einem Monat", Many: "%d Monaten"},
		{D: Year, One: "einem Jahr", Many: "%d Jahren"},
	},

	Zero: "einer Sekunde",

	Max:           73 * time.Hour,
	DefaultLayout: "02.01.2006",
}

//Predefined turkish configuration
var Turkish = Config{
	PastPrefix:   "",
	PastSuffix:   " önce",
	FuturePrefix: "",
	FutureSuffix: " içinde",

	Periods: []FormatPeriod{
		{D: time.Second, One: "yaklaşık bir saniye", Many: "%d saniye"},
		{D: time.Minute, One: "yaklaşık bir dakika", Many: "%d dakika"},
		{D: time.Hour, One: "yaklaşık bir saat", Many: "%d saat"},
		{D: Day, One: "bir gün", Many: "%d gün"},
		{D: Month, One: "bir ay", Many: "%d ay"},
		{D: Year, One: "bir yıl", Many: "%d yıl"},
	},

	Zero: "yaklaşık bir saniye",

	Max:           73 * time.Hour,
	DefaultLayout: "02/01/2006",
}

// Korean support
var Korean = Config{
	PastPrefix:   "",
	PastSuffix:   " 전",
	FuturePrefix: "",
	FutureSuffix: " 후",

	Periods: []FormatPeriod{
		{D: time.Second, One: "약 1초", Many: "%d초"},
		{D: time.Minute, One: "약 1분", Many: "%d분"},
		{D: time.Hour, One: "약 한시간", Many: "%d시간"},
		{D: Day, One: "하루", Many: "%d일"},
		{D: Month, One: "1개월", Many: "%d개월"},
		{D: Year, One: "1년", Many: "%d년"},
	},

	Zero: "방금",

	Max:           10 * 365 * 24 * time.Hour,
	DefaultLayout: "2006-01-02",
}

//Format returns a textual representation of the time value formatted according to the layout
//defined in the Config. The time is compared to time.Now() and is then formatted as a fuzzy
//timestamp (eg. "4 days ago")
func (cfg Config) Format(t time.Time) string {
	return cfg.FormatReference(t, time.Now())
}

//FormatReference is the same as Format, but the reference has to be defined by the caller
func (cfg Config) FormatReference(t time.Time, reference time.Time) string {

	d := reference.Sub(t)

	if (d >= 0 && d >= cfg.Max) || (d < 0 && -d >= cfg.Max) {
		return t.Format(cfg.DefaultLayout)
	}

	return cfg.FormatRelativeDuration(d)
}

//FormatRelativeDuration is the same as Format, but for time.Duration.
//Config.Max is not used in this function, as there is no other alternative.
func (cfg Config) FormatRelativeDuration(d time.Duration) string {

	isPast := d >= 0

	if d < 0 {
		d = -d
	}

	s, _ := cfg.getTimeText(d, true)

	if isPast {
		return strings.Join([]string{cfg.PastPrefix, s, cfg.PastSuffix}, "")
	}
	return strings.Join([]string{cfg.FuturePrefix, s, cfg.FutureSuffix}, "")

}

//Round the duration d in terms of step.
func round(d time.Duration, step time.Duration, roundCloser bool) time.Duration {

	if roundCloser {
		return time.Duration(float64(d)/float64(step) + 0.5)
	}

	return time.Duration(float64(d) / float64(step))
}

//Count the number of parameters in a format string
func nbParamInFormat(f string) int {
	return strings.Count(f, "%") - 2*strings.Count(f, "%%")
}

//Convert a duration to a text, based on the current config
func (cfg Config) getTimeText(d time.Duration, roundCloser bool) (string, time.Duration) {
	if len(cfg.Periods) == 0 || d < cfg.Periods[0].D {
		return cfg.Zero, 0
	}

	for i, p := range cfg.Periods {

		next := p.D
		if i+1 < len(cfg.Periods) {
			next = cfg.Periods[i+1].D
		}

		if i+1 == len(cfg.Periods) || d < next {

			r := round(d, p.D, roundCloser)

			if next != p.D && r == round(next, p.D, roundCloser) {
				continue
			}

			if r == 0 {
				return "", d
			}

			layout := p.Many
			if r == 1 {
				layout = p.One
			}

			if nbParamInFormat(layout) == 0 {
				return layout, d - r*p.D
			}

			return fmt.Sprintf(layout, r), d - r*p.D
		}
	}

	return d.String(), 0
}

//NoMax creates an new config without a maximum
func NoMax(cfg Config) Config {
	return WithMax(cfg, 9223372036854775807, time.RFC3339)
}

//WithMax creates an new config with special formatting limited to durations less than max.
//Values greater than max will be formatted by the standard time package using the defaultLayout.
func WithMax(cfg Config, max time.Duration, defaultLayout string) Config {
	n := cfg
	n.Max = max
	n.DefaultLayout = defaultLayout
	return n
}
