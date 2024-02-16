package go2js

import (
	"fmt"
	"net"
	"syscall/js"
	"time"
)

// JsConn is a connection to a JavaScript runtime through WASM.
type JsConn struct {
	onRead  string
	onWrite string
	onClose *string
}

// Option is a configuration option for a JsConn.
type Option func(c *JsConn)

// WithOnReadHandler is the Option to set the onRead handler.
func WithOnReadHandler(onRead string) Option { return func(c *JsConn) { c.onRead = onRead } }

// WithOnWriteHandler is the Option to set the onWrite handler.
func WithOnWriteHandler(onWrite string) Option { return func(c *JsConn) { c.onWrite = onWrite } }

// WithOnCloseHandler is the Option to set the onClose handler.
func WithOnCloseHandler(onClose string) Option { return func(c *JsConn) { c.onClose = &onClose } }

// NewJsConn returns a new
func NewJsConn(opts ...Option) (net.Conn, error) {
	conn := &JsConn{
		onRead:  "onWrite",
		onWrite: "onRead",
		onClose: nil,
	}

	for _, opt := range opts {
		opt(conn)
	}

	if _, ok := getJS(conn.onRead); !ok {
		return nil, fmt.Errorf("onRead handler function \"%s\" is undefined", conn.onRead)
	}
	if _, ok := getJS(conn.onWrite); !ok {
		return nil, fmt.Errorf("onWrite handler function \"%s\" is undefined", conn.onWrite)
	}
	if conn.onClose != nil {
		if _, ok := getJS(*conn.onClose); !ok {
			return nil, fmt.Errorf("onClose handler function \"%s\" is undefined", *conn.onClose)
		}
	}

	return conn, nil
}

// Read invokes the JavaScript onRead function, filling the buffer with data read.
func (c *JsConn) Read(b []byte) (int, error) {
	result, err := invokeJS(c.onRead)
	if err != nil {
		return 0, err
	}
	bytesRead := copy(b, []byte(result.String()))
	return bytesRead, nil
}

// Write invokes the JavaScript onWrite handler function.
func (c *JsConn) Write(b []byte) (int, error) {
	jsValue := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(jsValue, b)
	_, err := invokeJS(c.onWrite, jsValue)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Close invokes the JavaScript onClose handler function.
func (c *JsConn) Close() error {
	if c.onClose != nil {
		_, err := invokeJS(*c.onClose)
		return err
	}
	return nil
}

// LocalAddr returns a dummy address for the Go end of the connection.
func (c *JsConn) LocalAddr() net.Addr { return jsConnAddr{"golang end"} }

// RemoteAddr returns a dummy address for the JavaScript end of the connection.
func (c *JsConn) RemoteAddr() net.Addr { return jsConnAddr{"javascript end"} }

// SetDeadline is a stub (deadlines in the context of the js event loop arent straightforward to address).
func (c *JsConn) SetDeadline(deadline time.Time) error { return nil }

// SetReadDeadline is a stub (deadlines in the context of the js event loop arent straightforward to address).
func (c *JsConn) SetReadDeadline(deadline time.Time) error { return nil }

// SetWriteDeadline is a stub (deadlines in the context of the js event loop arent straightforward to address).
func (c *JsConn) SetWriteDeadline(deadline time.Time) error { return nil }
