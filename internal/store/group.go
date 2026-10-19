package store

import (
	"fmt"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/task"
)

// ErrCycle is returned when an operation would create a group cycle.
var ErrCycle = fmt.Errorf("store: operation would create a group cycle")

// groupsByID loads all groups indexed by ID.
func (s *Store) groupsByID() (map[string]domain.Group, []domain.Group, error) {
	groups, err := s.ListGroups()
	if err != nil {
		return nil, nil, err
	}
	return task.ByID(groups), groups, nil
}

// GroupChainEnabled reports whether groupID and all of its ancestors are
// enabled. An empty groupID (ungrouped) is enabled.
func (s *Store) GroupChainEnabled(groupID string) (bool, error) {
	if groupID == "" {
		return true, nil
	}
	byID, _, err := s.groupsByID()
	if err != nil {
		return false, err
	}
	return task.ChainEnabled(groupID, byID), nil
}

// ValidateParent ensures parentID exists and that assigning it to groupID would
// not create a cycle. groupID may be empty for a not-yet-created group.
func (s *Store) ValidateParent(groupID, parentID string) error {
	if parentID == "" {
		return nil
	}
	byID, _, err := s.groupsByID()
	if err != nil {
		return err
	}
	if _, ok := byID[parentID]; !ok {
		return fmt.Errorf("store: parent group %q does not exist", parentID)
	}
	if groupID != "" && task.WouldCycle(groupID, parentID, byID) {
		return ErrCycle
	}
	return nil
}

// SetGroupParent reparents a group, rejecting cycles.
func (s *Store) SetGroupParent(groupID, parentID string) error {
	return s.UpdateGroup(groupID, nil, &parentID)
}

// UpdateGroup atomically changes name and parent intent after validating the
// hierarchy within the same transaction. Nil fields retain their current value.
func (s *Store) UpdateGroup(id string, name, parent *string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin update group: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	group, err := scanGroup(tx.QueryRow(`SELECT id,name,parent_id,enabled,created_at,updated_at FROM groups WHERE id=?`, id))
	if err != nil {
		return err
	}
	if parent != nil && *parent != "" {
		var exists, cycle bool
		if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM groups WHERE id=?)`, *parent).Scan(&exists); err != nil {
			return fmt.Errorf("store: inspect group parent: %w", err)
		}
		if !exists {
			return fmt.Errorf("store: parent group %q does not exist", *parent)
		}
		if err := tx.QueryRow(`WITH RECURSIVE ancestors(id,parent_id) AS (SELECT id,parent_id FROM groups WHERE id=? UNION SELECT g.id,g.parent_id FROM groups g JOIN ancestors a ON g.id=a.parent_id) SELECT EXISTS(SELECT 1 FROM ancestors WHERE id=?)`, *parent, id).Scan(&cycle); err != nil {
			return fmt.Errorf("store: inspect group ancestry: %w", err)
		}
		if cycle {
			return ErrCycle
		}
	}
	if name != nil {
		group.Name = *name
	}
	if parent != nil {
		group.ParentID = *parent
	}
	res, err := tx.Exec(`UPDATE groups SET name=?,parent_id=?,updated_at=? WHERE id=?`, group.Name, nullStr(group.ParentID), fmtTime(s.now().UTC()), id)
	if err := affected(res, err, "update group"); err != nil {
		return err
	}
	if parent != nil {
		if err := resetNotificationConditionStatesForScope(tx, domain.NotificationScopeGroup, id); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit update group: %w", err)
	}
	return nil
}

// RenameGroup updates a group's name.
func (s *Store) RenameGroup(id, name string) error {
	return s.UpdateGroup(id, &name, nil)
}

// GroupTree returns the group hierarchy as a forest.
func (s *Store) GroupTree() ([]*task.TreeNode, error) {
	_, groups, err := s.groupsByID()
	if err != nil {
		return nil, err
	}
	return task.BuildForest(groups), nil
}
