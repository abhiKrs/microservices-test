package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
)

type CreateSourceRequest struct {
	Name       string `json:"name" validate:"required"`
	SourceType uint8  `json:"sourceType" validate:"required"`
	TeamId     string `json:"teamId" validate:"required"`
}

type UpdateSourceRequest struct {
	Name string `json:"name" validate:"required"`
}

func (r *UpdateSourceRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

// type UpdateSourceResponse struct {
// 	Name string `json:"name" validate:"required"`
// }

type CreateSourceResponse struct {
	IsSuccessful bool               `json:"isSuccessful"`
	Message      []string           `json:"message,omitempty"`
	Data         SourceBodyResponse `json:"data,omitempty"`
}

type DeleteSourceResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}

type GetSourceResponse struct {
	IsSuccessful bool               `json:"isSuccessful"`
	Message      []string           `json:"message,omitempty"`
	Data         SourceBodyResponse `json:"data,omitempty"`
}

type GetAllSourceResponse struct {
	IsSuccessful bool                 `json:"isSuccessful"`
	Message      []string             `json:"message,omitempty"`
	Data         []SourceBodyResponse `json:"data,omitempty"`
}

type SourceBodyResponse struct {
	SourceId    uuid.UUID `json:"sourceId"`
	ProfileId   uuid.UUID `json:"profileId"`
	Name        string    `json:"name"`
	SourceType  uint8     `json:"sourceType"`
	SourceToken string    `json:"sourceToken,omitempty"`
	TeamId      uuid.UUID `json:"teamId"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (r *CreateSourceRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
