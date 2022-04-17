package updater

import (
	"log"
	"net"
	"time"

	"github.com/Abbas-gheydi/hanoproxy/model"
)

func (rec *recordTable) checkTcp(r model.Record) (updatedRecord model.Record) {
	for i := range r.Ip {

		go func(iteration int) {

			healthstatus := tcp_raw_connect(r.Ip[iteration].Addr, r.Ip[iteration].Port, r.Options.RetryCount)

			rec.update(&r.Ip[iteration], healthstatus)

		}(i)

	}

	return r
}

func tcp_raw_connect(host string, port string, retryCount int) (isRechable bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("tcp panic occurred:", err)
		}
	}()
	for i := 0; i <= retryCount; i++ {

		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			log.Println("tcp Connecting error:", err)
		}
		if conn != nil {
			defer conn.Close()
			log.Println("tcp Success", net.JoinHostPort(host, port))
			isRechable = true
			break
		}
	}
	return

}
