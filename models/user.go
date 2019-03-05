package models

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type User struct {
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Uuid      string
	Email     string
	UserId    string
	CreatedAt time.Time
}

// 创建cookie
func (user *User) CreateSession() (sess Session, err error) {
	sess = Session{
		Uuid:      CreateUUID(),
		UserId:    user.Uuid,
		Email:     user.Email,
		CreatedAt: time.Now(),
	}
	db := DbSess.DB("blog").C("Session")
	err = db.Insert(&sess)
	return
}

// 返回用户的session
func (user *User) Session() (session Session, err error) {
	db := DbSess.DB("blog").C("Session")
	session = Session{}
	err = db.Find(bson.M{"userid": user.Uuid}).One(&session)
	return
}

// 校验session的正确性
func (session *Session) Check() (tsess Session, valid bool, err error) {
	db := DbSess.DB("blog").C("Session")
	user := User{}
	err = db.Find(bson.M{"uuid": session.Uuid}).One(&tsess)
	if err != nil {
		valid = false
		return
	}
	db = DbSess.DB("blog").C("User")
	err = db.Find(bson.M{"uuid": tsess.UserId}).One(&user)
	if err != nil {
		valid = false
		return
	}
	if tsess.UserId == user.Uuid {
		valid = true
		return
	}
	return
}

// 通过UUID删除session
func (sess Session) DeleteByUUID() {
	db := DbSess.DB("blog").C("Session")
	err := db.Remove(bson.M{"uuid": sess.Uuid})
	if err != nil {
		log.Fatal(err)
	}
}

// 通过session找到用户
func (sess Session) User() (user User, err error) {
	db := DbSess.DB("blog").C("User")
	err = db.Find(bson.M{"uuid": sess.UserId}).One(&user)
	return
}

// 新建 用户
func (user *User) Creat() (err error) {
	db := DbSess.DB("blog").C("User")

	//查看该注册邮箱是否已经存在用户
	exuser := User{}
	err = db.Find(bson.M{"email": user.Email}).One(&exuser)
	if err == nil {
		err = ErrUserExists
		return
	}

	//完善用户信息
	user.Uuid = CreateUUID()
	user.Password = Encrypt(user.Password)
	user.CreatedAt = time.Now()

	err = db.Insert(user)
	return
}

// 根据email寻找用户
func UserByEmail(email string) (user User, err error) {
	db := DbSess.DB("blog").C("User")
	user = User{}
	err = db.Find(bson.M{"email": email}).One(&user)
	return
}
