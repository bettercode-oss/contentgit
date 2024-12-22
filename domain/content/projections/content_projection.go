package projections

import (
	"contentgit/dtos"
	persistence "contentgit/ports/out/persistance"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ContentProjection struct {
	Id            string                `gorm:"primarykey"`
	TenantId      string                `gorm:"not null"`
	Content       persistence.JSONB     `gorm:"type:jsonb"`
	ContentType   string                `gorm:"type:varchar(100)"`
	FieldChanges  []ContentFieldChange  `gorm:"foreignKey:ContentId"`
	FieldComments []ContentFieldComment `gorm:"foreignKey:ContentId"`
	Version       uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func NewContentProjection(id string, tenantId string, content map[string]any, contentType string, version uint) ContentProjection {
	return ContentProjection{
		Id:          id,
		TenantId:    tenantId,
		Content:     content,
		ContentType: contentType,
		Version:     version,
	}
}

func (*ContentProjection) TableName() string {
	return "contents"
}

func (e *ContentProjection) UpdateField(fieldName string, updateField dtos.ContentUpdateField) {
	e.Content[fieldName] = updateField.AfterValue
	e.FieldChanges = append(e.FieldChanges, ContentFieldChange{
		Name: fieldName,
		Content: FieldUpdateVO{
			BeforeValue:   updateField.BeforeValue,
			AfterValue:    updateField.AfterValue,
			CreatedById:   updateField.CreatedById,
			CreatedByName: updateField.CreatedByName,
		},
	})
}

func (e *ContentProjection) AddFieldComment(fieldName, comment, createdById, createdByName string) {
	e.FieldComments = append(e.FieldComments, ContentFieldComment{
		Name:          fieldName,
		Comment:       comment,
		CreatedById:   createdById,
		CreatedByName: createdByName,
	})
}

type ContentFieldChange struct {
	gorm.Model
	ContentId string        `gorm:"not null"`
	Name      string        `gorm:"not null"`
	Content   FieldUpdateVO `gorm:"type:jsonb"`
}

func (ContentFieldChange) TableName() string {
	return "content_field_changes"
}

type ContentFieldComment struct {
	gorm.Model
	ContentId     string `gorm:"not null"`
	Name          string `gorm:"not null"`
	Comment       string `gorm:"not null;type:text"`
	CreatedById   string
	CreatedByName string
}

func (ContentFieldComment) TableName() string {
	return "content_field_comments"
}

type FieldUpdateVO struct {
	BeforeValue   any    `json:"beforeValue"`
	AfterValue    any    `json:"afterValue"`
	CreatedById   string `json:"createdById"`
	CreatedByName string `json:"createdByName"`
}

// Value Marshal
func (jsonField *FieldUpdateVO) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

// Scan Unmarshal
func (jsonField *FieldUpdateVO) Scan(value any) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &jsonField)
}
