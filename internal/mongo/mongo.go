package mongo

import (
	"FakeAPI/internal/db"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type Mongo struct {
	Client *mongo.Client
}

var instance *Mongo
var once sync.Once

func GetInstance() *Mongo {
	once.Do(func() {
		instance = &Mongo{}
	})
	return instance
}

const TIMEOUT = 5 * time.Second

func (m *Mongo) Connect(config *db.Config) (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+config.Username+":"+config.Password+"@"+config.Host+":"+config.Port))
	if err != nil {
		cancel()
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		cancel()
		return nil, err
	}
	m.Client = client
	return cancel, nil
}

func (m *Mongo) Get(query *db.Query, target interface{}) error {
	collection := m.Client.Database(query.Database).Collection(query.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	result := collection.FindOne(ctx, query.Filter)
	err := result.Decode(target)
	return err
}

func (m *Mongo) Insert(query *db.Query) error {
	collection := m.Client.Database(query.Database).Collection(query.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	_, err := collection.InsertOne(ctx, query.Object)
	return err
}

func (m *Mongo) Update(query *db.Query) error {
	collection := m.Client.Database(query.Database).Collection(query.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	_, err := collection.UpdateOne(ctx, query.Filter, bson.D{{
		Key:   "$set",
		Value: query.Object,
	}}, options.Update().SetUpsert(true))
	return err
}

func (m *Mongo) Close() {
	if m.Client != nil {
		m.Client.Disconnect(context.Background())
		m.Client = nil
	}
}
