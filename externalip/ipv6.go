package externalip

type IPv6 struct {
	Addr string
}

func GetIPv6() *IPv6 {
	return &IPv6{}
}
