package stat

import (
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/remap"
	"github.com/apache/incubator-trafficcontrol/grove/web"
)

func StatsInc(m *web.ConnMap, num int, addrs *[]string) {
	for i := 0; i < num; i++ {
		c := NewFakeConn()
		*addrs = append(*addrs, c.RemoteAddr().String())
		m.Add(c)
	}
}

func StatsDec(m *web.ConnMap, num int, addrs *[]string) {
	for i := 0; i < num; i++ {
		if len(*addrs) == 0 {
			return
		}
		m.Remove((*addrs)[0])
		*addrs = (*addrs)[1:]
	}
}

type FakeConn struct{ Addr net.Addr }

func (c FakeConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (c FakeConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (c FakeConn) Close() error                       { return nil }
func (c FakeConn) LocalAddr() net.Addr                { return c.Addr }
func (c FakeConn) RemoteAddr() net.Addr               { return c.Addr }
func (c FakeConn) SetDeadline(t time.Time) error      { return nil }
func (c FakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (c FakeConn) SetWriteDeadline(t time.Time) error { return nil }

type FakeAddr struct {
	addr    string
	network string
}

func (a FakeAddr) Network() string { return a.network }
func (a FakeAddr) Addr() string    { return a.addr }
func (a FakeAddr) String() string  { return a.addr }

func NewFakeConn() net.Conn {
	s := GenGUIDStr()
	a := FakeAddr{addr: s, network: s}
	return FakeConn{Addr: net.Addr(&a)}
}

func GenGUIDStr() string {
	length := 32 // 32 characters ought to be enough for anyone
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_")
	s := make([]rune, length)
	for i := range s {
		s[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(s)
}

func TestStatsCount(t *testing.T) {
	{
		httpConns := web.NewConnMap()
		httpsConns := web.NewConnMap()
		addrs := []string{}
		r := remap.RemapRule{RemapRuleBase: remap.RemapRuleBase{Name: "foo"}}
		stats := New([]remap.RemapRule{r}, nil, 0, httpConns, httpsConns)
		expected := 10
		StatsInc(httpConns, expected, &addrs)
		if actual := stats.Connections(); actual != uint64(expected) {
			t.Errorf("Stats.Connections() expected %v actual %v", expected, actual)
		}
	}
	{
		httpConns := web.NewConnMap()
		httpsConns := web.NewConnMap()
		addrs := []string{}
		r := remap.RemapRule{RemapRuleBase: remap.RemapRuleBase{Name: "foo"}}
		stats := New([]remap.RemapRule{r}, nil, 0, httpConns, httpsConns)
		expected := 10
		StatsInc(httpsConns, expected, &addrs)
		if actual := stats.Connections(); actual != uint64(expected) {
			t.Errorf("Stats.Connections() expected %v actual %v", expected, actual)
		}
	}
	{
		httpConns := web.NewConnMap()
		httpsConns := web.NewConnMap()
		addrs := []string{}
		r := remap.RemapRule{RemapRuleBase: remap.RemapRuleBase{Name: "foo"}}
		stats := New([]remap.RemapRule{r}, nil, 0, httpConns, httpsConns)
		expected := 10
		StatsInc(httpConns, expected, &addrs)
		StatsInc(httpsConns, expected, &addrs)
		if actual := stats.Connections(); actual != uint64(expected)*2 {
			t.Errorf("Stats.Connections() expected %v actual %v", expected, actual)
		}
	}

	{
		httpConns := web.NewConnMap()
		httpsConns := web.NewConnMap()
		addrs := []string{}
		r := remap.RemapRule{RemapRuleBase: remap.RemapRuleBase{Name: "foo"}}
		stats := New([]remap.RemapRule{r}, nil, 0, httpConns, httpsConns)
		count := 10
		StatsInc(httpConns, count, &addrs)
		StatsDec(httpConns, count, &addrs)
		if actual := stats.Connections(); actual != 0 {
			t.Errorf("Stats.Connections() expected %v actual %v", 0, actual)
		}
	}
	{
		httpConns := web.NewConnMap()
		httpsConns := web.NewConnMap()
		addrs := []string{}
		r := remap.RemapRule{RemapRuleBase: remap.RemapRuleBase{Name: "foo"}}
		stats := New([]remap.RemapRule{r}, nil, 0, httpConns, httpsConns)
		count := 10
		StatsInc(httpsConns, count, &addrs)
		StatsDec(httpsConns, count, &addrs)
		if actual := stats.Connections(); actual != 0 {
			t.Errorf("Stats.Connections() expected %v actual %v", 0, actual)
		}
	}
	{
		httpConns := web.NewConnMap()
		httpsConns := web.NewConnMap()
		addrs := []string{}
		r := remap.RemapRule{RemapRuleBase: remap.RemapRuleBase{Name: "foo"}}
		stats := New([]remap.RemapRule{r}, nil, 0, httpConns, httpsConns)
		count := 10
		StatsInc(httpConns, count, &addrs)
		StatsInc(httpsConns, count, &addrs)
		StatsDec(httpConns, count, &addrs)
		if actual := stats.Connections(); actual != uint64(count) {
			t.Errorf("Stats.Connections() expected %v actual %v", count, actual)
		}
	}

	{
		httpConns := web.NewConnMap()
		httpsConns := web.NewConnMap()
		addrs := []string{}
		r := remap.RemapRule{RemapRuleBase: remap.RemapRuleBase{Name: "foo"}}
		stats := New([]remap.RemapRule{r}, nil, 0, httpConns, httpsConns)
		count := 10
		StatsInc(httpConns, count, &addrs)
		StatsDec(httpConns, 1, &addrs)
		if actual := stats.Connections(); actual != uint64(count-1) {
			t.Errorf("stats.Connections() expected %v actual %v", count-1, actual)
		}
	}

}