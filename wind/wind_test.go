package wind

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string `@:"Autowired"`
}

func (p *Person) hello() {
	fmt.Println("Hello,", p.Name)
}

func TestWind(t *testing.T) {
	bf := CreateBeanFactory()

	if err := bf.RegisterBean("name", "Luncert"); err != nil {
		t.Error(err)
	}

	if err := bf.RegisterBean("Person", Person{}); err != nil {
		t.Error(err)
	} else {
		bean, err := bf.GetBean("Person")
		if err == nil {
			if p, ok := bean.(Person); ok {
				p.hello()
			} else {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	}
}
