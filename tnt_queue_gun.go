package main

import (
	"fmt"
	"log"
	"time"
	"math/rand"

	"github.com/tarantool/go-tarantool"

	"github.com/spf13/afero"
	"github.com/yandex/pandora/cli"
	"github.com/yandex/pandora/core"
	"github.com/yandex/pandora/core/aggregator/netsample"
	"github.com/yandex/pandora/core/import"
	"github.com/yandex/pandora/core/register"
)

type Ammo struct {
	Method   string
	TubeName string
	Params   map[string] interface {}
}

func customAmmoProvider() core.Ammo {
	return &Ammo{}
}

type GunConfig struct {
	Target []string `validate:"required"`
	User   string   `validate:"required"`
	Pass   string   `validate:"required"`
}

type Gun struct {
	conn *tarantool.Connection
	conf   GunConfig
	aggr   core.Aggregator
}

func NewGun(conf GunConfig) *Gun {
	return &Gun{
		conf: conf,
	}
}

func (g *Gun) Bind(aggr core.Aggregator, deps core.GunDeps) error {
	conn, err := tarantool.Connect(
		g.conf.Target[rand.Intn(len(g.conf.Target))],
		tarantool.Opts{
			User: g.conf.User,
			Pass: g.conf.Pass,
		},
	)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	g.conn = conn
	g.aggr = aggr

	return nil
}

func makeArgs(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return []interface{}{}
	}
	return args
}

func (g *Gun) queueCall(tube string, method string, args ...interface{}) (*tarantool.Response, error) {
	return g.conn.Call(fmt.Sprintf("queue.tube.%s:%s", tube, method), makeArgs(args))
}

func (g *Gun) Shoot(coreAmmo core.Ammo) {
	ammo := coreAmmo.(*Ammo)
	sample := netsample.Acquire(ammo.Method)

	code := 200
	var err error
	startTime := time.Now()
	switch ammo.Method {
	case "put":
		_, err = g.queueCall(ammo.TubeName, "put", ammo.Params["data"])
	case "take":
		_, err = g.queueCall(ammo.TubeName, "take")
	}
	sample.SetLatency(time.Since(startTime))
	if err != nil {
		log.Printf("Error %s task: %s", ammo.Method, err)
		code = 500
	}

	defer func() {
		sample.SetProtoCode(code)
		sample.AddTag(ammo.TubeName)
		g.aggr.Report(sample)
	}()
}

func main() {
	fs := afero.NewOsFs()
	coreimport.Import(fs)

	coreimport.RegisterCustomJSONProvider("tarantool_call_provider", customAmmoProvider)
	register.Gun("tnt_queue_gun", NewGun, func() GunConfig {
		return GunConfig{
			Target: []string{"localhost:3301"},
			User:   "guest",
			Pass:   "",
		}
	})
	cli.Run()
}
