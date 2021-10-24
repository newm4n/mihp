package helper

import (
	"fmt"
	"github.com/newm4n/mihp/pkg/errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// NewCronStruct create a cron time checker from a string.
// Cron field as follows
//     * * * * * * *
//     | | | | | | +-- Year
//     | | | | | +---- Month (1-12)
//     | | | | +------ Day of Week (0-6)
//     | | | +-------- Day of Month (1-12)
//     | | +---------- Hour (0-23)
//     | +------------ Minute (0-59)
//     +-------------- Second (0-59)
//
// Each field can support multiple entry such as:
//     -  1,4,6 -->  either 1,4, or 5
//     -  1-5   -->  any number from 1 inclusive to 5 inclusive.
//     -  5-    -->  any number from 5 and beyond
//     -  -5    -->  any number up to 5 inclusive
//     -  */5   -->  any number of mod 5
//     -  *     -->  any number
//
// Within the field, you can mix and match. e.g:
//     -  3,6,10-15,30-
//
// So the entire cron string could be as complex as:
//     - 1,4,6,12-23    1,4,6,12-23 1,4,6,12-23  1,4,6,12-23 1,4,6,12-23 1,4,6,12-23 1,4,6,12-23,40-,-50
//
func NewCronStruct(cron string) (*CronStruct, error) {
	numOnly := regexp.MustCompile(`[ \n\r\t]+`)
	nCron := numOnly.ReplaceAllString(cron, " ")

	tok := strings.Split(nCron, " ")
	if len(tok) != 7 {
		return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, cron)
	}

	ret := &CronStruct{
		syntax:            cron,
		secondInterval:    nil,
		minuteInterval:    nil,
		hourInterval:      nil,
		dayInterval:       nil,
		dayOfWeekInterval: nil,
		monthInterval:     nil,
		yearInterval:      nil,
	}

	var err error
	ret.secondInterval, err = StringToInterval(tok[0])
	if err != nil {
		return nil, err
	}
	ret.minuteInterval, err = StringToInterval(tok[1])
	if err != nil {
		return nil, err
	}
	ret.hourInterval, err = StringToInterval(tok[2])
	if err != nil {
		return nil, err
	}
	ret.dayInterval, err = StringToInterval(tok[3])
	if err != nil {
		return nil, err
	}
	ret.dayOfWeekInterval, err = StringToInterval(tok[4])
	if err != nil {
		return nil, err
	}
	ret.monthInterval, err = StringToInterval(tok[5])
	if err != nil {
		return nil, err
	}
	ret.yearInterval, err = StringToInterval(tok[6])
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type CronStruct struct {
	syntax            string
	secondInterval    *Interval // 0-59
	minuteInterval    *Interval // 0-59
	hourInterval      *Interval // 0-23
	dayInterval       *Interval // 1-31
	dayOfWeekInterval *Interval // 0-6
	monthInterval     *Interval // 1-12
	yearInterval      *Interval
}

func (c *CronStruct) IsIn(t time.Time) bool {
	if !c.secondInterval.IsIn(t.Second()) {
		logrus.Tracef("Second %d not in cron %s, %s", t.Second(), c.syntax, c.secondInterval.String())
		return false
	}
	if !c.minuteInterval.IsIn(t.Minute()) {
		logrus.Tracef("Minute %d not in cron %s", t.Minute(), c.syntax)
		return false
	}
	if !c.hourInterval.IsIn(t.Hour()) {
		logrus.Tracef("Hour %d not in cron %s", t.Hour(), c.syntax)
		return false
	}
	if !c.dayInterval.IsIn(t.Day()) {
		logrus.Tracef("Day %d not in cron %s", t.Day(), c.syntax)
		return false
	}
	if !c.monthInterval.IsIn(int(t.Month())) {
		logrus.Tracef("Month %d not in cron %s", int(t.Month()), c.syntax)
		return false
	}
	if !c.yearInterval.IsIn(t.Year()) {
		logrus.Tracef("Year %d not in cron %s", t.Year(), c.syntax)
		return false
	}
	if !c.dayOfWeekInterval.IsIn(int(t.Weekday())) {
		logrus.Tracef("Week day %d not in cron %s", int(t.Weekday()), c.syntax)
		return false
	}
	return true
}
