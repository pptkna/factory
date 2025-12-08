package part

import (
	"context"
	"errors"

	"github.com/pptkna/rocket-factory/inventory/internal/model"
	repoConverter "github.com/pptkna/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/pptkna/rocket-factory/inventory/internal/repository/model"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	var repoPart repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrNotFound
		}

		return nil, err
	}

	return lo.ToPtr(repoConverter.PartToModel(repoPart)), nil
}
