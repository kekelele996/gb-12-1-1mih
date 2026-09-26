package service

import (
	"testing"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestMissingRequiredMaterials(t *testing.T) {
	items := []model.MaterialItem{
		{Name: "成绩单", IsRequired: true, Status: constants.MaterialApproved},
		{Name: "推荐信", IsRequired: true, Status: constants.MaterialUploaded},
		{Name: "护照扫描件", IsRequired: true, Status: constants.MaterialPending},
		{Name: "获奖证书", IsRequired: false, Status: constants.MaterialPending},
	}
	missing := MissingRequiredMaterials(items)
	if len(missing) != 2 {
		t.Fatalf("missing = %v, want 2 items", missing)
	}
	if missing[0] != "推荐信" || missing[1] != "护照扫描件" {
		t.Errorf("missing = %v, want [推荐信 护照扫描件]", missing)
	}
}

func TestMissingRequiredMaterialsAllApproved(t *testing.T) {
	items := []model.MaterialItem{
		{Name: "成绩单", IsRequired: true, Status: constants.MaterialApproved},
		{Name: "获奖证书", IsRequired: false, Status: constants.MaterialPending},
	}
	if missing := MissingRequiredMaterials(items); len(missing) != 0 {
		t.Errorf("missing = %v, want empty", missing)
	}
}
