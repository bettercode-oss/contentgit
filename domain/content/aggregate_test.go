package content

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewContentAggregate(t *testing.T) {
	t.Run("id가 없으면 error를 반환한다", func(t *testing.T) {
		// given
		aggregateId := ""
		tenantId := "bettercode"

		// when
		_, err := NewContentAggregate(aggregateId, tenantId)

		// then
		assert.Error(t, err)
	})

	t.Run("tenantId가 없으면 error를 반환한다", func(t *testing.T) {
		// given
		aggregateId := uuid.New().String()
		tenantId := ""

		// when
		_, err := NewContentAggregate(aggregateId, tenantId)

		// then
		assert.Error(t, err)
	})

	t.Run("필수 값이 있으면 aggregate를 반환한다", func(t *testing.T) {
		// given
		aggregateId := uuid.New().String()
		tenantId := "bettercode"

		// when
		aggregate, err := NewContentAggregate(aggregateId, tenantId)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, aggregate)
		assert.Equal(t, ContentAggregateType, aggregate.GetType())
		assert.Equal(t, aggregateId, aggregate.GetID())
		assert.Equal(t, tenantId, aggregate.GetTenantId())
	})
}

func TestContentAggregate_CreateContent(t *testing.T) {
	t.Run("Content가 nil이면 error를 반환한다", func(t *testing.T) {
		// given
		sut, _ := NewContentAggregate(uuid.New().String(), "bettercode")

		// when
		err := sut.CreateContent(context.Background(), nil)

		// then
		assert.Error(t, err)
	})
	t.Run("Content를 생성한다", func(t *testing.T) {
		// given
		sut, _ := NewContentAggregate(uuid.New().String(), "bettercode")

		// when
		err := sut.CreateContent(context.Background(), map[string]any{"name": "test"})

		// then
		assert.NoError(t, err)
		assert.Equal(t, "test", sut.Content["name"])
	})
}

func TestContentAggregate_UpdateField(t *testing.T) {
	t.Run("필드명이 없으면 error를 반환한다", func(t *testing.T) {
		// given
		sut, _ := NewContentAggregate(uuid.New().String(), "bettercode")
		sut.Content = map[string]any{"name": "홍길동"}

		// when
		err := sut.UpdateField(context.Background(), "unknownField", "홍길동", "고길동", "testerId", "testerName")

		// then
		assert.Equal(t, ErrFieldNotFound, err)
	})

	t.Run("값이 충돌하면 error를 반환한다", func(t *testing.T) {
		// given
		sut, _ := NewContentAggregate(uuid.New().String(), "bettercode")
		sut.Content = map[string]any{"name": "홍길동"}

		// when
		err := sut.UpdateField(context.Background(), "name", "고길동", "둘리", "testerId", "testerName")

		// then
		assert.Equal(t, ErrFieldUpdateConflict, err)
	})

	t.Run("필드를 업데이트한다", func(t *testing.T) {
		// given
		sut, _ := NewContentAggregate(uuid.New().String(), "bettercode")
		sut.Content = map[string]any{"name": "홍길동"}

		// when
		err := sut.UpdateField(context.Background(), "name", "홍길동", "고길동", "testerId", "testerName")

		// then
		assert.NoError(t, err)
		assert.Equal(t, "고길동", sut.Content["name"])
	})
}

func TestContentAggregate_AddFieldComment(t *testing.T) {
	t.Run("필드명이 없으면 error를 반환한다", func(t *testing.T) {
		// given
		sut, _ := NewContentAggregate(uuid.New().String(), "bettercode")
		sut.Content = map[string]any{"name": "홍길동"}

		// when
		err := sut.AddFieldComment(context.Background(), "unknownField", "comment", "testerId", "testerName")

		// then
		assert.Equal(t, ErrFieldNotFound, err)
	})

	t.Run("필드에 댓글을 추가한다", func(t *testing.T) {
		// given
		sut, _ := NewContentAggregate(uuid.New().String(), "bettercode")
		sut.Content = map[string]any{"name": "홍길동"}

		// when
		err := sut.AddFieldComment(context.Background(), "name", "comment", "testerId", "testerName")

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, len(sut.FieldComments))
		actualFieldComment := sut.FieldComments[0]
		assert.Equal(t, "name", actualFieldComment.FieldName)
		assert.Equal(t, "comment", actualFieldComment.Comments[0].Comment)
		assert.Equal(t, "testerId", actualFieldComment.Comments[0].CreatedById)
		assert.Equal(t, "testerName", actualFieldComment.Comments[0].CreatedByName)
	})
}
