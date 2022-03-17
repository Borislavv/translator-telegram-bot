package service

import (
	"io/ioutil"
	"strings"
	"sync"
)

var zoneDir string
var timeZones []string
var zoneDirs = []string{
	// Update path according to your OS
	"/usr/share/zoneinfo/",
	"/usr/share/lib/zoneinfo/",
	"/usr/lib/locale/TZ/",
}

type TimeZoneFetcherService struct {
	// deps.
	mx *sync.Mutex
}

// NewtimeZoneFetcherService - contructor of TimeZoneFetcherService structure.
func NewtimeZoneFetcherService(mx *sync.Mutex) *TimeZoneFetcherService {
	return &TimeZoneFetcherService{
		mx: mx,
	}
}

// GetTimeZones - wiil return a list with all possible timezones extracted from files of OS (linux).
func (service *TimeZoneFetcherService) GetTimeZones() []string {
	service.mx.Lock()
	defer func() {
		service.mx.Unlock()
	}()

	for _, zoneDir = range zoneDirs {
		service.readDirRecursive("")
	}

	return timeZones
}

// readDirRecursive - scan system dir. and files with timezones (recursive + ~agr.).
func (service *TimeZoneFetcherService) readDirRecursive(path string) {
	files, _ := ioutil.ReadDir(zoneDir + path)
	for _, f := range files {
		if f.Name() != strings.ToUpper(f.Name()[:1])+f.Name()[1:] {
			continue
		}
		if f.IsDir() {
			service.readDirRecursive(path + "/" + f.Name())
		} else {
			timeZones = append(timeZones, (path + "/" + f.Name())[1:])
		}
	}
}
