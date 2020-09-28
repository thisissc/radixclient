package radixclient

import (
	"log"
	"strings"
	"time"

	"github.com/mediocregopher/radix/v3"
)

var (
	_RadixPoolMap    = make(map[string]*radix.Pool, 0)
	_RadixClusterMap = make(map[string]*radix.Cluster, 0)
)

func Radix() radix.Client {
	name := "DEFAULT"
	cfg := _RadixConfigMap[name]

	if cfg.IsCluster {
		cluster, ok := _RadixClusterMap[name]
		if !ok {
			var err error
			poolFunc := func(network, addr string) (radix.Client, error) {
				return radix.NewPool(network, addr, cfg.MinPool,
					radix.PoolOnFullBuffer(cfg.MaxPool, time.Duration(cfg.DrainInterval)*time.Second))
			}
			cluster, err = radix.NewCluster(strings.Split(cfg.Addr, ","), radix.ClusterPoolFunc(poolFunc))
			if err != nil {
				log.Println(err)
			} else {
				_RadixClusterMap[name] = cluster
			}
		}

		return cluster
	} else {
		pool, ok := _RadixPoolMap[name]
		if !ok || pool.NumAvailConns() == 0 {
			var err error
			pool, err = radix.NewPool("tcp", cfg.Addr, cfg.MinPool,
				radix.PoolOnFullBuffer(cfg.MaxPool, time.Duration(cfg.DrainInterval)*time.Second))
			if err != nil {
				log.Println(err)
			} else {
				_RadixPoolMap[name] = pool
			}
		}

		return pool
	}
}
