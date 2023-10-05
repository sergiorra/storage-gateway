package services

import (
	"context"
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
	"time"

	"storage-gateway/domain/models"
	"storage-gateway/domain/ports"
	"storage-gateway/internal/context-wrapper"
	"storage-gateway/internal/log"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
)

// RingNode represents a node in the object storage ring with its associated hash ID
type RingNode struct {
	Node   ports.ObjectStorage
	HashID uint32
}

// NodePoolService manages a pool of object storage nodes and provides methods for refreshing and balancing the nodes using consistent hashing
type NodePoolService struct {
	ds        ports.DiscoveryService
	scheduler *gocron.Scheduler
	nodes     []*RingNode
	mu        sync.Mutex
}

// NewNodePoolService creates a new instance of NodePoolService with the provided discovery service for node discovery
func NewNodePoolService(ds ports.DiscoveryService) *NodePoolService {
	return &NodePoolService{
		ds:        ds,
		scheduler: gocron.NewScheduler(time.UTC),
		nodes:     make([]*RingNode, 0),
	}
}

// StartRefreshingNodes starts a periodic task for refreshing the pool of nodes discovered by the discovery service and balances them
func (nps *NodePoolService) StartRefreshingNodes() error {
	_, err := nps.scheduler.Every(2).Minutes().Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
		defer cancel()

		correlationID := uuid.New().String()
		ctx = context_wrapper.WithCorrelationID(ctx, correlationID)

		log.Infot(context_wrapper.GetCorrelationID(ctx), "refreshing pool nodes")

		nodes, err := nps.ds.DiscoverNodes(ctx)
		if err != nil {
			log.Errort(context_wrapper.GetCorrelationID(ctx), fmt.Sprintf("could not discover nodes with error %s", err))
			return
		}

		nps.BalanceNodes(nodes)
	})
	if err != nil {
		return err
	}

	nps.scheduler.StartAsync()

	return nil
}

// BalanceNodes updates the pool of nodes with the provided list of new nodes and recalculates their hash IDs for load balancing using consistent hashing
func (nps *NodePoolService) BalanceNodes(newNodes []ports.ObjectStorage) {
	nps.mu.Lock()
	defer nps.mu.Unlock()

	nps.nodes = make([]*RingNode, 0)

	for _, newNode := range newNodes {
		nps.nodes = append(nps.nodes, &RingNode{
			Node:   newNode,
			HashID: crc32.ChecksumIEEE([]byte(newNode.ID())),
		})
	}

	sort.Slice(nps.nodes, func(i, j int) bool {
		return nps.nodes[i].HashID < nps.nodes[j].HashID
	})
}

// GetNode returns the object storage node responsible for the given key based on the consistent hash ring
func (nps *NodePoolService) GetNode(key string) (ports.ObjectStorage, error) {
	nps.mu.Lock()
	defer nps.mu.Unlock()

	if len(nps.nodes) == 0 {
		return nil, fmt.Errorf("no nodes in the pool")
	}

	i := sort.Search(len(nps.nodes), func(i int) bool {
		return nps.nodes[i].HashID >= crc32.ChecksumIEEE([]byte(key))
	})

	if i >= len(nps.nodes) {
		i = 0
	}

	if !nps.nodes[i].Node.IsOnline() {
		return nil, models.ErrObjectStorageNotAvailable
	}

	return nps.nodes[i].Node, nil
}

// StopRefreshingNodes stops the periodic node refreshing task
func (nps *NodePoolService) StopRefreshingNodes() {
	nps.scheduler.Stop()
}
