# MessageBroker

![GitHub release (latest by date)](https://img.shields.io/github/v/release/lunarKettle/MessageBroker?display_name=tag)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lunarKettle/MessageBroker)

There is a simple message broker written in Go.

## Run in docker compose *(Recommended)*

Firstly you need to clone the project repository on your local machine:
```shell
git clone https://github.com/lunarKettle/message-broker.git
```

It's simple to run it with docker compose
```shell
make run
```

After this you will see the following message:
```shell
message-broker  | 2025/03/20 19:52:32 INFO Starting server Address=:8080
```

There is also some environment variables you can set, to configure project.

List of available environment variables above.

| Variable  | Description          | Default |
|-----------|----------------------|---------|
| `ADDRESS` | Address to listen on | `:8080` |




