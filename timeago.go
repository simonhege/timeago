// Copyright 2013 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// timeago allows the formatting of time in terms of fuzzy timestamps.
// For example:
//		one minute ago
//		3 years ago
//		in 2 minutes
package timeago

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goinggo/beego-mgo/go-i18n/i18n/language"
	"github.com/goinggo/beego-mgo/go-i18n/i18n/plural"
)

const (
	Day   time.Duration = time.Hour * 24
	Month time.Duration = Day * 30
	Year  time.Duration = Day * 365
)

type FormatPeriod struct {
	D     time.Duration
	One   string
	Few   string
	Many  string
	Other string
}

// Config allows the customization of timeago.
// You may configure string items (language, plurals, ...) and
// maximum allowed duration value for fuzzy formatting.
type Config struct {
	LocaleID     string
	PastPrefix   string
	PastSuffix   string
	FuturePrefix string
	FutureSuffix string

	Periods []FormatPeriod

	Zero          string
	Max           time.Duration // Maximum duration for using the special formatting.
	DefaultLayout string        // Layout to use if delta is greater than the minimum of last period in Periods and Max
}

var (
	// https://github.com/goinggo/beego-mgo/tree/master/go-i18n/i18n
	// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html#ru
	pluralRussian = &language.Language{
		ID:               "ru",
		PluralCategories: newSet(plural.One, plural.Few, plural.Many, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			mod10 := ops.I % 10
			mod100 := ops.I % 100
			if ops.V == 0 && mod10 == 1 && mod100 != 11 {
				return plural.One
			}
			if ops.V == 0 && (mod10 >= 2 && mod10 <= 4) && (mod100 < 12 || mod100 > 14) {
				return plural.Few
			}
			if ops.V == 0 && (mod10 == 0 || (mod10 >= 5 && mod10 <= 9) || (mod100 >= 11 && mod100 <= 14)) {
				return plural.Many
			}
			return plural.Other
		},
	}
)

// Predefined english configuration
var English = Config{
	LocaleID:     "en",
	PastPrefix:   "",
	PastSuffix:   " ago",
	FuturePrefix: "in ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		FormatPeriod{
			D:     time.Second,
			One:   "about a second",
			Other: "%d seconds",
		},
		FormatPeriod{
			D:     time.Minute,
			One:   "about a minute",
			Other: "%d minutes",
		},
		FormatPeriod{
			D:     time.Hour,
			One:   "about an hour",
			Other: "%d hours",
		},
		FormatPeriod{
			D:     Day,
			One:   "one day",
			Other: "%d days",
		},
		FormatPeriod{
			D:     Month,
			One:   "one month",
			Other: "%d months",
		},
		FormatPeriod{
			D:     Year,
			One:   "one year",
			Other: "%d years",
		},
	},

	Zero: "about a second",

	Max:           73 * time.Hour,
	DefaultLayout: "2006-01-02",
}

var Chinese = Config{
	LocaleID:     "zh",
	PastPrefix:   "",
	PastSuffix:   "前",
	FuturePrefix: "于 ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		FormatPeriod{
			D:     time.Second,
			One:   "1 秒",
			Other: "%d 秒",
		},
		FormatPeriod{
			D:     time.Minute,
			One:   "1 分钟",
			Other: "%d 分钟",
		},
		FormatPeriod{
			D:     time.Hour,
			One:   "1 小时",
			Other: "%d 小时",
		},
		FormatPeriod{
			D:     Day,
			One:   "1 天",
			Other: "%d 天",
		},
		FormatPeriod{
			D:     Month,
			One:   "1 月",
			Other: "%d 月",
		},
		FormatPeriod{
			D:     Year,
			One:   "1 年",
			Other: "%d 年",
		},
	},

	Zero: "1 秒",

	Max:           73 * time.Hour,
	DefaultLayout: "2006-01-02",
}

// Predefined french configuration
var French = Config{
	LocaleID:     "fr",
	PastPrefix:   "il y a ",
	PastSuffix:   "",
	FuturePrefix: "dans ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		FormatPeriod{
			D:     time.Second,
			One:   "environ une seconde",
			Other: "moins d'une minute",
		},
		FormatPeriod{
			D:     time.Minute,
			One:   "environ une minute",
			Other: "%d minutes",
		},
		FormatPeriod{
			D:     time.Hour,
			One:   "environ une heure",
			Other: "%d heures",
		},
		FormatPeriod{
			D:     Day,
			One:   "un jour",
			Other: "%d jours",
		},
		FormatPeriod{
			D:     Month,
			One:   "un mois",
			Other: "%d mois",
		},
		FormatPeriod{
			D:     Year,
			One:   "un an",
			Other: "%d ans",
		},
	},

	Zero: "environ une seconde",

	Max:           73 * time.Hour,
	DefaultLayout: "02/01/2006",
}

