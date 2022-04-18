package updater

import (
	"strings"
	"sync"
	"time"

	"github.com/Abbas-gheydi/hanoproxy/confs"
	"github.com/Abbas-gheydi/hanoproxy/dns/records"
	"github.com/Abbas-gheydi/hanoproxy/model"
)

type recordTable struct {
	Records []model.Record
	Mutex   sync.Mutex
}

func (r *recordTable) update(ip *model.IpAddr, healthStatus bool) {
	r.Mutex.Lock()
	ip.IsHealthy = healthStatus
	r.Mutex.Unlock()

}

var srcRecordTable recordTable

func init() {

	srcRecordTable.Records = make([]model.Record, 0)

}

func prepareRecordLists() {
	for i, src := range confs.HaNoProxy.DnsRecords {
		//add domain name and dot to records
		src.Name = src.Name + "." + confs.HaNoProxy.GlobalOptions.HADomain
		if !strings.HasSuffix(src.Name, ".") {
			src.Name = src.Name + "."
			src.Name = strings.ToLower(src.Name)
		}
		//make all ip healty for the beginning
		/*
			for it := range src.Ip {
				src.Ip[it].IsHealthy = true
			}
		*/

		confs.HaNoProxy.DnsRecords[i] = src

	}
}
func createNewRecordList() []model.Record {
	var allrecords []model.Record

	for _, record := range confs.HaNoProxy.DnsRecords {
		var r model.Record
		var ips []model.IpAddr
		var sentinels []model.IpAddr

		ips = append(ips, record.Ip...)
		sentinels = append(sentinels, record.Sentinels...)

		r.Ip = ips
		r.Sentinels = sentinels
		r.Name = record.Name
		r.Options = record.Options
		r.ServiceType = record.ServiceType
		allrecords = append(allrecords, r)

	}
	return allrecords
}
func FillRecordTables() {
	prepareRecordLists()

	srcRecordTable.Records = createNewRecordList()
	for _, src := range createNewRecordList() {
		records.DnsRecords.Add(src)

	}

}

func HealthCheck() {

	for {
		for _, src := range srcRecordTable.Records {
			if src.Options.CheckForHealth {

				switch src.ServiceType {

				case model.Postgres:

					go records.DnsRecords.Update(srcRecordTable.checkPostgresHealth(src))

				case model.Sentinel:
					go records.DnsRecords.Update(srcRecordTable.checkRedisHealth(src))
				case model.HTTP:
					go records.DnsRecords.Update(srcRecordTable.checkHttpHealth(src))
				case model.TCP:
					go records.DnsRecords.Update(srcRecordTable.checkTcp(src))
				}
			}
		}
		time.Sleep(time.Second * time.Duration(confs.HaNoProxy.GlobalOptions.UpdateInterval))
	}

}
