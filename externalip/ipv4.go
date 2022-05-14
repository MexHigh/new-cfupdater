package externalip

type IPv4 struct {
	Addr string
}

func GetIPv4() *IPv4 {
	return &IPv4{}
}
