package autotest

import (
	"testing"
	"time"

	"gitlab.meitu.com/platform/thanos/tools/autotest/cmd"

	"github.com/gomodule/redigo/redis"
)

//AutoClient check redis comman
type AutoClient struct {
	es *cmd.ExampleString
	el *cmd.ExampleList
	ek *cmd.ExampleKey
	*cmd.ExampleSystem
	em *cmd.ExampleMulti
	// addr string
	pool *redis.Pool
	conn redis.Conn
}

//NewAutoClient creat auto client
func NewAutoClient() *AutoClient {
	return &AutoClient{}
}

//Start run client
func (ac *AutoClient) Start(addr string) {
	// ac.pool = newPool(addr)
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	_, err = redis.String(conn.Do("auth", "test-1541501672-1-98d9882bb7a8ba2c16974e"))
	if err != nil {
		panic(err)
	}
	ac.conn = conn
	ac.es = cmd.NewExampleString(conn)
	ac.ek = cmd.NewExampleKey(conn)
	ac.el = cmd.NewExampleList(conn)
	ac.ExampleSystem = cmd.NewExampleSystem(conn)
	ac.em = cmd.NewExampleMulti(conn)
}

//Close shut client
func (ac *AutoClient) Close() {
	// ac.pool.Close()
	ac.conn.Close()
}

//StringCase check string case
//TODO
func (ac *AutoClient) StringCase(t *testing.T) {
	ac.es.SetEqual(t, "key-set", "value")
	ac.es.AppendEqual(t, "key-set", "value")
	ac.es.AppendEqual(t, "append", "value")
	ac.es.StrlenEqual(t, "key-set")
	ac.es.MSetEqual(t, "key-set", "value")
	ac.es.MGetEqual(t, "key-not-exist")
	ac.es.IncrEqual(t, "incr")
	ac.es.IncrEqual(t, "incr")
	ac.es.StrlenEqual(t, "heng")
}

//ListCase check list case
//TODO
func (ac *AutoClient) ListCase(t *testing.T) {
	ac.el.LpushEqual(t, "key-list", "v1", "v2", "v3", "v4")
	ac.el.LlenEqual(t, "key-list")
	ac.el.LsetEqual(t, "key-list", 3, "v0")
	ac.el.LindexEqual(t, "key-list", 3)
	ac.el.LrangeEqual(t, "key-list", 0, 10)
	ac.el.LrangeEqual(t, "key-list", 99, 100)
	ac.el.LpopEqual(t, "key-list")
	ac.el.LpopEqual(t, "key-list-l")
}

//KeyCase check key case
//TODO
func (ac *AutoClient) KeyCase(t *testing.T) {
	ac.ek.TTLEqual(t, "key-set", -1)
	ac.ek.RandomKeyEqual(t)
	ac.ek.ScanEqual(t, "", 3)
	ac.ek.ExistsEqual(t, 3, "key-set", "incr", "foo", "append")
	ac.ek.DelEqual(t, 3, "key-set", "incr", "foo", "append")
	ac.ek.TTLEqual(t, "key-set", -2)
	ac.ek.RandomKeyEqual(t)
	ac.ek.ScanEqual(t, "", 0)

	ac.es.SetEqual(t, "key-set", "value")
	ac.ek.ExpireEqual(t, "key-set", 2, 1)
	ac.ek.TTLEqual(t, "key-set", 1)
	time.Sleep(time.Second * 2)
	ac.ek.ExpireEqual(t, "key-set", 1, 0)
	ac.ek.ExpireEqual(t, "key-set", 0, 0)
}

//SystemCase check system case
func (ac *AutoClient) SystemCase(t *testing.T) {
	//auth
	ac.AuthEqual(t, "test-1541501672-1-98d9882bb7a8ba2c16974e")
	//ping
	ac.PingEqual(t)
}

//MultiCase bug mutil exec repley msg is error
func (ac *AutoClient) MultiCase(t *testing.T) {
	//multi
	ac.em.MultiEqual(t)
	ac.em.Cmd(t)
	// exec
	ac.em.ExecEqual(t)
}
