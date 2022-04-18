package records

import (
	"errors"
	"strings"
	"sync"

	"github.com/Abbas-gheydi/hanoproxy/model"
)

var (
	DnsRecords recordspool
)

type recordspool interface {
	Lookup(requestQuery string) (responseIp string, count int, err error)
	Add(newR model.Record)
	Update(rec model.Record)
	Initialize()
}

func newInmemoryRecordsPool() recordspool {

	return &recordMap{}
}

func init() {

	DnsRecords = newInmemoryRecordsPool()
	DnsRecords.Initialize()

}

type recordMap struct {
	Table map[string]model.Record
	Mutex sync.Mutex
}

func (r *recordMap) Initialize() {
	r.Table = make(map[string]model.Record)

}

func (r *recordMap) Lookup(requestQuery string) (responseIp string, count int, err error) {
	//convert all name to lowercase
	requestQuery = strings.ToLower(requestQuery)
	r.Mutex.Lock()
	count = len(r.Table[requestQuery].Ip)

	if count > 0 {
		//choose load balancing method
		if r.Table[requestQuery].Options.LBmethod == model.ActivePassive {

			responseIp, count = r.activePassiveLB(requestQuery)
		} else {

			responseIp, count = r.roundRobinLB(requestQuery)
		}

	}

	r.Mutex.Unlock()
	if responseIp == "" {
		err = errors.New(requestQuery + " not found")
	}

	return
}

func (r *recordMap) Add(newR model.Record) {

	r.Mutex.Lock()
	r.Table[newR.Name] = newR

	r.Mutex.Unlock()
}

func (r *recordMap) Update(rec model.Record) {
	r.Mutex.Lock()
	r.Table[rec.Name] = rec
	r.Mutex.Unlock()

}
