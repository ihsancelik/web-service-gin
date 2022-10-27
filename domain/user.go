package domain

import "time"

type User struct {
	Id                   string    `json:"id"`
	Sign                 string    `json:"sign"`
	TotalDailyLoginCount int       `json:"totalDailyLoginCount"`
	RegisteredDate       time.Time `json:"registeredDate"`
	LastLoginDate        time.Time `json:"lastLoginDate"`
}
