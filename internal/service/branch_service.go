package service

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type BranchService struct {
	BranchRepository *repository.BranchRepository
}

func NewBranchService(db *sql.DB) *BranchService {
	return &BranchService{BranchRepository: repository.NewBranchRepository(db)}
}

func (bs *BranchService) AddBranch(cinema_chain_id int, name, address, imageURL, imageID string) (int64, error) {
	return bs.BranchRepository.AddBranch(cinema_chain_id, name, address, imageURL, imageID)
}

func (bs *BranchService) GetAllBranch() ([]model.Branch, error) {
	return bs.BranchRepository.GetAllBranch()
}

func (bs *BranchService) GetBranchByID(branchID int) (*model.Branch, error) {
	return bs.BranchRepository.GetBranchByID(branchID)
}

func (bs *BranchService) GetBranchesByCinemaChainID(cinema_chain_id int) ([]model.Branch, error) {
	return bs.BranchRepository.GetBranchesByCinemaID(cinema_chain_id)
}

func (bs *BranchService) UpdateBranch(branchID int, name, address, imageURL, imageID string) error {
	return bs.BranchRepository.UpdateBranch(branchID, name, address, imageURL, imageID)
}

func (bs *BranchService) DeleteBranch(branchID int) error {
	return bs.BranchRepository.DeleteBranch(branchID)
}
