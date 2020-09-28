package radixclient

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mediocregopher/radix/v3"
	"github.com/pkg/errors"
)

func Save2Redis(key string, expire int, data interface{}) error {
	c := Radix()

	bytesData, _ := json.Marshal(data)
	pipe := radix.Pipeline(
		radix.Cmd(nil, "SET", key, string(bytesData)),
		radix.FlatCmd(nil, "EXPIRE", key, expire),
	)
	err := c.Do(pipe)
	if err != nil {
		return errors.Wrap(err, "Save2Redis failed")
	}
	return nil
}

func LoadFromRedis(key string, data interface{}) error {
	c := Radix()

	var bytesData []byte
	err := c.Do(radix.Cmd(&bytesData, "GET", key))
	if err != nil {
		return errors.Wrap(err, "LoadFromRedis failed")
	}

	err = json.Unmarshal(bytesData, data)
	if err != nil {
		return errors.Wrap(err, "Json unmarshal failed")
	}

	return nil
}

func Save2RedisMutex(key string, expire int, data interface{}) error {
	c := Radix()

	replicaKey := fmt.Sprintf("_replica_%s", key)
	mutexKey := fmt.Sprintf("_mutex_%s", key)

	bytesData, _ := json.Marshal(data)

	pipe := radix.Pipeline(
		radix.Cmd(nil, "SET", key, string(bytesData)),
		radix.FlatCmd(nil, "EXPIRE", key, expire),
		radix.Cmd(nil, "SET", replicaKey, string(bytesData)),
		radix.FlatCmd(nil, "EXPIRE", replicaKey, expire+120), // 2min more
		radix.Cmd(nil, "DEL", mutexKey),
	)

	err := c.Do(pipe)
	if err != nil {
		return errors.Wrap(err, "Save2Redis failed")
	}
	return nil
}

func LoadFromRedisMutex(key string, data interface{}) (bool, error) {
	c := Radix()

	var err error
	var bytesData []byte
	dataExpired := false

	err = c.Do(radix.Cmd(&bytesData, "GET", key))
	if err != nil {
		dataExpired = true

		replicaKey := fmt.Sprintf("_replica_%s", key)
		err = c.Do(radix.Cmd(&bytesData, "GET", replicaKey))
		if err != nil {
			return true, errors.Wrap(err, "LoadFromRedis failed")
		}
	}

	err = json.Unmarshal(bytesData, data)
	if err != nil {
		return true, errors.Wrap(err, "Json unmarshal failed")
	}

	if dataExpired {
		mutexKey := fmt.Sprintf("_mutex_%s", key)
		var cnt int
		err2 := c.Do(radix.Cmd(&cnt, "INCR", mutexKey))
		if err2 != nil {
			log.Println(err2)
		}

		dataExpired = cnt == 1

		if dataExpired {
			err2 = c.Do(radix.FlatCmd(nil, "EXPIRE", mutexKey, 10))
			if err2 != nil {
				log.Println(err2)
			}
		}
	}

	return dataExpired, err
}
