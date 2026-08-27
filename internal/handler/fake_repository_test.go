package handler

import (
	"context"
	"sync"

	"github.com/adaravaks/URLshortener/internal/model"
	"github.com/adaravaks/URLshortener/internal/repository"
)

type fakeLinkRepository struct {
	mu     sync.Mutex
	links  map[string]*model.Link
	nextID int64
}

func newFakeLinkRepository() *fakeLinkRepository {
	return &fakeLinkRepository{
		links: make(map[string]*model.Link),
	}
}

func (f *fakeLinkRepository) Create(ctx context.Context, link *model.Link) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.nextID++
	link.ID = f.nextID
	f.links[link.ShortCode] = link
	return nil
}

func (f *fakeLinkRepository) GetByCode(ctx context.Context, code string) (*model.Link, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	link, ok := f.links[code]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return link, nil
}

func (f *fakeLinkRepository) IncrementClickCount(ctx context.Context, code string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	link, ok := f.links[code]
	if !ok {
		return repository.ErrNotFound
	}
	link.ClickCount++
	return nil
}
