package handles

type TCP_handle struct {
	Src_IP   string
	Src_Port uint32
	Out_Port uint32
	Active   string
	Close    bool
}

func (tcp *TCP_handle) Get_src() string {
	return tcp.Src_IP
}
