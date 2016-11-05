## Cistern

cistern collects binlog from each pump in cluster, and store them on disk in order of commitTS.

## How to use

```
Usage of cistern:
  -addr string
       	addr (i.e. 'host:port') to listen on for drainer connections (default "127.0.0.1:8249")
  -cluster-id uint
       	specifies the ID of TiDB cluster that cistern in charge of
  -collect-batch int
       	the max number of binlog items in a pulling batch (default 5000)
  -collect-interval int
       	the interval time (in seconds) of binlog collection loop (default 10)
  -config-file string
       	path to the configuration file
  -data-dir string
        path to the data directory of boltDB (default "data.cistern")
  -debug
       	whether to enable debug-level logging
  -deposit-window-period int
        a period of time (in minutes) after that the binlog items stored in boltDB will become to public state (default 10)
  -metrics-addr string
       	prometheus pushgateway address, leaves it empty will disable prometheus push.
  -metrics-interval int
       	prometheus client push interval in second, set "0" to disable prometheus push. (default 15)
  -pd-urls string
       	a comma separated list of PD endpoints (default "http://127.0.0.1:2379")
  -version
       	print version info
```


## Example

```
./bin/cistern -cluster-id 1
```
or use configuration file

```
./bin/cistern -config-file config.toml
```