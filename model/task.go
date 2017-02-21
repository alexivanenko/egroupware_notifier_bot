//model package uses Bolt DB connection for interacting with store
package model

import (
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
)

const BUCKET_TASKS = "tasks"

type Task struct {
	ID     uint64
	Number int
}

//LoadTasks returns already notified task numbers list
func LoadTasks() ([]*Task, error) {
	var result []*Task

	if err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_TASKS))

		if bucket != nil {

			bucket.ForEach(func(k, v []byte) error {
				task := new(Task)
				json.Unmarshal(v, &task)

				result = append(result, task)

				return nil
			})
		}
		return nil
	}); err != nil {
		return result, err
	}

	return result, nil
}

//Save stored task to DB
func (task *Task) Save() error {

	if err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_TASKS))
		if err != nil {
			return err
		}

		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		task.ID = id

		if buf, err := json.Marshal(task); err != nil {
			return err
		} else if err := bucket.Put([]byte(strconv.FormatUint(task.ID, 10)), buf); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
