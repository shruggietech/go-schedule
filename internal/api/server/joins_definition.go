package server

import "net/http"

// handleGetJoinDefinition serves complete writable text for deliberate authoring.
// Branches remain paged, while ordinary observation endpoints retain text limits.
func (s *Server) handleGetJoinDefinition(w http.ResponseWriter, r *http.Request) {
	offset, ok := nonNegativeQuery(w, r, "branch_offset", 0)
	if !ok {
		return
	}
	limit, ok := nonNegativeQuery(w, r, "branch_limit", 100)
	if !ok {
		return
	}
	if limit < 1 || limit > 1000 {
		writeError(w, http.StatusBadRequest, CodeValidation, "branch_limit", "branch_limit must be between 1 and 1000")
		return
	}
	join, err := s.store.GetParallelJoinDefinitionPageContext(r.Context(), r.PathValue("id"), offset, limit)
	if err != nil {
		s.joinError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, join)
}
