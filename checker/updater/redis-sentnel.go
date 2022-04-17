package updater

import (
	"context"
	"log"

	"github.com/Abbas-gheydi/hanoproxy/model"
	"github.com/go-redis/redis/v8"
)

func (rec *recordTable) checkRedisHealth(r model.Record) (updatedRecord model.Record) {
	for _, s := range r.Sentinels {
		master, err := GetRedisMaster(s.Addr, s.Port, s.Password, r.Options.SentinelMonitorMasterName)
		if err == nil {

			if len(r.Ip) == 0 {
				r.Ip = make([]model.IpAddr, 1)
			}
			for i := range r.Ip {
				rec.Mutex.Lock()
				r.Ip[i].Addr = master
				r.Ip[i].IsHealthy = true
				rec.Mutex.Unlock()
			}
			break
		}
	}
	return r
}

func GetRedisMaster(ip string, port string, password string, masterName string) (master string, err error) {

	sentinel := redis.NewSentinelClient(&redis.Options{
		Addr:     ip + ":" + port,
		Password: password,
	})

	addr, err := sentinel.GetMasterAddrByName(context.Background(), masterName).Result()
	if err == nil {
		master = addr[0]
		log.Println("sentinel", ip, port, "connected successfully.", "redis master is:", addr)
	} else {
		log.Println("sentinel", addr, err)
	}
	return
}
