# HAnoProxY

## About The Project:

HAnoProxY is a DNS server which is built to bring high availability and load balancing for the services without using any kind of proxy server.

  

Proxy servers are handy and can make our systems more robust but sometimes the proxy server itself can be a bottleneck or single point of failure.

  

DNS protocol is designed to translate names to IP addresses, So if a DNS server is smart enough to detect service failure and automatically change the pointed IP address to healthy service then our entire system can be more reliable.

Furthermore, DNS is capable of acting as a load balancer and routes users to multiple servers and increases total performance.

  

## How it Works:

  

HAnoProxY is developed purley in Go and it can act as an authoritative and recursive DNS nameserver at the same time. It checks the health of listed endpoints and returns a healthy ip address to the Users/Services. It has two mechanisms for load balancing, by default it is **round robin** but it can be configured as an **Active-Passive** load balancing too.

health Check protocols:

-   HTTP Check
    
-   TCP Check
    
-   Redis Sentinel master/replica detection
    
-   PostgreSQL master-slave detectaion
    

  
  

## Redis Sentinel Example:

In this scenario we have 3 redis servers and redis-sentinel installed on them as well. Our application does not understand redis Sentinel and it only accepts one ip address/dns for redis connection.

The application uses HAnoProxY as a Dns Resolver.

  
  

Redis master = 10.10.10.1   
Redis replica1 = 10.10.10.2   
Redis replica2 = 10.10.10.3   
DnsRecord Name = redis   
Domain Name = ha.local   



So Servers can access to redis-master node by using this address:   
redis.ha.local

  
  

hanoproxy.yaml :

  

```yaml

GlobalOptions:

EnableRecursiveDnsServer: true

UpStreamDnsServer: "8.8.8.8"

HADomain: "ha.local"  #domain.

ListenPort: "53"  #udp

ListenIP: "0.0.0.0"

TTL: "10"  #Seconds

UpdateInterval: "60"  #Seconds



DnsRecords:

- Name: "redis"

  ServiceType: "sentinel"

  Options:

   CheckForHealth: true

   SentinelMonitorMasterName: "mymaster"

  Ip:

  - Addr: "10.10.10.1"  #predefined master ip (optional, it automatically detected by server)


  Sentinels:

  - Addr: "10.10.10.1"

    Port: "26379"

    Password: "auth_pass"
    

  - Addr: "10.10.10.2"

    Port: "26379"

    Password: "auth_pass"
    

  - Addr: "10.10.10.3"

    Port: "26379"

    Password: "auth_pass"

  

```

Before master failure:

```



+-----------+   +-----------+   +-----------+
|  redis    |   |  redis    |   |  redis    |
|  master   |   |  Replica1 |   |  Replica2 |
|10.10.10.1 |   |10.10.10.2 |   |10.10.10.3 |
+-----------+   +-----------+   +-----------+
      |
      |
      +---------+ (redis.ha.local = 10.10.10.1)
                |
                |
            +-----------+
            |  backend  |
            |           |
            +-----------+




```

After master failure:

```


+-----------+   +-----------+   +-----------+
| Unhealthy |   |  redis    |   |  redis    |
|    \ /    |   |new master |   |  Replica2 |
|    / \    |   |10.10.10.2 |   |10.10.10.3 |
+-----------+   +-----------+   +-----------+
                |               
                |   
                + (redis.ha.local = 10.10.10.2)
                |
                |
            +-----------+
            |  backend  |
            |           |
            +-----------+



```

## How to use it:

  

Download [hanoproxy](https://github.com/Abbas-gheydi/hanoproxy/releases) and install it by

```bash

sudo install ./hanoproxy /usr/local/bin

```

edit hanoproxy.yaml
Put hanoproxy.yaml into current directory or /etc/hanoproxy/
Then run hanoproxy.

By default, non-root users can not bind port 53(privileged-ports) so you can use this command to obtain the required permissions,   
or just run it as root (not recommended).


```bash

setcap 'cap_net_bind_service=+ep' /path/to/program

```

## License

MIT
