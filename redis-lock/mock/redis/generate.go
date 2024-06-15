package redismock

//go:generate mockgen -package=redismock -destination=./cmd.mock.go github.com/redis/go-redis/v9 Cmdable
