package commenter

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentWithChildren struct {
	Comment
	Children []*CommentWithChildren `json:"children,omitempty"`
}

type Store struct {
	mu       sync.RWMutex
	nextID   int64
	byID     map[int64]*Comment
	children map[int64][]int64
}

func NewStore() *Store {
	return &Store{
		nextID:   1,
		byID:     make(map[int64]*Comment),
		children: make(map[int64][]int64),
	}
}

var ErrNotFound = errors.New("comment not found")

func (s *Store) Create(parentID *int64, content string) *Comment {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++
	c := &Comment{
		ID:        id,
		ParentID:  parentID,
		Content:   strings.TrimSpace(content),
		CreatedAt: time.Now().UTC(),
	}
	s.byID[id] = c
	if parentID != nil {
		s.children[*parentID] = append(s.children[*parentID], id)
	}
	return c
}

func (s *Store) buildTree(id int64) *CommentWithChildren {
	c := s.byID[id]
	node := &CommentWithChildren{Comment: *c}
	for _, cid := range s.children[id] {
		node.Children = append(node.Children, s.buildTree(cid))
	}
	return node
}

func (s *Store) GetTree(id int64) (*CommentWithChildren, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.byID[id]; !ok {
		return nil, ErrNotFound
	}
	return s.buildTree(id), nil
}

func (s *Store) ListTop(page, pageSize int, sortBy string) ([]*Comment, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var tops []*Comment
	for _, c := range s.byID {
		if c.ParentID == nil {
			cc := *c
			tops = append(tops, &cc)
		}
	}
	switch sortBy {
	case "created_asc":
		sort.Slice(tops, func(i, j int) bool { return tops[i].CreatedAt.Before(tops[j].CreatedAt) })
	default:
		sort.Slice(tops, func(i, j int) bool { return tops[i].CreatedAt.After(tops[j].CreatedAt) })
	}
	total := len(tops)
	start := (page - 1) * pageSize
	if start > total {
		return []*Comment{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return tops[start:end], total
}

func (s *Store) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return ErrNotFound
	}
	toDelete := []int64{}
	var collect func(int64)
	collect = func(cid int64) {
		toDelete = append(toDelete, cid)
		for _, child := range s.children[cid] {
			collect(child)
		}
	}
	collect(id)
	for _, del := range toDelete {
		delete(s.byID, del)
		delete(s.children, del)
	}
	for pid, list := range s.children {
		newList := list[:0]
		for _, v := range list {
			keep := true
			for _, d := range toDelete {
				if v == d {
					keep = false
					break
				}
			}
			if keep {
				newList = append(newList, v)
			}
		}
		s.children[pid] = newList
	}
	return nil
}

func (s *Store) Search(query string, page, pageSize int) ([]*Comment, int) {
	q := strings.ToLower(strings.TrimSpace(query))
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []*Comment
	for _, c := range s.byID {
		if strings.Contains(strings.ToLower(c.Content), q) {
			cc := *c
			res = append(res, &cc)
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i].CreatedAt.After(res[j].CreatedAt) })
	total := len(res)
	start := (page - 1) * pageSize
	if start > total {
		return []*Comment{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return res[start:end], total
}

type Service struct {
	store *Store
}

func NewService() *Service {
	return &Service{store: NewStore()}
}

func (s *Service) Create(parentID *int64, content string) *Comment {
	return s.store.Create(parentID, content)
}

func (s *Service) GetTree(id int64) (*CommentWithChildren, error) {
	return s.store.GetTree(id)
}

func (s *Service) ListTop(page, pageSize int, sortBy string) ([]*Comment, int) {
	return s.store.ListTop(page, pageSize, sortBy)
}

func (s *Service) Delete(id int64) error {
	return s.store.Delete(id)
}

func (s *Service) Search(query string, page, pageSize int) ([]*Comment, int) {
	return s.store.Search(query, page, pageSize)
}
