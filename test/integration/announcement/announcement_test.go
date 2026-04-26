//go:build integration

package announcement_integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	announcementsvc "github.com/medincident/medincident-backend/internal/service/command/announcement"
)

func codeOf(t *testing.T, err error) string {
	t.Helper()
	var oe oops.OopsError
	if !errors.As(err, &oe) {
		t.Fatalf("expected oops error, got %T: %v", err, err)
	}
	code, _ := oe.Code().(string)
	return code
}

func strPtr(s string) *string { return &s }

func TestCreate_OrgLevel(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тестовое объявление",
			Content:        "Содержание объявления достаточной длины",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	a := loadAnnouncement(t, res.ID)
	assert.Equal(t, orgID, a.OrganizationID)
	assert.Equal(t, "Тестовое объявление", a.Title)
	assert.False(t, a.IsArchived)

	var vc int64
	require.NoError(t, testDB.Raw(
		`SELECT view_count FROM projections.announcement_views WHERE announcement_id = ?`, res.ID,
	).Scan(&vc).Error)
	assert.Equal(t, int64(0), vc)
}

func TestCreate_ClinicLevel(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	clinicID := insertClinic(t, orgID)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			ClinicID:       strPtr(clinicID.String()),
			Title:          "Клиническое объявление",
			Content:        "Содержание объявления достаточной длины",
			Priority:       "high",
		},
	})
	require.NoError(t, err)

	a := loadAnnouncement(t, res.ID)
	assert.True(t, a.ClinicID.Valid)
	assert.Equal(t, clinicID, a.ClinicID.UUID)
}

func TestCreate_DeptLevel(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	clinicID := insertClinic(t, orgID)
	deptID := insertDept(t, clinicID)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			ClinicID:       strPtr(clinicID.String()),
			DepartmentID:   strPtr(deptID.String()),
			Title:          "Отделенческое объявление",
			Content:        "Содержание объявления достаточной длины",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)

	a := loadAnnouncement(t, res.ID)
	assert.True(t, a.DepartmentID.Valid)
	assert.Equal(t, deptID, a.DepartmentID.UUID)
}

func TestCreate_OrgNotFound(t *testing.T) {
	resetDB(t)
	_, err := svc.Create(context.Background(), &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: uuid.New().String(),
			Title:          "Тест",
			Content:        "Содержимое длиннее десяти символов",
			Priority:       "normal",
		},
	})
	require.Error(t, err)
	assert.Equal(t, announcementsvc.ErrCodeAnnouncementOrgNotFound, codeOf(t, err))
}

func TestCreate_DeptWithoutClinic(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)
	clinicID := insertClinic(t, orgID)
	deptID := insertDept(t, clinicID)

	_, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			DepartmentID:   strPtr(deptID.String()),
			Title:          "Тест",
			Content:        "Содержимое длиннее десяти символов",
			Priority:       "normal",
		},
	})
	require.Error(t, err)
	assert.Equal(t, announcementsvc.ErrCodeAnnouncementInvalidScope, codeOf(t, err))
}

func TestCreate_InvalidTimeRange(t *testing.T) {
	resetDB(t)
	orgID := insertOrg(t)
	starts := time.Now().Add(2 * time.Hour)
	ends := time.Now().Add(1 * time.Hour)

	_, err := svc.Create(context.Background(), &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тест",
			Content:        "Содержимое длиннее десяти символов",
			Priority:       "normal",
			StartsAt:       &starts,
			EndsAt:         &ends,
		},
	})
	require.Error(t, err)
	assert.Equal(t, announcementsvc.ErrCodeAnnouncementTimeRange, codeOf(t, err))
}

func TestUpdate_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Старый заголовок",
			Content:        "Старое содержание объявления",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)

	require.NoError(t, svc.Update(ctx, &announcementsvc.UpdateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.UpdateAnnouncementPayload{
			ID:      res.ID.String(),
			Title:   "Новый заголовок",
			Content: "Новое содержание объявления",
		},
	}))

	a := loadAnnouncement(t, res.ID)
	assert.Equal(t, "Новый заголовок", a.Title)
	assert.Equal(t, "Новое содержание объявления", a.Content)
}

