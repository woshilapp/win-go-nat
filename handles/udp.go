package handles

type UDP_handle struct {
	Src_IP   string
	Src_Port uint32
	Out_Port uint32
	Close    bool
}

func (udp *UDP_handle) Get_src() string {
	return udp.Src_IP
}
