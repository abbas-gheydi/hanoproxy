package records

import (
	"errors"
	"log"
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

func roundRobind(hint int, lengh int) int {
	if lengh > 1 {
		if hint+1 >= lengh {
			hint = 0
		} else {
			hint++
		}
	}
	return hint
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

		responseIpAddr := r.Table[requestQuery]

		for i := 0; i < count; i++ {

			responseIpAddr.LastHint = roundRobind(responseIpAddr.LastHint, len(responseIpAddr.Ip))
			r.Table[requestQuery] = responseIpAddr
			if responseIpAddr.Ip[responseIpAddr.LastHint].IsHealthy || !responseIpAddr.Options.CheckForHealth {
				responseIp = responseIpAddr.Ip[responseIpAddr.LastHint].Addr
				break

			} else {
				//if all ip address are unhealthy then return first one
				if i+1 == count {

					responseIp = responseIpAddr.Ip[0].Addr
					log.Println("Warining, All Ip Addresses are unHealthy,then just return first record for: ", requestQuery)
				}

			}
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
