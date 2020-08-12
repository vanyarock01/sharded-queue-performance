# #!/bin/bash

# build image with gun binary
docker build -t tnt-queue-gun-build -f build.Dockerfile . 
# run container
docker run -d --name temp-gun-build tnt-queue-gun-build
# get gun from container
docker cp temp-gun-build:/go/src/tnt_queue_gun .
# stop and remove container
docker stop temp-gun-build && docker container rm temp-gun-build
