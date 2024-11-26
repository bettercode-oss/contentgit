package app

import (
	"contentgit/domain/content/projections"
	"contentgit/ports/out/persistance/eventsourcing"
	"log"
)

func (a *App) migrateDatabase() error {
	log.Println(">>> Database Migrate")
	// 테이블 생성
	if err := a.gormDB.AutoMigrate(&eventsourcing.Event{}, &eventsourcing.Snapshot{},
		&projections.ContentProjection{}, &projections.ContentFieldChange{}, &projections.ContentFieldComment{}); err != nil {
		return err
	}

	return nil
}
