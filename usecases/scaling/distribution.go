package scaling

import (
	"fmt"

	"github.com/semi-technologies/weaviate/usecases/sharding"
)

type (
	// ShardDist shard distribution over nodes
	ShardDist map[string][]string
	// nodeShardDist map a node its shard distribution
	nodeShardDist map[string]ShardDist
)

// distributions returns shard distribution for local node as well as remote nodes
func distributions(before, after *sharding.State) (ShardDist, nodeShardDist) {
	localDist := make(ShardDist, len(before.Physical))
	nodeDist := make(map[string]ShardDist)
	for name := range before.Physical {
		newNodes := difference(after.Physical[name].BelongsToNodes, before.Physical[name].BelongsToNodes)
		if before.IsShardLocal(name) {
			localDist[name] = newNodes
		} else {
			belongsTo := before.Physical[name].BelongsToNode()
			dist := nodeDist[belongsTo]
			if dist == nil {
				dist = make(map[string][]string)
				nodeDist[belongsTo] = dist
			}
			dist[name] = newNodes
		}
	}
	return localDist, nodeDist
}

// nodes return node names
func (m nodeShardDist) nodes() []string {
	ns := make([]string, 0, len(m))
	for node := range m {
		ns = append(ns, node)
	}
	return ns
}

// hosts resolve node names into host addresses
func hosts(nodes []string, resolver clusterState) ([]string, error) {
	hs := make([]string, len(nodes))
	for i, node := range nodes {
		host, ok := resolver.NodeHostname(node)
		if !ok {
			return nil, fmt.Errorf("%w, %q", ErrUnresolvedName, node)
		}
		hs[i] = host
	}
	return hs, nil
}

// shards return names of all shards
func (m ShardDist) shards() []string {
	ns := make([]string, 0, len(m))
	for node := range m {
		ns = append(ns, node)
	}
	return ns
}


