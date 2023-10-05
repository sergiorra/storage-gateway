package services

import (
	"context"

	"storage-gateway/domain/models"
)

type PutObjectService struct {
	nps *NodePoolService
}

func NewPutObjectService(nps *NodePoolService) *PutObjectService {
	return &PutObjectService{
		nps: nps,
	}
}

func (pos *PutObjectService) PutObject(ctx context.Context, obj *models.Object) error {
	if !obj.ID.IsValidID() {
		return models.ErrObjectIDNotValid
	}

	objectStorageNode, err := pos.nps.GetNode(obj.ID.Value())
	if err != nil {
		return err
	}

	return objectStorageNode.PutObject(ctx, obj)
}
