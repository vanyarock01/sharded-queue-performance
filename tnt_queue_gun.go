package main

import (
	// "fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tarantool/go-tarantool"

	"github.com/spf13/afero"
	"github.com/yandex/pandora/cli"
	"github.com/yandex/pandora/core"
	"github.com/yandex/pandora/core/aggregator/netsample"
	"github.com/yandex/pandora/core/import"
	"github.com/yandex/pandora/core/register"
)

type Ammo struct {
	Collector string
	Min       float64
	Max       float64
}

func customAmmoProvider() core.Ammo {
	return &Ammo{}
}

type GunConfig struct {
	Target string `validate:"required"`
	User   string `validate:"required"`
	Pass   string
}

type Gun struct {
	conn *tarantool.Connection
	conf GunConfig
	aggr core.Aggregator
}

func NewGun(conf GunConfig) *Gun {
	return &Gun{
		conf: conf,
	}
}

func (g *Gun) Bind(aggr core.Aggregator, deps core.GunDeps) error {
	conn, err := tarantool.Connect(
		g.conf.Target,
		tarantool.Opts{
			User:    g.conf.User,
			Pass:    g.conf.Pass,
			Timeout: 3000 * time.Millisecond,
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

func ParseData(data []interface{}) uint64 {
	return data[0].([]interface{})[0].(uint64)
}

func (g *Gun) Shoot(coreAmmo core.Ammo) {
	ammo := coreAmmo.(*Ammo)
	sample := netsample.Acquire(ammo.Collector)

	code := 200
	var latency time.Duration
	log.Println(ammo.Min + (rand.Float64() * (ammo.Max - ammo.Min)))
	resp, err := g.conn.Call("observe", []interface{}{ammo.Collector, ammo.Min + (rand.Float64() * (ammo.Max - ammo.Min))})
	if err != nil {
		log.Printf("Error %s task: %s", "observe", err)
		code = 500
	} else {
		latency = time.Duration(ParseData(resp.Data) * 1000)
		// log.Println(latency)
		sample.SetLatency(latency)
	}

	defer func() {
		sample.SetProtoCode(code)
		sample.SetUserDuration(latency)
		g.aggr.Report(sample)
	}()
}

func main() {
	fs := afero.NewOsFs()
	coreimport.Import(fs)

	coreimport.RegisterCustomJSONProvider("tarantool_call_provider", customAmmoProvider)
	register.Gun("tnt_queue_gun", NewGun, func() GunConfig {
		return GunConfig{
			Target: "localhost:3302",
			User:   "U",
			Pass:   "X",
		}
	})
	cli.Run()
}
