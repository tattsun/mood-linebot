package model

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID     bson.ObjectId `bson:"_id"`
	UserID string        `bson:"user_id"`
}

type UserRepository struct {
	col *mgo.Collection
}

func NewUserRepository(db *mgo.Database) (*UserRepository, error) {
	col := db.C("users")

	userIDIndex := mgo.Index{
		Key:    []string{"user_id"},
		Unique: true,
	}
	if err := col.EnsureIndex(userIDIndex); err != nil {
		return nil, errors.Wrap(err, "failed to create index 'user_id' to collection 'users'")
	}

	return &UserRepository{
		col: col,
	}, nil
}

func (r *UserRepository) Create(u *User) error {
	return r.col.Insert(u)
}

func (r *UserRepository) FindAll() ([]User, error) {
	var users []User
	if err := r.col.Find(nil).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}
