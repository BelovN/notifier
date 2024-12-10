package models

type User struct {
	Id           uint
	Username     string
	IsSubscribed bool
	ChannelId    int64
}
