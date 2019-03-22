package wind

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	Name string `@:"Autowired"`
}

func (p *Person) hello() {
	fmt.Println("Hello,", p.Name)
}

func TestWind(t *testing.T) {
	_type := reflect.TypeOf(Person{})
	bf := CreateBeanFactory()
	err := bf.RegisterBean(_type)
	if err != nil {
		fmt.Println(err)
	} else if bean, ok := bf.GetBean("Person"); ok {
		if p, ok := bean.(Person); ok {
			p.hello()
		}
	}
}
