package soh

import (
	"fmt"
	"net"
	"time"
)

type Xymon struct {
	Host     string
	Test     string
	Colour   string
	Label    string
	Message  string
	Epoch    time.Time
	Interval time.Duration
}

func (x *Xymon) FormatInterval() string {
	return fmt.Sprintf("%.0f", x.Interval.Seconds())
}

func (x *Xymon) FormatDate() string {
	return x.Epoch.Format("Mon Jan 02 03:04:05 2006")
}

func (x *Xymon) FormatLabel() string {
	msg := x.Label
	if x.Label != "" {
		msg = " " + msg
	}
	return msg
}

func (x *Xymon) FormatColour() string {
	colour := "green"
	if x.Colour != "" {
		colour = x.Colour
	}
	return colour
}

func (x *Xymon) FormatMessage() string {
	msg := x.Message
	if x.Message != "" {
		msg = "\n" + msg + "\n"
	}
	return msg
}

func (x *Xymon) Report() string {
	if x.Interval > 0 {
		return fmt.Sprintf("status+%s %s.%s %s %s%s\n%s", x.FormatInterval(), x.Host, x.Test, x.FormatColour(), x.FormatDate(), x.FormatLabel(), x.FormatMessage())
	} else {
		return fmt.Sprintf("status %s.%s %s %s%s\n%s", x.Host, x.Test, x.FormatColour(), x.FormatDate(), x.FormatLabel(), x.FormatMessage())
	}
}

func (x *Xymon) SendReport(server string) error {

	host := server

	_, _, err := net.SplitHostPort(server)
	if err != nil {
		host = net.JoinHostPort(server, "1984")
	}

	conn, err := net.Dial("udp", host)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := []byte(x.Report())
	_, err = conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
