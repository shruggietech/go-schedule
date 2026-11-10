package store

import "fmt"

// prepareHotStatements retains query plans, never admission results, for the store lifetime.
func (s *Store) prepareHotStatements() error {
	var err error
	s.completionAdmission, err = s.db.Prepare(completionAdmissionSQL)
	if err != nil {
		return fmt.Errorf("store: prepare completion admission: %w", err)
	}
	s.runInsert, err = s.db.Prepare(runInsertSQL)
	if err != nil {
		return fmt.Errorf("store: prepare run insertion: %w", err)
	}
	s.automationMatches, err = s.db.Prepare(automationMatchesSQL)
	if err != nil {
		return fmt.Errorf("store: prepare completion matches: %w", err)
	}
	s.chainDeliveryInsert, err = s.db.Prepare(chainDeliveryInsertSQL)
	if err != nil {
		return fmt.Errorf("store: prepare chain delivery insertion: %w", err)
	}
	return nil
}

const automationMatchesSQL = "SELECT id,target_task_id,0 FROM completion_chains WHERE source_task_id=? AND (on_outcome=? OR on_outcome=?) UNION ALL SELECT id,'',1 FROM parallel_joins WHERE entry_task_id=? AND ?='success'"
const chainDeliveryInsertSQL = "INSERT OR IGNORE INTO completion_deliveries(id,chain_id,source_task_id,target_task_id,source_run_id,state,created_at) VALUES(?,?,?,?,?,?,?)"
