package godingtalk


type Recordresult struct{
	BaseCheckTime int64 `json:"baseCheckTime"`
	CheckType string `json:"checkType"`
	CorpID string `json:"corpId"`
	GroupID int `json:"groupId"`
	ID int `json:"id",gorm:"primary_key"`
	LocationResult string `json:"locationResult"`
	PlanID int `json:"planId"`
	RecordID int `json:"recordId,omitempty"`
	TimeResult string `json:"timeResult"`
	UserCheckTime int64 `json:"userCheckTime"`
	UserID string `json:"userId"`
	WorkDate int64 `json:"workDate"`
}

type Attendance struct {
	OAPIResponse
	Recordresult []Recordresult
}


func (c *DingTalkClient) Attendances() (Attendance, error) {
	var data Attendance
	err := c.httpRPC("attendance/list", nil, nil, &data)
	return data, err
}