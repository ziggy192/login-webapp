package model

import (
	"encoding/json"
	"time"
)

type BaseResponse struct {
	Message string         `json:"message"`
	Time    time.Time      `json:"time"`
	Data    map[string]any `json:"data"`
}

func (b *BaseResponse) UnmarshalData(v any) error {
	marshal, err := json.Marshal(b.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, v)
}
