package go2js

type jsConnAddr struct {
	onString string
}

func (a jsConnAddr) Network() string { return "go2js" }
func (a jsConnAddr) String() string  { return a.onString }
