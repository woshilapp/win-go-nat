package handles

type ICMP_handle struct {
	Src_IP string
	Id     int32
	Close  bool
}

func (icmp *ICMP_handle) Get_src() string {
	return icmp.Src_IP
}
