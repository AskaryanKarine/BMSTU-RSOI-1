package models

type Person struct {
	ID      int32  `json:"id" validate:"omitempty"`
	Name    string `json:"name" validate:"required"`
	Age     int32  `json:"age" validate:"omitempty,gt=0"`
	Address string `json:"address" validate:"omitempty"`
	Work    string `json:"work" validate:"omitempty"`
}
