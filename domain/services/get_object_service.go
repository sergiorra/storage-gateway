package services

import (
	"context"
	"storage-gateway/domain/models"
)

type GetObjectService struct {
	nps *NodePoolService
}

func NewGetObjectService(nps *NodePoolService) *GetObjectService {
	return &GetObjectService{
		nps: nps,
	}
}

func (gos *GetObjectService) GetObject(ctx context.Context, objectID models.ObjectID) (*models.Object, error) {
	if !objectID.IsValidID() {
		return nil, models.ErrObjectIDNotValid
	}

	objectStorageNode, err := gos.nps.GetNode(objectID.Value())
	if err != nil {
		return nil, err
	}

	return objectStorageNode.GetObject(ctx, objectID.Value())
}
