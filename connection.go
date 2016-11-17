package reimpl

import (
	"net"
	"time"

	"github.com/rpoletaev/respio"
)

type ConConfig struct {
	Password string
	DB       int
	WTimeOut time.Duration
	RTimeOut time.Duration
}

// Conn connection wrapper with RESP support
type conn struct {
	*ConConfig
	conn   net.Conn
	reader *respio.RESPReader
	writer *respio.RESPWriter
}

func (c *conn) Close() {
	c.conn.Close()
}

// Dial setup connection with default config
func Dial(host string, port string) (*conn, error) {
	defaultConfig := &ConConfig{
		DB: 0,
	}

	return DialWithConfig(host, port, defaultConfig)
}

// DialWithConfig dial server with config
func DialWithConfig(host string, port string, config *ConConfig) (*conn, error) {
	tcp, err := net.Dial("tcp", host+port)
	if err != nil {
		return nil, err
	}

	con := &conn{
		config,
		tcp,
		respio.NewReader(tcp),
		respio.NewWriter(tcp),
	}

	return con, nil
}

// Cmd send command to server and return response or error
func (c *conn) Cmd(cmd string, prs ...interface{}) (interface{}, error) {
	if cmd == "" {
		return nil, nil
	}

	if c.WTimeOut != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.WTimeOut))
	}

	err := c.writer.SendCmd(cmd, prs)
	if err != nil {
		return nil, err
	}

	err = c.writer.Flush()

	if c.RTimeOut != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.RTimeOut))
	}

	return c.reader.Read()
}
