package service

import (
	"github.com/cs161079/monorepo/common/repository"
)

type SequenceService interface {
}

type sequenceService struct {
	Repo repository.SequenceRepository
}

func NewSequenceService(repo repository.SequenceRepository) SequenceService {
	return sequenceService{
		Repo: repo,
	}
}

func (s sequenceService) ZeroAllSequence() error {
	selectedData, err := s.Repo.SequenceList01()
	if err != nil {
		return err
	}
	for _, seq := range selectedData {
		seq.SEQ_COUNT = 0
		s.Repo.UpdateSequence(seq)
	}
	return nil
}
