package util

import (
	"errors"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

var zoneDirs = []string{
	// Update path according to your OS
	"/usr/share/zoneinfo/",
	"/usr/share/lib/zoneinfo/",
	"/usr/lib/locale/TZ/",
}

var zoneDir string
var timezones []string
var mx sync.Mutex
var dateTimeLayout = "2006-01-02 15:04:05"
var dateTimeShortLayout = "2006-01-02 15"

// GetTimeZones - will return a list with all possible timezones
func GetTimeZones() []string {
	mx.Lock()
	defer func() {
		mx.Unlock()
	}()

	for _, zoneDir = range zoneDirs {
		readDirRecursive("")
	}

	return timezones
}

// GetUserTimeZone - goes over all timezones and trying determine the target one.
func GetUserTimeZone(dateTime string) (string, error) {
	var loc *time.Location

	userDateTime, err := time.Parse(dateTimeLayout, dateTime)
	if err != nil {
		return "", err
	}

	now := time.Now()
	for _, tz := range GetTimeZones() {
		loc, _ = time.LoadLocation(tz)

		if now.In(loc).Format(dateTimeShortLayout) == userDateTime.Format(dateTimeShortLayout) {
			return tz, nil
		}
	}

	return "", errors.New("Unable determine timezone of user provided string.")
}

// readDirRecursive - scan system dir. and files with timezones (recursive + ~agr.).
func readDirRecursive(path string) {
	files, _ := ioutil.ReadDir(zoneDir + path)
	for _, f := range files {
		if f.Name() != strings.ToUpper(f.Name()[:1])+f.Name()[1:] {
			continue
		}
		if f.IsDir() {
			readDirRecursive(path + "/" + f.Name())
		} else {
			timezones = append(timezones, (path + "/" + f.Name())[1:])
		}
	}
}
