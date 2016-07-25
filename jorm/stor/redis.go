// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package stor

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type RedisStor struct {
	Conn     redis.Conn
	DBPrefix string
}

type RedisConf struct {
	Network  string
	Address  string
	DBPrefix string
}

func NewRedisStor(conf *RedisConf) (*RedisStor, error) {

	conn, err := redis.Dial(conf.Network, conf.Address)
	if err != nil {
		return nil, err
	}

	return &RedisStor{
		Conn:     conn,
		DBPrefix: conf.DBPrefix + ":",
	}, nil
}

// Override io.Closer.Close
func (db *RedisStor) Close() error { return db.Conn.Close() }

func (db *RedisStor) dbKey(key string) string { return db.DBPrefix + key }

func (db *RedisStor) ReadBfKey(ipId uint32) (BfKeyInfo, bool) {

	kinfo := BfKeyInfo{}

	vals, err := redis.Values(db.Conn.Do("HMGET", db.dbKey("bf_keys:"+strconv.FormatUint(uint64(ipId), 16)), "key", "ttl"))
	if err != nil {
		return kinfo, false
	}

	_, err = redis.Scan(vals, &kinfo.Key, &kinfo.TTL)
	if err != nil {
		return kinfo, false
	}

	return kinfo, true
}
