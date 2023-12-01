//@program: cumulus-go-alfred
//@author: YL
//@create: 2023-12-01 17:08
//@desc
/*

other:
*/

package milkman_udp

import (
	"context"
	"fmt"
	"net"
	trans "tm-200/backend/lib/go_trans/ifaces"
)

type (
	UDPClient struct {
		addr   string
		conn   *net.UDPConn
		ctx    context.Context
		cancel context.CancelFunc
		isOpen bool

		txChan chan []byte         // 发送信道
		rxChan chan trans.IContext // 接收信道
	}
)

func NewUDPClient(addr string, c context.Context) *UDPClient {
	_, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithCancel(c)
	return &UDPClient{
		addr:   addr,
		ctx:    ctx,
		cancel: cancel,
		txChan: make(chan []byte),
		rxChan: make(chan trans.IContext, 1),
	}
}

func (c *UDPClient) Start() error {
	addr, err := net.ResolveUDPAddr("udp", c.addr)
	if err != nil {
		return err
	}

	raddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8080,
	}
	c.conn, err = net.DialUDP("udp", addr, raddr)
	if err != nil {
		return err
	}
	c.isOpen = true
	go c.Read()
	return nil
}

func (c *UDPClient) Stop() {
	c.cancel()
	c.isOpen = false
	_ = c.conn.Close()
	close(c.rxChan)
	close(c.txChan)
}

func (c *UDPClient) Read() {
	var data [512]byte
	for {
		n, addr, err := c.conn.ReadFromUDP(data[:])
		if err != nil {
			break
		}
		if addr == c.conn.LocalAddr() {
			continue
		}
		ctx := c.NewUDPContext()
		ctx.data = data[:n]
		ctx.remote = addr.String()
		c.rxChan <- ctx
		ctx = nil
	}
}

func (c *UDPClient) Write(data []byte, addr *net.UDPAddr) error {
	if _, err := c.conn.WriteToUDP(data, addr); err != nil {
		return err
	} else {
		return nil
	}
}

func (c *UDPClient) Send(data []byte, remote string) error {
	if data == nil {
		return fmt.Errorf("data nil")
	}

	_, err := c.conn.Write(data)
	return err
}

func (c *UDPClient) NewUDPContext() *UDPCContext {
	return NewCContext()
}

func (c *UDPClient) Receive() <-chan trans.IContext {
	return c.rxChan
}

func (c *UDPClient) GetContext() trans.IContext {
	if c.conn == nil {
		return nil
	}

	return NewCContext()
}
