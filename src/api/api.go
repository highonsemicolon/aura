package api

import (
	"context"
)

type API struct{}

func (a *API) LookupResources(ctx context.Context, x LookupResourcesRequestObject) (LookupResourcesResponseObject, error) {
	return LookupResources200JSONResponse{}, nil
}

func (a *API) CreateRelationship(ctx context.Context, x CreateRelationshipRequestObject) (CreateRelationshipResponseObject, error) {
	return CreateRelationship201Response{}, nil
}

func (a *API) CheckPermission(ctx context.Context, x CheckPermissionRequestObject) (CheckPermissionResponseObject, error) {
	return CheckPermission200JSONResponse{}, nil
}

func (a *API) DeleteRelationship(ctx context.Context, x DeleteRelationshipRequestObject) (DeleteRelationshipResponseObject, error) {
	return DeleteRelationship200Response{}, nil
}

func (a *API) ExpandRelationships(ctx context.Context, x ExpandRelationshipsRequestObject) (ExpandRelationshipsResponseObject, error) {
	return ExpandRelationships200JSONResponse{}, nil
}
