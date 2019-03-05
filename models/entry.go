package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entry struct {
	Title     string
	Uuid      string
	Topic     string
	Body      string
	UserId    string
	CreatedAt time.Time
	Abstract  string
	UserName  string
}

const PageCount = 5

//返回文章列表
func Entries() (entries []Entry, err error) {
	//Sess := DbSess.Copy()
	//defer Sess.Close()
	db := DbSess.DB("blog").C("Entry")
	err = db.Find(bson.M{}).Sort("-_id").Limit(PageCount).All(&entries)
	return
}

func EntriesCount() (enum int, err error) {
	db := DbSess.DB("blog").C("Entry")
	enum, err = db.Find(bson.M{}).Count()
	return
}

func EntriesPage(page int) (entries []Entry, err error) {
	db := DbSess.DB("blog").C("Entry")
	err = db.Find(bson.M{}).Sort("-_id").Skip((page - 1) * PageCount).Limit(PageCount).All(&entries)
	return
}

//创建文章
func (user User) CreateEntry(title string, topic string, body string, abstract string) (entry Entry, err error) {
	db := DbSess.DB("blog").C("Entry")
	entry = Entry{
		Uuid:      CreateUUID(),
		Topic:     topic,
		Body:      body,
		UserId:    user.Uuid,
		CreatedAt: time.Now(),
		Abstract:  abstract,
		UserName:  user.Name,
		Title:     title,
	}
	err = db.Insert(&entry)
	return
}

func ReadEntry(uuid string) (entry Entry, err error) {
	db := DbSess.DB("blog").C("Entry")
	err = db.Find(bson.M{"uuid": uuid}).One(&entry)
	return
}

func (user User) UpdateEntry(uuid string, title string, topic string, body string, abstract string) (err error) {
	entry := Entry{}
	db := DbSess.DB("blog").C("Entry")
	err = db.Find(bson.M{"uuid": uuid}).One(&entry)
	if err != nil {
		return
	}
	if user.Uuid != entry.UserId {
		return ErrWrongUser
	}
	entry.Title = title
	entry.Body = body
	entry.Abstract = abstract
	entry.Topic = topic
	err = db.Update(bson.M{"uuid": uuid}, &entry)
	return
}

func (user User) DeleteEntry(uuid string) error {
	entry := Entry{}
	db := DbSess.DB("blog").C("Entry")
	err := db.Find(bson.M{"uuid": uuid}).One(&entry)
	if err != nil {
		return err
	}
	if user.Uuid != entry.UserId {
		return ErrWrongUser
	}
	err = db.Remove(bson.M{"uuid": uuid})
	return err
}

func GetUserEntry(uuid string) (entries []Entry, err error) {
	db := DbSess.DB("blog").C("Entry")
	err = db.Find(bson.M{"userid": uuid}).Sort("-_id").Limit(PageCount).All(&entries)
	return
}

func UserEntriesCount(uuid string) (enum int, err error) {
	db := DbSess.DB("blog").C("Entry")
	enum, err = db.Find(bson.M{"userid": uuid}).Count()
	return
}

func UserEntriesPage(uuid string, page int) (entries []Entry, err error) {
	db := DbSess.DB("blog").C("Entry")
	err = db.Find(bson.M{"userid": uuid}).Sort("-_id").Skip((page - 1) * PageCount).Limit(PageCount).All(&entries)
	return
}
