#!/usr/bin/env tarantool

require('strict').on()

local cartridge = require('cartridge')
local ok, err = cartridge.cfg({
    workdir = 'tmp/db',
    roles = {
        'cartridge.roles.vshard-storage',
        'cartridge.roles.vshard-router',
        'sharded_queue.storage',
        'sharded_queue.api',
    },
    cluster_cookie = 'queue-app-cluster-cookie',
})

assert(ok, tostring(err))
