package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/usecase"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type CouncilNodesHandler struct {
	logger usecase.Logger

	routePath       RoutePath
	councilNodeView viewrepo.CouncilNodeViewRepo
}

func NewCouncilNodesHandler(logger usecase.Logger, routePath RoutePath, councilNodeView viewrepo.CouncilNodeViewRepo) *CouncilNodesHandler {
	return &CouncilNodesHandler{
		logger: logger.WithFields(usecase.LogFields{
			"module": "CouncilNodesHandler",
		}),

		routePath:       routePath,
		councilNodeView: councilNodeView,
	}
}

func (handler *CouncilNodesHandler) ListActiveCouncilNodes(resp http.ResponseWriter, req *http.Request) {
	var err error

	pagination, err := ParsePagination(req)
	if err != nil {
		BadRequest(resp, err)
		return
	}

	councilNodes, paginationResult, err := handler.councilNodeView.ListActivities(pagination)
	if err != nil {
		handler.logger.Errorf("error listing council nodes: %v", err)
		InternalServerError(resp)
		return
	}

	SuccessWithPagination(resp, councilNodes, paginationResult)
}

func (handler *CouncilNodesHandler) FindCouncilNodeById(resp http.ResponseWriter, req *http.Request) {
	var err error

	routeVars := handler.routePath.Vars(req)
	councilNodeIdVar, ok := routeVars["id"]
	if !ok {
		BadRequest(resp, errors.New("missing council node id path parameter"))
		return
	}
	councilNodeId, err := strconv.ParseUint(councilNodeIdVar, 10, 64)
	if err != nil {
		BadRequest(resp, errors.New("invalid council node id path parameter"))
		return
	}

	councilNode, err := handler.councilNodeView.FindById(councilNodeId)
	if err != nil {
		if err == adapter.ErrNotFound {
			NotFound(resp)
			return
		}
		handler.logger.Errorf("error finding council node: %v", err)
		InternalServerError(resp)
		return
	}

	Success(resp, councilNode)
}

func (handler *CouncilNodesHandler) ListCouncilNodeActivitiesById(resp http.ResponseWriter, req *http.Request) {
	var err error

	pagination, err := ParsePagination(req)
	if err != nil {
		BadRequest(resp, err)
		return
	}

	routeVars := handler.routePath.Vars(req)
	councilNodeIdVar, ok := routeVars["id"]
	if !ok {
		BadRequest(resp, errors.New("missing council node id path parameter"))
		return
	}
	councilNodeId, err := strconv.ParseUint(councilNodeIdVar, 10, 64)
	if err != nil {
		BadRequest(resp, errors.New("invalid council node id path parameter"))
		return
	}

	activities, paginationResult, err := handler.councilNodeView.ListActivitiesById(councilNodeId, pagination)
	if err != nil {
		handler.logger.Errorf("error listing council node activities: %v", err)
		InternalServerError(resp)
		return
	}

	SuccessWithPagination(resp, activities, paginationResult)
}
