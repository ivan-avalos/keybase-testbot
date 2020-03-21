# Keybase Testbot

This is only a simple bot written in Go. It runs in a Docker container with a Keybase client already installed and ready to work. Just set some environment variables containing the bot login params, and run!

## Build

```bash
docker build \
    --build-arg KEYBASE_USERNAME=$KEYBASE_USERNAME \
    --build-arg KEYBASE_PAPERKEY=$KEYBASE_PAPERKEY \
    -t $DOCKER_TAG .
```

## Run

```bash
docker run --rm -it $DOCKER_TAG
```