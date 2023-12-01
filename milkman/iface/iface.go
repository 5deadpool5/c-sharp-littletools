//@program: cumulus-go-alfred
//@author: YL
//@create: 2023-12-01 17:07
//@desc
/*

other:
*/

package milkman_iface

type (
	IConfig interface {
		GetIP() string
		SetIP(ip string)
		GetPort() string
		SetPort(port int)
	}

	IContext interface {
		GetData() []byte
		SetRemoteAddr(string)
		GetRemoteAddr() string
		NewContext() IContext
	}

	ITrans interface {
		Start() error
		Stop()
		GetContext() IContext
		Receive() <-chan IContext
		Send(data []byte, remote string) error
	}
)
