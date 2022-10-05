package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri, dbName string

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func init() {
	// read config file
	if err := initConfig(); err != nil {
		fmt.Print("Error init config")
	}
	dbName = viper.GetString("mongodb.dbname")
	uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&readPreference=primary&directConnection=true&ssl=false",
		os.Getenv("MONGODB_USER"),
		os.Getenv("MONGODB_PASSWORD"),
		viper.GetString("mongodb.host"),
		viper.GetString("mongodb.port"),
		dbName,
	)
}

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// db := client.Database(dbName)

	// t.Run("Insertexample")
}
