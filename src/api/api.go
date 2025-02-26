package api

import (
	"context"

	"github.com/highonsemicolon/aura/src/service"
)

type API struct {
	object service.ObjectServiceInterface
}

func (a *API) DeleteObject(ctx context.Context, request DeleteObjectRequestObject) (DeleteObjectResponseObject, error) {
	return DeleteObject200Response{}, nil
}
func (a *API) CreateObject(ctx context.Context, request CreateObjectRequestObject) (CreateObjectResponseObject, error) {
	err := a.object.Create(ctx, request.Params.XUserId, request.Body.Object)
	if err != nil {
		return nil, err
	}

	return CreateObject201Response{}, nil
}

func (a *API) LookupResources(ctx context.Context, request LookupResourcesRequestObject) (LookupResourcesResponseObject, error) {
	return LookupResources200JSONResponse{}, nil
}

func (a *API) CreateRelationship(ctx context.Context, request CreateRelationshipRequestObject) (CreateRelationshipResponseObject, error) {
	return CreateRelationship201Response{}, nil
}

func (a *API) CheckPermission(ctx context.Context, request CheckPermissionRequestObject) (CheckPermissionResponseObject, error) {
	return CheckPermission200JSONResponse{}, nil
}

func (a *API) DeleteRelationship(ctx context.Context, request DeleteRelationshipRequestObject) (DeleteRelationshipResponseObject, error) {
	return DeleteRelationship200Response{}, nil
}

func (a *API) ExpandRelationships(ctx context.Context, request ExpandRelationshipsRequestObject) (ExpandRelationshipsResponseObject, error) {
	return ExpandRelationships200JSONResponse{}, nil
}
