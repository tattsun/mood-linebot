package model

import (
	"time"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type Mood struct {
	ID        bson.ObjectId `bson:"_id"`
	UserID    string        `bson:"user_id"`
	Mood      int           `bson:"mood"`
	Timestamp time.Time     `bson:"timestamp"`
}

type MoodRepository struct {
	col *mgo.Collection
}

func NewMoodRepository(db *mgo.Database) (*MoodRepository, error) {
	col := db.C("moods")
	return &MoodRepository{
		col: col,
	}, nil
}

func (r *MoodRepository) Create(m *Mood) error {
	return r.col.Insert(m)
}

func (r *MoodRepository) FindAll() ([]Mood, error) {
	var moods []Mood
	if err := r.col.Find(nil).All(&moods); err != nil {
		return nil, err
	}
	return moods, nil
}
