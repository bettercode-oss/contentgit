package app

import (
	"contentgit/app/cache"
	"contentgit/appservices"
	"contentgit/domain/content"
	"contentgit/ports/out/messaging/broker/pgmq"
	"contentgit/ports/out/persistance/eventsourcing"
	"contentgit/ports/out/persistance/rdb"
)

type ComponentRegistry struct {
	components map[string]any
}

func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		components: make(map[string]any),
	}
}

func (registry *ComponentRegistry) Register(name string, component any) {
	if registry.components[name] != nil {
		return
	}
	registry.components[name] = component
}

func (registry *ComponentRegistry) Get(name string) any {
	if registry.components[name] != nil {
		return registry.components[name]
	}
	return nil
}

func (a *App) registerComponents() error {
	a.componentRegistry.Register("InMemoryCache", cache.NewInMemoryCache())

	// register repositories
	a.componentRegistry.Register("ContentProjectionRepository", &rdb.ContentProjectionRepositoryImpl{})
	a.componentRegistry.Register("EventsBus", pgmq.NewPostgresMessagingQueue())

	// register services
	contentService := appservices.NewContentService(
		eventsourcing.NewRdbEventStore(
			a.componentRegistry.components["EventsBus"].(eventsourcing.EventsBus),
			content.NewEventSerializer(),
			&eventsourcing.EventRepository{},
			&eventsourcing.SnapshotRepository{},
		),
	)
	a.componentRegistry.Register("ContentService", contentService)

	contentQuery := appservices.NewContentQuery(a.componentRegistry.components["ContentProjectionRepository"].(content.ContentProjectionRepository))
	a.componentRegistry.Register("ContentQuery", contentQuery)

	// register event handlers
	contentEventHandler := content.NewContentEventHandler(content.NewEventSerializer(), a.componentRegistry.components["ContentProjectionRepository"].(content.ContentProjectionRepository))
	a.componentRegistry.Register("ContentEventHandler", contentEventHandler)

	return nil
}
