package model

import (
	"08_ihome/apiv1.0/pkg/setting"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var db * gorm.DB
var RE redis.Conn

func init()  {
	var (
		err error
		dbType, dbName, user, password, host string
	)

	sec ,err := setting.Setting.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database':%v",err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	host = sec.Key("HOST").String()
	password = sec.Key("PASSWORD").String()

	db,err = gorm.Open(dbType,fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil{
		log.Println(err)
	}

	db.SingularTable(true)
	db.AutoMigrate(&Orders{})
	db.LogMode(true)
	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(1000)

	RE,err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}

}

func RedeisGet(key string) (value map[string]interface{}) {
	value = make(map[string]interface{})
	reply, err :=redis.Bytes(RE.Do("get", key))
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(reply,&value)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func RedisSet(ex string,value_r []byte) bool {
	do, err := RE.Do("SET", "area_info", value_r, "EX", ex)
	if err != nil {
		fmt.Println(err)
	}
	if do == int64(1) {
		fmt.Println("success")
		return true
	} else {
		return false
	}
}

