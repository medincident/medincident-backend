package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateDepartmentDetailsPayload carries the new name and (optional)
// description for an existing department.
type UpdateDepartmentDetailsPayload struct {
	ID          string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// UpdateDepartmentDetailsCommand = caller + payload.
type UpdateDepartmentDetailsCommand struct {
	Caller  authz.Caller
	Payload UpdateDepartmentDetailsPayload
}

func buildDepartmentDetailsChangedEnvelope(d *model.Department) (*eventv1.Envelope, error) {
	var desc string
	if d.Description.Valid {
		desc = d.Description.String
	}
	msg := &deptv1.DepartmentDetailsChanged{
		Name:        d.Name,
		Description: desc,
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.orgstructure.department").
			Code(ErrCodeDepartmentSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(d.UpdatedAt),
		AggregateType: "department",
		AggregateId:   d.ID.String(),
		Payload:       payload,
	}, nil
}

// UpdateDetails changes a department's name and description.
//
// See: docs/services/OrgStructure.md
func (s *DepartmentService) UpdateDetails(
	ctx context.Context,
	cmd UpdateDepartmentDetailsCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Department(id)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dept model.Department
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&dept, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", id).
					Errorf("department not found")
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("department_id", id).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Payload.Name)
		var newDesc null.String
		if cmd.Payload.Description != nil {
			newDesc = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if dept.Name == newName && dept.Description == newDesc {
			return nil
		}
		dept.Name = newName
		dept.Description = newDesc
		if err := tx.Save(&dept).Error; err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", id).
				Wrap(err)
		}

		env, err := buildDepartmentDetailsChangedEnvelope(&dept)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.department.v1.details_changed", env)
	})
}
