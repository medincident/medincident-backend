package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	departmentv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/department/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

type UpdateDepartmentDetailsCommand struct {
	ID          uuid.UUID
	Name        string
	Description *string
}

func buildDepartmentDetailsChangedEvent(d *model.Department) *departmentv1.DepartmentDetailsChanged {
	ev := &departmentv1.DepartmentDetailsChanged{Name: d.Name}
	if d.Description.Valid {
		desc := d.Description.String
		ev.Description = &desc
	}
	return ev
}

func (s *DepartmentService) UpdateDetails(
	ctx context.Context,
	cmd UpdateDepartmentDetailsCommand,
) error {
	var errs []error
	if err := validateDepartmentName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateDepartmentDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dept model.Department
		if err := tx.First(&dept, "id = ?", cmd.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.department").
					Code(ErrCodeDepartmentNotFound).
					Public("Department not found.").
					With("department_id", cmd.ID).
					Errorf("department not found")
			}
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentLoadFailed).
				With("department_id", cmd.ID).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Name)
		newDesc := null.StringFromPtr(cmd.Description)
		if dept.Name == newName && dept.Description == newDesc {
			return nil
		}
		dept.Name = newName
		dept.Description = newDesc
		if err := tx.Save(&dept).Error; err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentSaveFailed).
				With("department_id", cmd.ID).
				Wrap(err)
		}

		event := buildDepartmentDetailsChangedEvent(&dept)
		payload, err := anypb.New(event)
		if err != nil {
			return oops.In("services.orgstructure.department").
				Code(ErrCodeDepartmentEventBuildFailed).
				With("department_id", cmd.ID).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(dept.UpdatedAt),
			AggregateType: AggregateTypeDepartment,
			AggregateId:   dept.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectDepartmentDetailsChanged, envelope)
	})
}
