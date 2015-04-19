package main

import (
	"fmt"
	"github.com/ozym/soh"
	"strconv"
	"time"
)

func (c *Cusp) ClockStatus() (float64, string, error) {
	var msg string

	if c.DASTimeLock != "" {
		msg += fmt.Sprintf("Lock: %s\n", c.DASTimeLock)
	}
	if c.GPSState != "" {
		msg += fmt.Sprintf("State: %s\n", c.GPSState)
	}
	if c.GPSLossPeriod != "" {
		msg += fmt.Sprintf("Loss: %s\n", c.GPSLossPeriod)
	}

	m, err := strconv.ParseFloat(c.TimeSuccessA, 32)
	if err == nil {
		msg += fmt.Sprintf("Clock: %.0f\n", m)
	}

	return m, msg, nil
}

func (c *Cusp) Okay(bounds Bounds) (bool, string, error) {
	var msg string = ""

	_, err := c.CheckIn()
	if err != nil {
		return false, msg, err
	}
	rate, err := c.SamplingRate()
	if err != nil {
		return false, msg, err
	}
	pre, err := c.PreEvent()
	if err != nil {
		return false, msg, err
	}
	post, err := c.PostEvent()
	if err != nil {
		return false, msg, err
	}

	if rate != bounds.SamplingRate {
		msg += fmt.Sprintf("Invalid Sampling Rate: [%g != %g]\n", rate, bounds.SamplingRate)
	}
	if pre != bounds.PreEvent {
		msg += fmt.Sprintf("Invalid Pre-Event Time: [%g != %g]\n", pre, bounds.PreEvent)
	}
	if post != bounds.PostEvent {
		msg += fmt.Sprintf("Invalid Post-Event Time: [%g != %g]\n", post, bounds.PostEvent)
	}
	if msg != "" {
		return false, msg, nil
	}

	return true, msg, nil
}

func (c *Cusp) Xymon(host string, valid time.Duration, bounds Bounds) []soh.Xymon {
	var x []soh.Xymon

	e, err := c.CheckIn()
	if err != nil {
		return x
	}

	o, m, err := c.Okay(bounds)
	if err == nil {
		var q string

		if o {
			q = "green"
		} else {
			q = "red"
		}

		if o {
			m += fmt.Sprintf("System info: okay\n\nLast Checkin: %s\n", e.Local().Format(time.UnixDate))
		} else {
			m += fmt.Sprintf("System info: fault\n\nLast Checkin: %s\n", e.Local().Format(time.UnixDate))
		}

		x = append(x, soh.Xymon{
			Host:     host,
			Epoch:    e,
			Interval: valid,
			Colour:   q,
			Test:     "checkin",
			Label:    fmt.Sprintf("valid checkin registered %.1f hours ago", time.Now().Sub(e).Hours()),
			Message:  m,
		})
	}

	p, n, err := c.ClockStatus()
	if err == nil {
		q := "green"
		if p < bounds.ClockWarning {
			q = "yellow"
		}
		if p < bounds.ClockError {
			q = "red"
		}
		x = append(x, soh.Xymon{
			Host:     host,
			Epoch:    e,
			Interval: valid,
			Colour:   q,
			Test:     "clock",
			Label:    fmt.Sprintf("quality=%.1f", p),
			Message:  n,
		})
	}
	v, err := c.Power()
	if err == nil {
		q := "green"
		if v < bounds.VoltageWarning {
			q = "yellow"
		}
		if v < bounds.VoltageError {
			q = "red"
		}
		x = append(x, soh.Xymon{
			Host:     host,
			Epoch:    e,
			Interval: valid,
			Colour:   q,
			Test:     "bat",
			Label:    fmt.Sprintf("voltage=%.1f", v),
			Message:  fmt.Sprintf("Current voltage: %.2f volts\n", v),
		})
	}

	t, err := c.Temperature()
	if err == nil {
		x = append(x, soh.Xymon{
			Host:     host,
			Epoch:    e,
			Colour:   "green",
			Interval: valid,
			Test:     "temp",
			Label:    fmt.Sprintf("temp=%.1f", t),
			Message:  fmt.Sprintf("Current temperature: %.1f\n", t),
		})
	}

	s, err := c.Storage()
	if err == nil {
		q := "green"
		if v < (bounds.StorageWarning * 1024) {
			q = "yellow"
		}
		if v < (bounds.StorageError * 1024) {
			q = "red"
		}
		x = append(x, soh.Xymon{
			Host:     host,
			Epoch:    e,
			Interval: valid,
			Colour:   q,
			Test:     "storage",
			Label:    fmt.Sprintf("storage=%d", s),
			Message:  fmt.Sprintf("Current storage: %d kbytes\n", s),
		})
	}

	return x
}
