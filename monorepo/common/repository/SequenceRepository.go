package repository

import (
	"fmt"

	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"gorm.io/gorm"
)

type SequenceRepository interface {
	SequenceGetNextVal(seqName string) (*int64, error)
	SequenceList01() ([]models.Sequence, error)
	UpdateSequence(seq models.Sequence) error
}

type sequneceRepository struct {
	DB *gorm.DB
}

func (r sequneceRepository) SequenceGetNextVal(seqName string) (*int64, error) {
	var nextVal int64
	var sequnece models.Sequence
	r0 := r.DB.Table("SEQUENCE").Where("SEQ_GEN=?", seqName).Find(&sequnece)
	if r0 != nil && r0.RowsAffected == 0 {
		sequnece = models.Sequence{
			SEQ_GEN: seqName,
		}
		nextVal = 1
	} else {
		nextVal = sequnece.SEQ_COUNT + 1
	}
	sequnece.SEQ_COUNT = nextVal
	r1 := r.DB.Table("SEQUENCE").Save(&sequnece)
	if r1 != nil {
		if r1.Error != nil {
			return nil, r1.Error
		}
	}
	return &nextVal, nil
}

func (r sequneceRepository) SequenceList01() ([]models.Sequence, error) {
	var selectedData []models.Sequence
	res := r.DB.Table("SEQUENCE").Find(&selectedData)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		logger.WARN(fmt.Sprintf("DOES NOT EXIST DATA IN SEQUENCE TABLE"))
		return nil, nil
	}
	return selectedData, nil
}

func (r sequneceRepository) UpdateSequence(seq models.Sequence) error {
	res := r.DB.Table("SEQUENCE").Where("SEQ_GEN=?", seq.SEQ_GEN).Update("SEQ_COUNT", seq.SEQ_COUNT)
	if res.Error != nil {
		return res.Error
	}
	logger.WARN(fmt.Sprintf("ROWS AFFECTED %d", res.RowsAffected))
	return nil
}
