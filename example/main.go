package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/liyuan1125/gorm-cache"
	redis2 "github.com/liyuan1125/gorm-cache/store/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	db *gorm.DB

	redisClient *redis.Client

	cachePlugin *cache.Cache
)

func newDb() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/flock?charset=utf8&parseTime=True&loc=Local"
	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	redisClient = redis.NewClient(&redis.Options{Addr: ":6379",DB: 10})

	cacheConfig := &cache.Config{
		Store:      redis2.NewWithDb(redisClient), // OR redis2.New(&redis.Options{Addr:"6379"})
		Serializer: &cache.DefaultJSONSerializer{},
	}

	cachePlugin = cache.New(cacheConfig)

	if err = db.Use(cachePlugin); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func basic() {
	var name string
	ctx := context.Background()
	ctx = cache.NewExpiration(ctx, time.Hour)

	db.Table("users").WithContext(ctx).Where("id = 11642759").Limit(1).Pluck("name", &name)
	fmt.Println(name)
	// output gorm
}

func customKey() {
	var name string
	ctx := context.Background()
	ctx = cache.NewExpiration(ctx, time.Hour)
	ctx = cache.NewKey(ctx, "name")

	db.Table("users").WithContext(ctx).Where("id = 11642759").Limit(1).Pluck("name", &name)

	fmt.Println(name)
	// output gormwithmysql
}

func useTag() {
	var name string
	ctx := context.Background()
	ctx = cache.NewExpiration(ctx, time.Hour)
	ctx = cache.NewTag(ctx, "users")

	db.Table("users").WithContext(ctx).Where("id = 11642759").Limit(1).Pluck("name", &name)

	fmt.Println(name)
	// output gormwithmysql
}

func main() {
	newDb()
	//basic()
	//customKey()
	useTag()

	ctx := context.Background()
	fmt.Println(redisClient.Keys(ctx, "*").Val())

	//fmt.Println(cachePlugin.RemoveFromTag(ctx, "users"))
}
