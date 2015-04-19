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
)

func Process(bounds Bounds, xymon Xymon, f string) error {
	// decode cusp file ...
	c, err := DecodeCusp(f)
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
		m.SendReport(xymon.Server)
	}

	return nil
}

func main() {

	b := Bounds{
		PreEvent:       40.0,
		PostEvent:      60.0,
		SamplingRate:   200.0,
		ClockWarning:   50.0,
		ClockError:     30.0,
		StorageWarning: 10,
		StorageError:   10,
		VoltageWarning: 11.5,
		VoltageError:   11.0,
	}

	x := Xymon{
		Host:   "unknown",
		Server: "localhost:1984",
		Valid:  2 * time.Hour,
	}

	flag.BoolVar(&verbose, "verbose", false, "make noise")
	flag.BoolVar(&dryrun, "dry-run", false, "don't actually send the messages")
	flag.BoolVar(&unlink, "unlink", false, "remove file after processing")
	flag.StringVar(&match, "match", "*.xml", "provide file matching template")
	flag.DurationVar(&recent, "recent", 24*time.Hour, "checkin file minimum age")

	flag.StringVar(&x.Server, "xymon", x.Server, "provide a xymon server address")
	flag.DurationVar(&x.Valid, "valid", x.Valid, "checkin validity")
	flag.StringVar(&x.Host, "host", x.Host, "provide a xymon host name")

	flag.Float64Var(&b.PreEvent, "pre", b.PreEvent, "expected pre-event time")
	flag.Float64Var(&b.PostEvent, "post", b.PostEvent, "expected post-event time")
	flag.Float64Var(&b.SamplingRate, "rate", b.SamplingRate, "expected sampling rate")
	flag.Float64Var(&b.StorageWarning, "storage-warning", b.StorageWarning, "provide a storage warning level in Mbytes")
	flag.Float64Var(&b.StorageError, "storage-error", b.StorageWarning, "provide a storage error level in Mbytes")
	flag.Float64Var(&b.VoltageWarning, "voltage-warning", b.VoltageWarning, "provide a warning voltage level")
	flag.Float64Var(&b.VoltageError, "voltage-error", b.VoltageError, "provide a voltage error level")
	flag.Float64Var(&b.ClockWarning, "clock-warning", b.ClockWarning, "provide a clock warning quality")
	flag.Float64Var(&b.ClockError, "clock-error", b.ClockError, "provide a clock error quality")

	flag.Parse()

	// run through the given files
	for _, f := range flag.Args() {
		log.Printf("checking file: %s\n", f)

		// check that the filenames match ...
		_, err := filepath.Match(match, filepath.Base(f))
		if err != nil {
			log.Printf("skipping unmatched file: %s\n", f)
			continue
		}

		// check that they exist ...
		log.Printf("processing file: %s\n", f)
		if _, err := os.Stat(f); os.IsNotExist(err) {
			log.Printf("file missing: %s, skipping\n", f)
			continue
		}

		if err := Process(b, x, f); err != nil {
			log.Printf("processing error: %s [%s]\n", f, err)
		}

		if unlink {
			log.Printf("unlinking file: %s\n", f)
			if !dryrun {
				err = os.Remove(f)
				if err != nil {
					log.Println(err)
				}
			}
		}

	}

}
