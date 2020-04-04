package config

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type Test struct {
	Id      int64  `gorm:"-;primary_key;AUTO_INCREMENT;column:id" json:"id" bson:"id"`
	Msg     string `gorm:"column:msg" json:"msg" bson:"msg"`
	Name    string `gorm:"column:name" json:"name" bson:"name"`
	AddTime int64  `gorm:"column:add_time" json:"add_time" bson:"add_time"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func List(w http.ResponseWriter, r *http.Request) {
	client, err := MongoClient()
	FailOnError(err, "MongoClient")

	collection := client.Database("testing")
	numbersCollection := collection.Collection("numbers")
	optsLimit := options.Find().SetLimit(10)
	optsSort := options.Find().SetSort(bson.M{"add_time": 1})
	filter := bson.D{}
	cursor, err := numbersCollection.Find(nil, filter, optsLimit, optsSort)
	if cursor == nil {
		defer cursor.Close(nil)
	}
	if err != nil && err != mongo.ErrNilDocument {
		FailOnError(err, "numbersCollection.Find")
	}

	w.Header().Set("content-type", "application/json;charset=utf-8")

	var results []Test
	if cursor == nil || err == mongo.ErrNilDocument {
		db, _ := gorm.Open("mysql", "root:Nm123456.@/test?charset=utf8&parseTime=True&loc=Local")
		db.LogMode(true) //打印mysql 日子
		defer db.Close()

		db = db.Debug().Table("tests").Limit(10).Order("add_time", true).Find(&results)
		if err := json.NewEncoder(w).Encode(results); err != nil {
			FailOnError(err, "json.NewEncoder(w)")
		}
		return
	}
	for cursor.Next(nil) {
		var elem Test
		if err := cursor.Decode(&elem); err != nil {
			FailOnError(err, "cursor.Decode(&elem)")
		}
		results = append(results, elem)
	}

	if err := cursor.Err(); err != nil {
		FailOnError(err, "cursor.Err")
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		FailOnError(err, "json.NewEncoder(w)")
	}
}
