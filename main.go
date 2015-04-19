package main

import (
	"errors"
	"flag"
	"log"
	"os"
	// _ "github.com/samuel/go-librato/librato"
	"path/filepath"
	"time"
)

// xymon bounds
type Bounds struct {
	PreEvent       float64
	PostEvent      float64
	SamplingRate   float64
	ClockWarning   float64
	ClockError     float64
	StorageWarning float64
	StorageError   float64
	VoltageWarning float64
	VoltageError   float64
}

// xymon settings
type Xymon struct {
	Host   string
	Server string
	Valid  time.Duration
}

// runtime settings
var (
	verbose bool
	dryrun  bool
	unlink  bool
	match   string
	recent  time.Duration

	bounds Bounds
	xymon  Xymon
)

func Process(path string) error {
	// decode cusp file ...
	c, err := DecodeCusp(path)
	if err != nil {
		return err
	}
	// recover checkin time
	t, err := c.CheckIn()
	if err != nil {
		return err
	}
	// too late?
	if t.Before(time.Now().UTC().Add(-recent)) {
		return errors.New("file too old")
	}
	// recover xymon metrics
	metrics := c.Xymon(xymon.Host, xymon.Valid, bounds)
	for _, m := range metrics {
		if verbose {
			log.Print(m.Report())
		}
		m.SendReport(xymon.Server)
	}

	return nil
}

func Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("check: %s\n", path)
	}

	m, err := filepath.Match(match, filepath.Base(path))
	if err != nil {
		return err
	}
	if !m {
		if verbose {
			log.Printf("skip: %s\n", path)
		}
		return nil
	}

	if verbose {
		log.Printf("process: %s\n", path)
	}
	err = Process(path)
	if err != nil {
		log.Println(err)
	}

	if unlink {
		if verbose {
			log.Printf("unlink: %s\n", path)
		}
		if !dryrun {
			err = os.Remove(path)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func init() {

	bounds = Bounds{
		PreEvent:       40.0,
		PostEvent:      60.0,
		SamplingRate:   200.0,
		ClockWarning:   50.0,
		ClockError:     30.0,
		StorageWarning: 50,
		StorageError:   10,
		VoltageWarning: 11.5,
		VoltageError:   11.0,
	}

	xymon = Xymon{
		Host:   "unknown",
		Server: "localhost:1984",
		Valid:  2 * time.Hour,
	}
}

func main() {

	flag.BoolVar(&verbose, "verbose", false, "make noise")
	flag.BoolVar(&dryrun, "dry-run", false, "don't actually send the messages")
	flag.BoolVar(&unlink, "unlink", false, "remove file after processing")
	flag.StringVar(&match, "match", "Report_*.xml", "provide file matching template")
	flag.DurationVar(&recent, "recent", 24*time.Hour, "checkin file minimum age")

	flag.StringVar(&xymon.Server, "xymon", xymon.Server, "provide a xymon server address")
	flag.DurationVar(&xymon.Valid, "valid", xymon.Valid, "checkin validity")
	flag.StringVar(&xymon.Host, "host", xymon.Host, "provide a xymon host name")

	flag.Float64Var(&bounds.PreEvent, "pre", bounds.PreEvent, "expected pre-event time")
	flag.Float64Var(&bounds.PostEvent, "post", bounds.PostEvent, "expected post-event time")
	flag.Float64Var(&bounds.SamplingRate, "rate", bounds.SamplingRate, "expected sampling rate")
	flag.Float64Var(&bounds.StorageWarning, "storage-warning", bounds.StorageWarning, "provide a storage warning level in Mbytes")
	flag.Float64Var(&bounds.StorageError, "storage-error", bounds.StorageWarning, "provide a storage error level in Mbytes")
	flag.Float64Var(&bounds.VoltageWarning, "voltage-warning", bounds.VoltageWarning, "provide a warning voltage level")
	flag.Float64Var(&bounds.VoltageError, "voltage-error", bounds.VoltageError, "provide a voltage error level")
	flag.Float64Var(&bounds.ClockWarning, "clock-warning", bounds.ClockWarning, "provide a clock warning quality")
	flag.Float64Var(&bounds.ClockError, "clock-error", bounds.ClockError, "provide a clock error quality")

	flag.Parse()

	// run through the given directories ...
	for _, d := range flag.Args() {
		log.Printf("walking directory: %s\n", d)

		err := filepath.Walk(d, Walk)
		if err != nil {
			log.Printf("[%s] %s\n", d, err)

		}
	}

}
