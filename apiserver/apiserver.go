package main

import (
	"fmt"
	"net/http"

	"github.com/ICKelin/cframe/pkg/access"
	"github.com/ICKelin/cframe/pkg/edgemanager"
	log "github.com/ICKelin/cframe/pkg/logs"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	addr string
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{
		addr: addr,
	}
}

func (s *ApiServer) Run() {
	eng := gin.New()
	eng.POST("/api-service/v1/edge/add", s.addEdge)
	eng.DELETE("/api-service/v1/edge/del", s.delEdge)
	eng.GET("/api-service/v1/edge/list", s.getEdgeList)
	eng.GET("/api-service/v1/topology", s.getTopology)

	eng.POST("/api-service/v1/access/add", s.addAccess)
	eng.DELETE("/api-service/v1/access/del", s.delAccess)
	eng.GET("/api-service/v1/access/list", s.getAccessList)

	eng.Run(s.addr)
}

func (s *ApiServer) addEdge(ctx *gin.Context) {
	addForm := AddEdgeForm{}
	err := ctx.BindJSON(&addForm)
	if err != nil {
		log.Error("bind add edge form fail: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if len(addForm.Name) <= 0 {
		log.Error("invalid name", addForm.Name)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid name"))
		return
	}

	// verify cidr format and conflict
	ok := edgemanager.VerifyCidr(addForm.Cidr)
	if !ok {
		log.Error("verify cidr fail")
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid cidr"))
		return
	}

	edg := &edgemanager.Edge{
		Type:     addForm.Type,
		Name:     addForm.Name,
		HostAddr: addForm.HostAddr,
		Cidr:     addForm.Cidr,
	}
	edgemanager.AddEdge(edg.Name, edg)
	ctx.JSON(http.StatusOK, nil)
}

func (s *ApiServer) delEdge(ctx *gin.Context) {
	delForm := DeleteEdgeForm{}
	err := ctx.BindJSON(&delForm)
	if err != nil {
		log.Error("bind add edge form fail: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if len(delForm.Name) <= 0 {
		log.Error("invalid name", delForm.Name)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid name"))
		return
	}

	edgemanager.DelEdge(delForm.Name)
	ctx.JSON(http.StatusOK, nil)
}

func (s *ApiServer) getEdgeList(ctx *gin.Context) {
	edges := edgemanager.GetEdges()
	ctx.JSON(http.StatusOK, edges)
}

type topology struct {
	EdgeNode []*edgemanager.Edge     `json:"edge_node"`
	EdgeHost []*edgemanager.EdgeHost `json:"edge_host"`
}

func (s *ApiServer) getTopology(ctx *gin.Context) {
	edges := edgemanager.GetEdges()
	hosts := edgemanager.GetEdgeHosts()
	t := &topology{
		EdgeNode: edges,
		EdgeHost: hosts,
	}
	ctx.JSON(http.StatusOK, t)
}

func (s *ApiServer) addAccess(ctx *gin.Context) {
	var addAcForm AddAccessForm
	err := ctx.BindJSON(&addAcForm)
	if err != nil {
		log.Error("bind add access form fail: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// TODO: verify form
	info := &access.AccessInfo{
		CloudPlatform: access.CloudPlatform(addAcForm.Type),
		AccessKey:     addAcForm.AccessKey,
		AccessSecret:  addAcForm.AccessSecret,
	}
	access.Add(info)
	ctx.JSON(http.StatusOK, nil)
}

func (s *ApiServer) delAccess(ctx *gin.Context) {
	var delAcForm DeleteAccessForm
	err := ctx.BindJSON(&delAcForm)
	if err != nil {
		log.Error("bind del access form fail: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	access.Del(access.CloudPlatform(delAcForm.Type))
	ctx.JSON(http.StatusOK, nil)
}

func (s *ApiServer) getAccessList(ctx *gin.Context) {
	l, err := access.GetAccessList()
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}

	ctx.JSON(http.StatusOK, l)
}
