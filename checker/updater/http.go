package updater

import (
	"log"
	"net/http"

	"github.com/Abbas-gheydi/hanoproxy/model"
)

func (rec *recordTable) checkHttpHealth(r model.Record) (updatedRecord model.Record) {
	for i := range r.Ip {

		go func(iteration int) {

			healthstatus := isURLhealthy(r.Ip[iteration].Url, r.Options.Expected_Response_Code, r.Options.RetryCount)

			rec.update(&r.Ip[iteration], healthstatus)

		}(i)

	}

	return r
}

func isURLhealthy(url string, code int, retryCount int) bool {
	health := false
	defer func() {
		if err := recover(); err != nil {
			log.Println("http panic occurred:", err)
		}
	}()

	for i := 0; i <= retryCount; i++ {

		resp, err := http.Get(url)
		if err != nil {
			log.Print("http", err)
			health = false
		} else {

			defer resp.Body.Close()

			if resp.StatusCode == code {

				health = true
				log.Println(url, "Response status:", resp.StatusCode)
				break
			}

		}

	}

	return health
}
