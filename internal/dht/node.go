package dht

type Node struct {
	ID       string            `json:"id"`
	PeerAddr string            `json:"peerAddr"`
	JoinAddr string            `json:"joinAddr"`
	DataDir  string            `json:"dataDir"`
	Store    map[string]string `json:"store"`
	Peers    map[string]string `json:"peers"`
}

func (n *Node) Start() {

}
