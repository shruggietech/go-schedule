package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// CreateJoin creates one parallel block without changing task inputs.
func (c *Client) CreateJoin(ctx context.Context, req server.JoinCreateRequest) (domain.ParallelJoin, error) {
	req.Branches = writableJoinBranches(req.Branches)
	var out domain.ParallelJoin
	err := c.do(ctx, http.MethodPost, "/v1/joins", req, &out)
	return out, err
}

// UpdateJoin atomically changes a block with its optional revision precondition.
func (c *Client) UpdateJoin(ctx context.Context, id string, req server.JoinUpdateRequest) (domain.ParallelJoin, error) {
	if req.Branches != nil {
		branches := writableJoinBranches(*req.Branches)
		req.Branches = &branches
	}
	var out domain.ParallelJoin
	err := c.do(ctx, http.MethodPatch, "/v1/joins/"+url.PathEscape(id), req, &out)
	return out, err
}

// DeleteJoin removes a block and resolves outstanding work through the daemon.
func (c *Client) DeleteJoin(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/joins/"+url.PathEscape(id), nil, nil)
}
func joinQuery(offset, limit, branchOffset, branchLimit, textLimit int) url.Values {
	q := url.Values{}
	for k, v := range map[string]int{"offset": offset, "limit": limit, "branch_offset": branchOffset, "branch_limit": branchLimit, "text_limit": textLimit} {
		if v > 0 {
			q.Set(k, strconv.Itoa(v))
		}
	}
	return q
}

// ListJoinsPage retrieves bounded definitions and nested branches.
func (c *Client) ListJoinsPage(ctx context.Context, offset, limit, branchOffset, branchLimit, textLimit int) ([]domain.ParallelJoin, bool, error) {
	var out struct {
		Joins   []domain.ParallelJoin `json:"joins"`
		HasMore bool                  `json:"has_more"`
	}
	err := c.do(ctx, http.MethodGet, withQuery("/v1/joins", joinQuery(offset, limit, branchOffset, branchLimit, textLimit)), nil, &out)
	return out.Joins, out.HasMore, err
}

// GetJoinPage retrieves one bounded exact definition page.
func (c *Client) GetJoinPage(ctx context.Context, id string, branchOffset, branchLimit, textLimit int) (domain.ParallelJoin, error) {
	var out domain.ParallelJoin
	err := c.do(ctx, http.MethodGet, withQuery("/v1/joins/"+url.PathEscape(id), joinQuery(0, 0, branchOffset, branchLimit, textLimit)), nil, &out)
	return out, err
}

// ListJoins returns complete definitions from bounded daemon pages.
func (c *Client) ListJoins(ctx context.Context) ([]domain.ParallelJoin, error) {
	out := []domain.ParallelJoin{}
	for offset := 0; ; {
		page, more, err := c.ListJoinsPage(ctx, offset, 100, 0, 100, 0)
		if err != nil {
			return nil, err
		}
		for _, v := range page {
			if v.HasMoreBranches {
				v, err = c.GetJoin(ctx, v.ID)
				if err != nil {
					return nil, err
				}
			}
			out = append(out, v)
		}
		if !more {
			return out, nil
		}
		if len(page) == 0 {
			return nil, fmt.Errorf("join definition page made no progress")
		}
		offset += len(page)
	}
}

// GetJoin returns every branch from a consistent definition revision.
func (c *Client) GetJoin(ctx context.Context, id string) (domain.ParallelJoin, error) {
	out, err := c.GetJoinPage(ctx, id, 0, 100, 0)
	if err != nil {
		return out, err
	}
	for out.HasMoreBranches {
		page, e := c.GetJoinPage(ctx, id, len(out.Branches), 100, 0)
		if e != nil {
			return out, e
		}
		if !page.UpdatedAt.Equal(out.UpdatedAt) {
			return out, &StatusError{Code: server.CodeConflict, Message: "join changed while reading branches; refresh the definition"}
		}
		if len(page.Branches) == 0 {
			return out, fmt.Errorf("join branch page made no progress")
		}
		out.Branches = append(out.Branches, page.Branches...)
		out.HasMoreBranches = page.HasMoreBranches
	}
	return out, nil
}

// ListJoinInstancesPage reads a bounded instance collection.
func (c *Client) ListJoinInstancesPage(ctx context.Context, f domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, bool, error) {
	q := joinQuery(f.Offset, f.Limit, f.BranchOffset, f.BranchLimit, f.TextLimit)
	if f.DefinitionID != "" {
		q.Set("join_id", f.DefinitionID)
	}
	if f.EntryRunID != "" {
		q.Set("entry_run_id", f.EntryRunID)
	}
	if f.State != "" {
		q.Set("state", string(f.State))
	}
	var out struct {
		Instances []domain.ParallelJoinInstance `json:"instances"`
		HasMore   bool                          `json:"has_more"`
	}
	err := c.do(ctx, http.MethodGet, withQuery("/v1/join-instances", q), nil, &out)
	return out.Instances, out.HasMore, err
}

// ListJoinInstances returns a bounded page with complete exact branch evidence.
func (c *Client) ListJoinInstances(ctx context.Context, f domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, error) {
	out, _, err := c.ListJoinInstancesPage(ctx, f)
	if err != nil {
		return nil, err
	}
	for i := range out {
		if out[i].HasMoreBranches {
			out[i], err = c.GetJoinInstance(ctx, out[i].ID)
			if err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}

// GetJoinInstancePage reads exact bounded evidence without task-history inference.
func (c *Client) GetJoinInstancePage(ctx context.Context, id string, f domain.ParallelJoinInstanceFilter) (domain.ParallelJoinInstance, error) {
	var out domain.ParallelJoinInstance
	err := c.do(ctx, http.MethodGet, withQuery("/v1/join-instances/"+url.PathEscape(id), joinQuery(0, 0, f.BranchOffset, f.BranchLimit, f.TextLimit)), nil, &out)
	return out, err
}

// GetJoinInstance returns all branches from a consistent exact instance revision.
func (c *Client) GetJoinInstance(ctx context.Context, id string) (domain.ParallelJoinInstance, error) {
	out, err := c.GetJoinInstancePage(ctx, id, domain.ParallelJoinInstanceFilter{BranchLimit: 100})
	if err != nil {
		return out, err
	}
	for out.HasMoreBranches {
		page, e := c.GetJoinInstancePage(ctx, id, domain.ParallelJoinInstanceFilter{BranchOffset: len(out.Branches), BranchLimit: 100})
		if e != nil {
			return out, e
		}
		if !page.UpdatedAt.Equal(out.UpdatedAt) {
			return out, &StatusError{Code: server.CodeConflict, Message: "join instance changed while reading branches; refresh the instance"}
		}
		if len(page.Branches) == 0 {
			return out, fmt.Errorf("join instance branch page made no progress")
		}
		out.Branches = append(out.Branches, page.Branches...)
		out.HasMoreBranches = page.HasMoreBranches
	}
	return out, nil
}

func writableJoinBranches(values []domain.ParallelJoinBranch) []domain.ParallelJoinBranch {
	out := make([]domain.ParallelJoinBranch, len(values))
	for i, v := range values {
		out[i] = domain.ParallelJoinBranch{ID: v.ID, Name: v.Name, TaskID: v.TaskID}
	}
	return out
}
