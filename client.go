package osc

import "net"

type Client struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Address() string {
	if c.conn == nil {
		return ""
	}

	addr := c.conn.RemoteAddr()
	return addr.Network() + " - " + addr.String()
}

func (c *Client) Connect(network string, address string) error {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return err
	}
	c.addr = addr

	conn, err := net.DialUDP(network, nil, addr)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

func (c *Client) Send(packet Packet) (int, error) {
	packetBinary, _ := packet.MarshalBinary()
	return c.conn.Write(packetBinary)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
