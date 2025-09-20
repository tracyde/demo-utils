## Building

Build translator
```bash
docker build -t demo/translator -f docker/translator/Dockerfile .
```

Build receiver
```bash
docker build -t demo/receiver -f docker/receiver/Dockerfile .
```

## Running

Run receiver

```bash
docker run -p 8081:8081 -e SERVER_PORT=8081 --name receiver demo/receiver
```

Run translator
```bash
docker run -e SERVER_PORT=8081 -e SERVER_HOST='receiver.orb.local' -e SERVER_ENDPOINT='/ingest' -e SERVER_INPUTLOGPATH='/var/log/test.log' --name translator -v ./test.log:/var/log/test.log demo/translator
```
