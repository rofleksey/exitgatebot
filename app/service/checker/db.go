package checker

import "sync"

type CommentDatabase struct {
	m      sync.RWMutex
	mapSet map[string]map[string]struct{}
}

func newCommentDatabase() *CommentDatabase {
	return &CommentDatabase{
		mapSet: make(map[string]map[string]struct{}),
	}
}

func (d *CommentDatabase) Store(username, commentID string) {
	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.mapSet[username]; !ok {
		d.mapSet[username] = make(map[string]struct{})
	}

	d.mapSet[username][commentID] = struct{}{}
}

func (d *CommentDatabase) Has(username, commentID string) bool {
	d.m.RLock()
	defer d.m.RUnlock()

	if _, ok := d.mapSet[username]; !ok {
		return false
	}

	_, ok := d.mapSet[username][commentID]
	return ok
}

func (d *CommentDatabase) Count(username string) int {
	d.m.RLock()
	defer d.m.RUnlock()

	profileSet, ok := d.mapSet[username]
	if !ok {
		return 0
	}

	return len(profileSet)
}
