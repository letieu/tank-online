## Terminal multiple players tank game

- A simple tank game just for learn golang
- Use redis as game server

![image](https://github.com/letieu/tank-online/assets/53562817/4aeffbc2-a13d-464f-b07e-2a62acf76cb4)

## Usage

```bash

# Start a redis server
docker run --name redis -p 6379:6379 -e REDIS_PASSWORD=secret bitnami/redis:latest

#Start game with redis server on localhost:6379, redis password is secret
./tank --name=letieu --host=localhost:6379 --pass=secret

```

## WIP
- [X] Play tank inside terminal
- [x] Multiple player
- [x] View port
- [ ] Leader board
- [ ] Configurable via UI (Charm)
- [ ] Play via SSH with wish
