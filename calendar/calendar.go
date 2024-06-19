package calendar

import (
	"birthdays/Db"
	"birthdays/models"
	"fmt"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
)

type Calendar struct {
	crn  *cron.Cron
	jobs []*func()
	st   *Db.Storage
}

func NewCalendar() *Calendar {
	return &Calendar{
		crn: cron.New(cron.WithSeconds()),
	}
}

func (c *Calendar) StartCal(st *Db.Storage) {

	c.CreateJobs(st)

	c.crn.Start()
}

func (c *Calendar) checkIfRecover(jobDate time.Time, lastRecord time.Time) bool {
	//utcJob := jobDate.UTC()
	//println(jobDate.String())
	if jobDate.After(lastRecord) && jobDate.Before(time.Now()) {

		return true
	}
	return false
}
func (c *Calendar) StopCal() {

}

func (c *Calendar) CreateJobs(st *Db.Storage) {
	var jobs []*models.CalendarJob
	jobs = st.SubscriptionsRows()
	lastRecord := st.GetLasCalRecord()

	var idToId map[int][]int
	idToId = make(map[int][]int)

	//
	//var mlt int
	//mlt = 12
	for _, i := range jobs {
		//mlt += -5
		//i.Date = testDate
		cronString := fmt.Sprintf("%d %d %d %d %d %d", i.Date.Second(), i.Date.Minute(), i.Date.Hour(), i.Date.Day(), i.Date.Month(), i.Date.Weekday())

		if c.checkIfRecover(i.Date, lastRecord) {

			cronString = fmt.Sprintf("%d %d %d %d %d %d", time.Now().Second()+2, time.Now().Minute(), time.Now().Hour(), time.Now().Day(), time.Now().Month(), time.Now().Weekday())
			i.Text = fmt.Sprintf("recover job for user %s", i.SubName)
		}

		id, _ := c.crn.AddFunc(cronString, func() {
			println(i.Text)
			st.SetLasCalRecord()
		})

		_, ok := idToId[i.Id]
		if !ok {
			idToId[i.Id] = []int{}
		}

		idToId[i.Id] = append(idToId[i.Id], int(id))
	}

	for i, v := range idToId {
		if err := st.UpdateIds(i, v); err != nil {
			println(err.Error())
		}
	}

}

func (c *Calendar) RemoveUserJobs(uid int, st *Db.Storage) {
	rm, err := st.FindJobsById(uid)
	if err != nil {
		println(err.Error())
	}
	jid, _ := strconv.Atoi(rm[0])
	c.crn.Remove(cron.EntryID(jid))
}

func job() cron.FuncJob {

	return func() {
		print("job")
	}
}
