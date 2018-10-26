#! /bin/sh
docker run --env REDIRECT_HOST=https://cnn.com --rm -it -p 80:80 cycleplatform/goredirector:latest