package bench

import "testing"

const url = "http://localhost:8080"

var (
	cfg  = config{url: url}
	cfgP = &config{url: url}
)

type config struct {
	url string
}

func (c config) vURL() string {
	return c.url
}

func (c *config) pURL() string {
	return c.url
}

type Iv interface {
	vURL() string
}

type Ip interface {
	pURL() string
}

type service struct {
	url                string
	configByValue      config
	configByPointer    *config
	configByInterfaceV Iv
	configByInterfaceP Ip
}

func newService() *service {
	return &service{
		url:                url,
		configByValue:      cfg,
		configByPointer:    cfgP,
		configByInterfaceV: cfg,
		configByInterfaceP: cfgP,
	}
}

func (s *service) strURL() string {
	return s.url
}

func (s *service) cfgByValueURL() string {
	return s.configByValue.vURL()
}

func (s *service) cfgByPointerURL() string {
	return s.configByPointer.pURL()
}

func (s *service) cfgByInterfaceVURL() string {
	return s.configByInterfaceV.vURL()
}

func (s *service) cfgByInterfacePURL() string {
	return s.configByInterfaceP.pURL()
}

func BenchmarkByStr(b *testing.B) {
	s := newService()
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("str", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.strURL()
		}
	})
	b.Run("cfgByValue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.cfgByValueURL()
		}
	})
	b.Run("cfgByPointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.cfgByPointerURL()
		}
	})
	b.Run("cfgByInterfaceV", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.cfgByInterfaceVURL()
		}
	})
	b.Run("cfgByInterfaceP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.cfgByInterfacePURL()
		}
	})
}

//goos: darwin
//goarch: arm64
//pkg: github.com/Ayupov-Ayaz/bench/configs
//BenchmarkByStr/str-8         	1000000000	         0.7260 ns/op
//BenchmarkByStr/cfgByValue-8  	1000000000	         0.7078 ns/op
//BenchmarkByStr/cfgByPointer-8         	1000000000	         0.7215 ns/op
//BenchmarkByStr/cfgByInterfaceV-8      	266082614	         4.401 ns/op
//BenchmarkByStr/cfgByInterfaceP-8      	286566992	         4.307 ns/op
