package response

type NodeInfo struct {
	ID   string `json:"id"`   // Unique node identifier (also the encryption key)
	Name string `json:"name"` // Name of the node, including client type, version, OS, custom data
	Addr string `json:"addr"`
}

type Blockchain struct {
	Height       int    `json:"height"`
	CurrentHash  string `json:"current_hash"`
	CurrentEpoch string `json:"current_epoch"`
}

// Progress progress of synchronization
type Progress struct {
	StartingBlock uint32 `json:"starting_block"`
	CurrentBlock  uint32 `json:"current_block"`
	HighestBlock  uint32 `json:"highest_block"`
	PulledStates  uint32 `json:"pulled_states"`
	KnownStates   uint32 `json:"known_states"`
}

type Node struct {
	P2PInfo    *NodeInfo   `json:"p2p_info"`
	P2PPeers   []*NodeInfo `json:"p2p_peers"`
	Blockchain *Blockchain `json:"blockchain"`
	IsSync     bool        `json:"is_sync"`
	Progress   *Progress   `json:"progress"`
}
