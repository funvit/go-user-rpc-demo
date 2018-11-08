package repository

import (
	"context"
	"time"

	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/google/uuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

const userCollectionName = "users"
const repoName = "mongodb"

// UserRepositoryMongo is struct of params
type UserRepositoryMongo struct {
	db           *mongo.Database
	client       *mongo.Client
	datebaseName string
	timeout      time.Duration
}

type userBson struct {
	ID           string    `bson:"_id"`
	Login        string    `bson:"l"`
	RegisteredAt time.Time `bson:"r"`
}

func (s *userBson) ToDomainUser() *domain.User {
	return &domain.User{
		ID:           uuid.MustParse(s.ID),
		Login:        s.Login,
		RegisteredAt: s.RegisteredAt,
	}
}

func userToBson(user domain.User) *userBson {
	return &userBson{
		ID:           user.ID.String(),
		Login:        user.Login,
		RegisteredAt: user.RegisteredAt,
	}
}

// NewUserRepositoryMongo create new instance of UserRepositoryMongo
func NewUserRepositoryMongo(conn string, databaseName string, timeout time.Duration) *UserRepositoryMongo {
	client, err := mongo.Connect(context.Background(), conn, nil)
	if err != nil {
		panic(err)
	}

	pingContext, pingCancel := context.WithTimeout(context.Background(), timeout)
	defer pingCancel()
	if pingErr := client.Ping(pingContext, nil); pingErr != nil {
		logrus.Fatal("Cant connect to mongodb!")
	} else {
		logrus.Debugf("Mongodb: ping pong")
	}

	if databaseName == "" {
		logrus.Fatal("Param 'databaseName' cant be empty!")
	}

	return &UserRepositoryMongo{
		client:       client,
		datebaseName: databaseName,
		db:           client.Database(databaseName),
		timeout:      timeout,
	}
}

// AddUser implements user adding to DB
func (repo *UserRepositoryMongo) AddUser(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	coll := repo.db.Collection(userCollectionName)

	// при обновлении - тупо переписать поверх
	result, err := coll.InsertOne(
		ctx,
		userToBson(user),
	)
	if err != nil {
		return err
	}
	logrus.Debugf("UserRepository: AddUser result: %s", result)
	return nil
}

//UpdateUser
func (repo *UserRepositoryMongo) UpdateUser(user domain.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	coll := repo.db.Collection(userCollectionName)

	// при обновлении - тупо переписать поверх
	result, err := coll.UpdateOne(
		ctx,
		bson.NewDocument(
			bson.EC.String("_id", user.ID.String()),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.String("l", user.Login),
			),
		),
	)
	if err != nil {
		return false, err
	}
	// logrus.Debugf("UserRepository: AddUser result: %s", result)
	return result.ModifiedCount == 1, nil
}

//GetUserByID implements get user from DB
func (repo *UserRepositoryMongo) GetUserByID(id uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	coll := repo.db.Collection(userCollectionName)

	result := coll.FindOne(
		ctx,
		bson.NewDocument(
			bson.EC.String("_id", id.String()),
		),
	)
	u := &userBson{}
	if err := result.Decode(u); err != nil {
		return nil, err
	}
	return u.ToDomainUser(), nil
}
