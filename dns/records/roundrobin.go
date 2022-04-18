package records

import "log"

func (r *recordMap) roundRobinLB(requestQuery string) (responseIp string, count int) {
	count = len(r.Table[requestQuery].Ip)
	responseIpAddr := r.Table[requestQuery]

	for i := 0; i < count; i++ {

		responseIpAddr.LastHint = roundRobinCounter(responseIpAddr.LastHint, len(responseIpAddr.Ip))
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
	return

}

func roundRobinCounter(hint int, lengh int) int {
	if lengh > 1 {
		if hint+1 >= lengh {
			hint = 0
		} else {
			hint++
		}
	}
	return hint
}
