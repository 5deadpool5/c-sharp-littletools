//@program: cumulus-go-alfred
//@author: YL
//@create: 2023-12-01 17:10
//@desc
/*

other:
*/

package milkman_udp

import (
	"context"
	"net"
	trans "tm-200/backend/lib/go_trans/ifaces"
	"tm-200/backend/utils/commrecoder"
)

type (
	UDPServer struct {
		addr   string
		conn   *net.UDPConn
		ctx    context.Context
		cancel context.CancelFunc
		isOpen bool

		txChan chan []byte         // 发送信道
		rxChan chan trans.IContext // 接收信道
	}
)

var (
	addrs []string
)

func NewUDPServer(addr string, c context.Context) *UDPServer {
	_, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithCancel(c)
	return &UDPServer{
		addr:   addr,
		ctx:    ctx,
		cancel: cancel,
		txChan: make(chan []byte),
		rxChan: make(chan trans.IContext, 1),
	}
}

func (us *UDPServer) Start() error {
	ad, err := net.ResolveUDPAddr("udp", us.addr)
	if err != nil {
		return err
	}

	us.conn, err = net.ListenUDP("udp", ad)
	if err != nil {
		return err
	}

	as, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range as {
			if ipNet, ok := a.(*net.IPNet); ok {
				if ipNet.IP.To4() != nil {
					addrs = append(addrs, ipNet.IP.String())
				}
			}
		}
	}

	us.isOpen = true
	go us.Read()
	return nil
}

func (us *UDPServer) Stop() {
	if us.isOpen == false {
		return
	}

	us.isOpen = false
	us.cancel()
	defer func() {
		_ = us.conn.Close()
		close(us.rxChan)
		close(us.txChan)
	}()
}

func (us *UDPServer) Read() {
	for {
		select {
		case <-us.ctx.Done():
			return
		default:
			var data [512]byte
			n, addr, err := us.conn.ReadFromUDP(data[:])
			if err != nil {
				break
			}

			if us.isMyIP(addr.IP.String()) == true {
				continue
			}
			ctx := us.NewUDPContext()
			ctx.data = data[:n]
			ctx.SetRemoteAddr(addr.String())
			us.rxChan <- ctx
			commrecoder.Infof("receive udp: %+ 2x", ctx.data)
		}
	}
}

func (us *UDPServer) Write(data []byte, addr *net.UDPAddr) error {
	if _, err := us.conn.WriteToUDP(data, addr); err != nil {
		return err
	} else {
		return nil
	}
}

func (us *UDPServer) Send(data []byte, remote string) error {
	if remote == "" {
		remote = "255.255.255.255:8080"
	}
	rAddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return err
	}
	if _, err := us.conn.WriteToUDP(data, rAddr); err != nil {
		commrecoder.Infof("send udp error:%s", err.Error())
		return err
	} else {
		commrecoder.Infof("send udp: %+ 2x", data)
		return nil
	}
}

func (c *UDPServer) NewUDPContext() *UDPSContext {
	return NewSContext()
}

func (us *UDPServer) Receive() <-chan trans.IContext {
	return us.rxChan
}

func (us *UDPServer) GetContext() trans.IContext {
	if us.conn == nil {
		return nil
	}

	return NewSContext()
}

func (us *UDPServer) isMyIP(addr string) bool {
	if addrs == nil {
		return false
	}
	for _, a := range addrs {
		if a == addr {
			return true
		}
	}
	return false
}
