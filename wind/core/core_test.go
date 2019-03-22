package core

import (
	"fmt"
	"testing"
)

// 不支持方法级别的注解

type HelloController struct {
	Name    string `@:"Autowired"`
	SayWhat string `@:"Autowired"`
}

func (p *HelloController) Hello() {
	fmt.Println(p.SayWhat, p.Name)
}

func TestWind(t *testing.T) {
	bf := CreateBeanFactory()

	if err := bf.RegisterBean("name", "Luncert"); err != nil {
		t.Error(err)
	}

	if err := bf.RegisterBean("sayWhat", "Hi"); err != nil {
		t.Error(err)
	}

	if err := bf.RegisterBean("", HelloController{}); err != nil {
		t.Error(err)
	} else {
		bean, err := bf.GetBean("HelloController")
		if err == nil {
			if p, ok := bean.(HelloController); ok {
				p.Hello()
			} else {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	}
}
