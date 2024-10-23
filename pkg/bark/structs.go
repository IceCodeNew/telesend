package bark

import "time"

type BarkSender struct {
	// Required; telegram user ID
	Creator int64 `xorm:"notnull"`
	// Required; generate by function UniqueID of pkg/uniqueID
	ID string `xorm:"pk notnull unique"`
	// Required; default: "https://api.day.app/"
	Server string `xorm:"notnull"`

	// Required
	DeviceKey []byte `xorm:"varchar(255) notnull"`
	// Required
	PreSharedSHA256IV []byte `xorm:"varchar(255) notnull"`
	// Required
	PreSharedSHA256Key []byte `xorm:"varchar(255) notnull"`

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type BarkMessage struct {
	// The number displayed next to App icon
	// Number greater than 9999 will be displayed as 9999+
	Badge int `json:"badge,omitempty"`
	// The content of the notification
	Body string `json:"body,omitempty"`
	// The value to be copied
	Copy string `json:"copy,omitempty"`
	// The group of the notification
	Group string `json:"group,omitempty"`
	// An url to the icon, available only on iOS 15 or later
	Icon string `json:"icon,omitempty"`
	// Value from https://github.com/Finb/Bark/tree/master/Sounds
	Sound string `json:"sound,omitempty"`
	// Notification title, optionally set by the sender
	Title string `json:"title,omitempty"`
	// Url that will jump when click notification
	URL string `json:"url,omitempty"`
}
