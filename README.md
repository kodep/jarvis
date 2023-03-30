# kodep-jarvis

# FAQ

## How to start?

1. Setup golang v1.20
2. Run Mattermost local instance. Configure the bot and get credentials.
3. Configure `.env`
4. Run start the bot: `go run ./cmd/jarvis`

## How to update DI containers?

If you need to update the existing DI containers, you need to run:

```bash
make gen
```

If created a new one, then you need to install [wire](https://github.com/google/wire) first, then run:
```bash
wire <folder>
```

## How to deploy?

Simply use the existing Dockerfile. You will find environment variables inside `.env.example`
