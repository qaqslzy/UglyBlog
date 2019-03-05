package models

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
	"log"
)

var DbSess *mgo.Session

var (
	ErrUserExists    = errors.New("User already exists")
	ErrPasswordWorng = errors.New("Password Wrong")
	ErrWrongUser     = errors.New("User is Wrong")
)

func init() {
	var err error
	DbSess, err = mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatal(err)
	}

	DbSess.SetMode(mgo.Eventual, true)
	return
}

func CreateUUID() (id string) {

	tradeID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return tradeID.PureString()
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
