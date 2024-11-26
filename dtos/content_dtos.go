package dtos

import "time"

type ContentDetails struct {
	Id            string                       `json:"id"`
	Content       map[string]any               `json:"content"`
	FieldChanges  []ContentDetailsFieldChange  `json:"fieldChanges"`
	FieldComments []ContentDetailsFieldComment `json:"fieldComments"`
	CreatedAt     time.Time                    `json:"createdAt"`
	UpdatedAt     time.Time                    `json:"updatedAt"`
}

type ContentDetailsFieldChange struct {
	Field   string                      `json:"field"`
	Changes []ContentDetailsUpdateField `json:"changes"`
}

type ContentDetailsUpdateField struct {
	Id            uint      `json:"id"`
	BeforeValue   any       `json:"beforeValue"`
	AfterValue    any       `json:"afterValue"`
	CreatedById   string    `json:"createdById"`
	CreatedByName string    `json:"createdByName"`
	CreatedAt     time.Time `json:"createdAt"`
}

type ContentDetailsFieldComment struct {
	Field    string                  `json:"field"`
	Comments []ContentDetailsComment `json:"comments"`
}

type ContentDetailsComment struct {
	Id            uint      `json:"id"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedById   string    `json:"createdById"`
	CreatedByName string    `json:"createdByName"`
}

type ContentSummary struct {
	Id        string         `json:"id"`
	Content   map[string]any `json:"content"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type ContentUpdateField struct {
	BeforeValue   any    `json:"beforeValue" binding:"required"`
	AfterValue    any    `json:"afterValue" binding:"required"`
	CreatedById   string `json:"createdById" binding:"required"`
	CreatedByName string `json:"createdByName" binding:"required"`
}

type ContentFieldComment struct {
	Comment       string `json:"comment" binding:"required"`
	CreatedById   string `json:"createdById" binding:"required"`
	CreatedByName string `json:"createdByName" binding:"required"`
}
