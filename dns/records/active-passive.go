package records

import "log"

func (r *recordMap) activePassiveLB(requestQuery string) (responseIp string, count int) {
	count = len(r.Table[requestQuery].Ip)
	responseIpAddr := r.Table[requestQuery]

	for i := 0; i < count; i++ {

		responseIpAddr.LastHint = roundRobinCounter(responseIpAddr.LastHint, len(responseIpAddr.Ip))
		r.Table[requestQuery] = responseIpAddr

		if responseIpAddr.Ip[responseIpAddr.LastHint].IsHealthy || !responseIpAddr.Options.CheckForHealth {
			responseIp = responseIpAddr.Ip[responseIpAddr.LastHint].Addr
			responseIpAddr.LastHint = activePassiveCounter(responseIpAddr.LastHint, len(responseIpAddr.Ip))
			r.Table[requestQuery] = responseIpAddr
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
func activePassiveCounter(hint int, lengh int) int {
	if lengh > 1 {
		if hint != 0 {
			hint--
		} else {
			hint = lengh
		}
	}
	return hint
}
