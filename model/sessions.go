package model

import (
	"time"

	common "github.com/copernet/whccommon/model"
)

func GetSessionById(id string) (*common.Session, bool) {
	var session = common.Session{}
	exist := walletdb.Where("session_id = ?", id).First(&session).RecordNotFound()
	// if not found the record, return nil and false directly
	if exist {
		return nil, false
	}

	return &session, true
}

func CreateSession(session *common.Session) error {
	return walletdb.Create(session).Error
}

func UpdateSessionById(session *common.Session) error {
	return walletdb.Model(session).Select("challenge", "p_challenge",
		"pub_key", "updated_at").
		Updates(map[string]interface{}{
			"challenge":   session.Challenge,
			"p_challenge": session.PChallenge,
			"pub_key":     session.PubKey,
			"updated_at":  time.Now(),
		}).Error
}
