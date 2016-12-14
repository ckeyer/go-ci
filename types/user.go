package types

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Email    string        `json:"email" bson:"email"`
	Phone    string        `json:"phone" bson:"phone"`
	Password string        `json:"-" bson:"password"`

	GithubAccount *GithubAccount `json:"github,omitempty" bson:"github,omitempty"`
	WechatAccount *WechatAccount `json:"wechat,omitempty" bson:"wechat,omitempty"`
}

type WechatAccount struct {
	Name string `json:"name" bson:"name"`
}

type GithubAccount struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"`
	Token    string `json:"token" bson:"token"`
}

func (g *GithubAccount) GetToken() string {
	if g.Token != "" {
		return g.Token
	}
	return g.Password
}
