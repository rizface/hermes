package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"strings"
	"time"
)

type config struct {
	username,password,host,port,dbname,command,seederPath,collectionName *string
}

func initConfig() *config {
	conf := new(config)

	conf.username = flag.String("username", "","username for connect to mongodb")
	conf.password = flag.String("password", "", "password for connect to mongodb")
	conf.host = flag.String("host", "localhost", "mongodb host")
	conf.port = flag.String("port","27017", "mongodb port")
	conf.dbname = flag.String("dbanme", "test", "database name")
	conf.command = flag.String("command", "up", "command to execute [up / down]")
	conf.seederPath = flag.String("path", "./seed", "seeder path")
	conf.collectionName = flag.String("collection", "all", "collection name [all / collectionName]")

	flag.Parse()

	return conf
}

func readSeeder(path *string) map[string]*options.CreateCollectionOptions {
	seeder := make(map[string]*options.CreateCollectionOptions)

	if path != nil {
		/*
			read all files inside dir
		*/
		dir,err := ioutil.ReadDir(*path)
		if err != nil  {
			logrus.Error(err)
			return nil
		}

		for _, file := range dir {
			/*
				loop and read file inside dir, except template.json file
			*/
			if file.Name() != "template.json" {
				var opt *options.CreateCollectionOptions
				// build filename & collection name
				fileName := fmt.Sprintf("%s/%s",*path,file.Name())
				collection := strings.Split(file.Name(),".")[0]

				// read seeder file
				content,err := ioutil.ReadFile(fileName)
				if err != nil {
					logrus.Error(err)
					return nil
				}

				// decode seeder file from json to *options.CreateCollectionOptions
				err = json.Unmarshal(content,&opt)
				if err != nil {
					logrus.Error(err)
					return nil
				}

				// set seeder inside map with key base on filename adn value is Unmarshal result
				seeder[collection] = opt
			}
		}

		return seeder
	}

	logrus.Error("seeder path is nil")

	return nil
}

func db(config *config) *mongo.Database {
	// initialize connection uri
	var uri string
	if *config.username == "" && *config.password == "" {
		uri = fmt.Sprintf("mongodb://%s:%s",*config.host,*config.port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s",*config.username,*config.password,*config.host,*config.port)
	}

	// connect to mongodb using uri
	opt := options.Client().ApplyURI(uri)
	client,err := mongo.Connect(context.Background(),opt)
	if err != nil {
		logrus.Error(err)

		return nil
	}

	// use db
	db := client.Database(*config.dbname)
	err = db.Client().Ping(context.Background(),&readpref.ReadPref{})
	if err != nil {
		logrus.Error(err)

		return nil
	}
	return db
}

func migrate(db *mongo.Database,config *config,seeders map[string]*options.CreateCollectionOptions) {
	for key, seeder := range seeders {
		logrus.Info("MIGRATE ", key, " ",time.Now())

		err := db.CreateCollection(context.Background(),key,seeder); if err != nil {
			if strings.Contains(err.Error(), "exists") {
				logrus.Error("FAILED MIGRATE : ", key, "IS EXISTS ", time.Now())
			} else {
				logrus.Error(err)

				return
			}
		}
	}
	logrus.Info("MIGRATION FINISHED ", time.Now())
}

func takeDownAll(db *mongo.Database) {
	err := db.Drop(context.Background()); if err != nil {
		logrus.Error(err)

		return
	}

	logrus.Info("ALL COLLECTION IS DELETED")
}

func takeDownCollections(db *mongo.Database,collections string) {
	items := strings.Split(collections,",")

	for _, item := range items {
		err := db.Collection(item).Drop(context.Background()); if err != nil {
			logrus.Error(err)
		} else {
			logrus.Info(item, " IS DELETED")
		}
	}
}

func main() {
	config := initConfig()
	db := db(config)

	if *config.command == "up" {
		logrus.Info("MIGRATION STARTED ", time.Now())
		seeder := readSeeder(config.seederPath)
		migrate(db,config,seeder)

		return
	} else if *config.command == "down" {
		if *config.collectionName == "all" {
			takeDownAll(db)
		} else {
			takeDownCollections(db,*config.collectionName)
		}

		return
	}
	logrus.Errorf("%s is invalid command", config.command)
}