func TestUpdate_NotFound(t *testing.T) {
	resetDB(t)
	err := svc.Update(context.Background(), &announcementsvc.UpdateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.UpdateAnnouncementPayload{
			ID:      uuid.New().String(),
			Title:   "Тест",
			Content: "Содержание объявления",
		},
	})
	require.Error(t, err)
	assert.Equal(t, announcementsvc.ErrCodeAnnouncementNotFound, codeOf(t, err))
}

func TestUpdatePriority_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тест",
			Content:        "Содержание объявления длиннее 10",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)

	require.NoError(t, svc.UpdatePriority(ctx, &announcementsvc.UpdateAnnouncementPriorityCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.UpdateAnnouncementPriorityPayload{
			ID:       res.ID.String(),
			Priority: "high",
		},
	}))
	assert.Equal(t, "high", string(loadAnnouncement(t, res.ID).Priority))
}

func TestUpdatePriority_ArchivedRejected(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тест",
			Content:        "Содержание объявления длиннее 10",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)
	require.NoError(t, svc.Archive(ctx, &announcementsvc.ArchiveAnnouncementCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.ArchiveAnnouncementPayload{ID: res.ID.String()},
	}))

	err = svc.UpdatePriority(ctx, &announcementsvc.UpdateAnnouncementPriorityCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.UpdateAnnouncementPriorityPayload{ID: res.ID.String(), Priority: "high"},
	})
	require.Error(t, err)
	assert.Equal(t, announcementsvc.ErrCodeAnnouncementArchived, codeOf(t, err))
}

func TestArchive_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тест",
			Content:        "Содержание объявления длиннее 10",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)

	require.NoError(t, svc.Archive(ctx, &announcementsvc.ArchiveAnnouncementCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.ArchiveAnnouncementPayload{ID: res.ID.String()},
	}))
	assert.True(t, loadAnnouncement(t, res.ID).IsArchived)
}

func TestArchive_Idempotent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тест",
			Content:        "Содержание объявления длиннее 10",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)

	require.NoError(t, svc.Archive(ctx, &announcementsvc.ArchiveAnnouncementCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.ArchiveAnnouncementPayload{ID: res.ID.String()},
	}))
	require.NoError(t, svc.Archive(ctx, &announcementsvc.ArchiveAnnouncementCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.ArchiveAnnouncementPayload{ID: res.ID.String()},
	}))
}

func TestUnarchive_HappyPath(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тест",
			Content:        "Содержание объявления длиннее 10",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)
	require.NoError(t, svc.Archive(ctx, &announcementsvc.ArchiveAnnouncementCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.ArchiveAnnouncementPayload{ID: res.ID.String()},
	}))
	require.NoError(t, svc.Unarchive(ctx, &announcementsvc.UnarchiveAnnouncementCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.UnarchiveAnnouncementPayload{ID: res.ID.String()},
	}))
	assert.False(t, loadAnnouncement(t, res.ID).IsArchived)
}

func TestUnarchive_Idempotent(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	orgID := insertOrg(t)

	res, err := svc.Create(ctx, &announcementsvc.CreateAnnouncementCommand{
		Caller: sysadminCaller,
		Payload: announcementsvc.CreateAnnouncementPayload{
			OrganizationID: orgID.String(),
			Title:          "Тест",
			Content:        "Содержание объявления длиннее 10",
			Priority:       "normal",
		},
	})
	require.NoError(t, err)

	// Not archived — Unarchive must be a no-op.
	require.NoError(t, svc.Unarchive(ctx, &announcementsvc.UnarchiveAnnouncementCommand{
		Caller:  sysadminCaller,
		Payload: announcementsvc.UnarchiveAnnouncementPayload{ID: res.ID.String()},
	}))
	assert.False(t, loadAnnouncement(t, res.ID).IsArchived)
}
