package cron

import (
	"context"
	"fmt"
	"github.com/newm4n/mihp/pkg/errors"
	"github.com/newm4n/mihp/pkg/helper"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// NewSchedule create a cron time checker from a string.
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
func NewSchedule(cron string) (*Schedule, error) {
	numOnly := regexp.MustCompile(`[ \n\r\t]+`)
	nCron := numOnly.ReplaceAllString(cron, " ")

	tok := strings.Split(nCron, " ")
	if len(tok) != 7 {
		return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, cron)
	}

	ret := &Schedule{
		cronSyntax:        cron,
		secondInterval:    nil,
		minuteInterval:    nil,
		hourInterval:      nil,
		dayInterval:       nil,
		dayOfWeekInterval: nil,
		monthInterval:     nil,
		yearInterval:      nil,
	}

	var err error
	ret.secondInterval, err = helper.StringToInterval(tok[0])
	if err != nil {
		return nil, err
	}
	ret.minuteInterval, err = helper.StringToInterval(tok[1])
	if err != nil {
		return nil, err
	}
	ret.hourInterval, err = helper.StringToInterval(tok[2])
	if err != nil {
		return nil, err
	}
	ret.dayInterval, err = helper.StringToInterval(tok[3])
	if err != nil {
		return nil, err
	}
	ret.dayOfWeekInterval, err = helper.StringToInterval(tok[4])
	if err != nil {
		return nil, err
	}
	ret.monthInterval, err = helper.StringToInterval(tok[5])
	if err != nil {
		return nil, err
	}
	ret.yearInterval, err = helper.StringToInterval(tok[6])
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Schedule define a specific events in time, this specific point in time can be recurrence every minutes, hours, days, weeks, months and time.
type Schedule struct {
	cronSyntax        string
	secondInterval    *helper.Interval // 0-59
	minuteInterval    *helper.Interval // 0-59
	hourInterval      *helper.Interval // 0-23
	dayInterval       *helper.Interval // 1-31
	dayOfWeekInterval *helper.Interval // 0-6
	monthInterval     *helper.Interval // 1-12
	yearInterval      *helper.Interval
}

// IsIn will check if the give time argument is matching to the time specified within the schedule.
func (c *Schedule) IsIn(t time.Time) bool {
	if !c.secondInterval.IsIn(t.Second()) {
		return false
	}
	if !c.minuteInterval.IsIn(t.Minute()) {
		return false
	}
	if !c.hourInterval.IsIn(t.Hour()) {
		return false
	}
	if !c.dayInterval.IsIn(t.Day()) {
		return false
	}
	if !c.monthInterval.IsIn(int(t.Month())) {
		return false
	}
	if !c.yearInterval.IsIn(t.Year()) {
		return false
	}
	if !c.dayOfWeekInterval.IsIn(int(t.Weekday())) {
		return false
	}

	logrus.Tracef("Time %s is in cron %s", t, c.cronSyntax)
	return true
}

var (
	jobs       = make(map[string]*Job)
	cronTicker *time.Ticker
	alive      = false
	stopChan   chan bool
	cronLogger = logrus.WithField("module", "Cron")
)

// Job specifies a function to be execution in a specific Schedule
type Job struct {
	Cron     *Schedule
	JobFunc  func(ctx context.Context)
	Deadline time.Duration
}

// AddJob adds a specific job into this scheduler engine, so their function can be invoked on the specified schedule.
func AddJob(jobId string, job *Job) {
	jobs[jobId] = job
}

// RemoveJob will remove an existing Job with specified id.
func RemoveJob(jobId string) {
	delete(jobs, jobId)
}

func tickerEvent(t time.Time) {
	cronLogger.Tracef("Scheduler ticks %s", t)
	for n, j := range jobs {
		if j.Cron.IsIn(t) {
			cronLogger.Debugf("executing job %s with cron cronSyntax %s at %s. Deadline for %s", n, j.Cron.cronSyntax, j.Deadline, t)
			ctx, _ := context.WithTimeout(context.Background(), j.Deadline)
			go j.JobFunc(ctx)
		}
	}
}

// Start the scheduler server
func Start() {
	cronLogger.Info("Starting module")
	stopChan = make(chan bool)
	if !alive {
		cronTicker = time.NewTicker(1 * time.Second)
		alive = true
		go func() {
			for {
				select {
				case <-stopChan:
					return
				case t := <-cronTicker.C:
					tickerEvent(t)
				}
			}
		}()
	}
}

// Stop the scheduler server
func Stop() {
	cronLogger.Info("Stopping module")
	if alive {
		alive = false
		cronTicker.Stop()
		stopChan <- true
	}
}