var Russian = Config{
	LocaleID:     "ru",
	PastPrefix:   "",
	PastSuffix:   " назад",
	FuturePrefix: "через ",
	FutureSuffix: "",

	Periods: []FormatPeriod{
		FormatPeriod{
			D:     time.Second,
			One:   "секунду",
			Few:   "%d секунды",
			Many:  "%d секунд",
			Other: "%d секунды",
		},
		FormatPeriod{
			D:     time.Minute,
			One:   "около минуты",
			Few:   "%d минуты",
			Many:  "%d минут",
			Other: "%d минуты",
		},
		FormatPeriod{
			D:     time.Hour,
			One:   "около часа",
			Few:   "%d часа",
			Many:  "%d часов",
			Other: "%d часа",
		},
		FormatPeriod{
			D:     Day,
			One:   "день",
			Few:   "%d дня",
			Many:  "%d дней",
			Other: "%d дня",
		},
		FormatPeriod{
			D:     Month,
			One:   "месяц",
			Few:   "%d месяца",
			Many:  "%d месяцев",
			Other: "%d месяца",
		},
		FormatPeriod{
			D:     Year,
			One:   "год",
			Few:   "%d года",
			Many:  "%d лет",
			Other: "%d года",
		},
	},

	Zero: "около секунды",

	Max:           73 * time.Hour,
	DefaultLayout: "2006-01-02",
}

// Register locales
func init() {
	language.Register(pluralRussian)
}

// newSet plural categories
// https://github.com/goinggo/beego-mgo/blob/master/go-i18n/i18n/language/language.go#L240
func newSet(pluralCategories ...plural.Category) map[plural.Category]struct{} {
	set := make(map[plural.Category]struct{}, len(pluralCategories))
	for _, pc := range pluralCategories {
		set[pc] = struct{}{}
	}
	return set
}

// Format returns a textual representation of the time value formatted according to the layout
// defined in the Config. The time is compared to time.Now() and is then formatted as a fuzzy
// timestamp (eg. "4 days ago")
func (cfg Config) Format(t time.Time) string {
	return cfg.FormatReference(t, time.Now())
}

// FormatReference is the same as Format, but the reference has to be defined by the caller
func (cfg Config) FormatReference(t time.Time, reference time.Time) string {
	d := reference.Sub(t)

	if (d >= 0 && d >= cfg.Max) || (d < 0 && -d >= cfg.Max) {
		return t.Format(cfg.DefaultLayout)
	}

	return cfg.FormatRelativeDuration(d)
}

// FormatRelativeDuration is the same as Format, but for time.Duration.
// Config.Max is not used in this function, as there is no other alternative.
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

// Round the duration d in terms of step.
func round(d time.Duration, step time.Duration, roundCloser bool) time.Duration {
	if roundCloser {
		return time.Duration(float64(d)/float64(step) + 0.5)
	}
	return time.Duration(float64(d) / float64(step))
}

// Count the number of parameters in a format string
func nbParamInFormat(f string) int {
	return strings.Count(f, "%") - 2*strings.Count(f, "%%")
}

// durationNumber Get number from time.Duration
func durationNumber(d time.Duration) (int, error) {
	re := regexp.MustCompile("[^0-9]")
	return strconv.Atoi(re.ReplaceAllString(d.String(), ""))
}

// Convert a duration to a text, based on the current config
func (cfg Config) getTimeText(d time.Duration, roundCloser bool) (string, time.Duration) {
	var layout string

	if len(cfg.Periods) == 0 || d < cfg.Periods[0].D {
		return cfg.Zero, 0
	}

	lang := language.LanguageWithID(cfg.LocaleID)

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

			pluralNumber, _ := durationNumber(r)
			pluralCategory, _ := lang.PluralCategory(pluralNumber)
			switch pluralCategory {
			case plural.One:
				layout = p.One
			case plural.Few:
				layout = p.Few
			case plural.Many:
				layout = p.Many
			case plural.Invalid:
			default:
				layout = p.Other
			}

			if nbParamInFormat(layout) == 0 {
				return layout, d - r*p.D
			}

			return fmt.Sprintf(layout, r), d - r*p.D
		}
	}

	return d.String(), 0
}

// NoMax creates an new config without a maximum
func NoMax(cfg Config) Config {
	return WithMax(cfg, 9223372036854775807, time.RFC3339)
}

// WithMax creates an new config with special formatting limited to durations less than max.
// Values greater than max will be formatted by the standard time package using the defaultLayout.
func WithMax(cfg Config, max time.Duration, defaultLayout string) Config {
	n := cfg
	n.Max = max
	n.DefaultLayout = defaultLayout
	return n
}
