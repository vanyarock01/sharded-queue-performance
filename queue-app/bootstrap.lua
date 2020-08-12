local conn = require('net.box').connect(
    'localhost:3301',
    { user = 'admin', password = 'queue-app-cluster-cookie' })


local bootstrap_cluster_str = [[
function bootstrap_cluster()
    require('fiber').create(function()
        require('cartridge.lua-api.edit-topology').edit_topology({
            replicasets = {
                {
                    alias = 'router-1', 
                    roles = {'sharded_queue.api'},
                    join_servers = {
                        {uri = 'localhost:3301'},
                    }
                },
                {
                    alias = 'router-2',
                    roles = {'sharded_queue.api'},
                    join_servers = {
                        {uri = 'localhost:3302'},
                    }
                },
                {
                    alias = 'storage-1',
                    roles = {'sharded_queue.storage'},
                    join_servers = {
                        {uri = 'localhost:3303'},
                        {uri = 'localhost:3304'},
                    }
                },
                {
                    alias = 'storage-2',
                    roles = {'sharded_queue.storage'},
                    join_servers = {
                        {uri = 'localhost:3305'},
                        {uri = 'localhost:3306'},
                    }
                }
            }
        })
        require("cartridge.lua-api.vshard").bootstrap_vshard()
        return true
    end)
end; return bootstrap_cluster()]]

conn:eval(bootstrap_cluster_str)
conn:close()

os.exit(0)
