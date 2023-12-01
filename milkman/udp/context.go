//@program: cumulus-go-alfred
//@author: YL
//@create: 2023-12-01 17:10
//@desc
/*

other:
*/

package milkman_udp

import trans "tm-200/backend/lib/go_trans/ifaces"

type (
	UDPSContext struct {
		//Client *net.UDPConn
		data   []byte
		remote string
	}
)

func NewSContext() *UDPSContext {
	ctx := UDPSContext{
		//Client: c,
	}
	return &ctx
}

func (c *UDPSContext) NewContext() trans.IContext {
	ctx := NewSContext()
	return ctx
}

func (c *UDPSContext) GetRemoteAddr() string {
	return c.remote
}

func (c *UDPSContext) SetRemoteAddr(a string) {
	c.remote = a
}

func (c *UDPSContext) GetData() []byte {
	return c.data
}

type (
	UDPCContext struct {
		//Client *net.UDPConn
		data   []byte
		remote string
	}
)

func NewCContext() *UDPCContext {
	ctx := UDPCContext{
		//Client: c,
	}
	return &ctx
}

func (c *UDPCContext) NewContext() trans.IContext {
	ctx := NewCContext()
	return ctx
}

func (c *UDPCContext) SetRemoteAddr(a string) {
	c.remote = a
}

func (c *UDPCContext) GetRemoteAddr() string {
	return c.remote
}

func (c *UDPCContext) GetData() []byte {
	return c.data
}
