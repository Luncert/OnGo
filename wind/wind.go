package wind

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Bean ...
type BeanDefinition struct {
	Name     string
	beanType reflect.Type
}

// BeanFactory ...
type BeanFactory struct {
	beans map[string]BeanDefinition
}

// CreateBeanFactory ...
func CreateBeanFactory() *BeanFactory {
	return &BeanFactory{beans: make(map[string]Bean)}
}

// RegisterBean ...
func (bf *BeanFactory) RegisterBean(t reflect.Type) (err error) {
	if t.Kind() != reflect.Struct {
		err = errors.New("argument must be a Struct Type")
	} else {
		beanName := t.Name()
		bf.beans[beanName] = BeanDefinition{Name: beanName, beanType: t}
	}
	return
}

// GetBean ...
func (bf *BeanFactory) GetBean(name string) (ins interface{}, ok bool) {
	bean, ok := bf.beans[name]
	if ok {
		// reflect.New返回的是一个ptr类型的Value，可以用reflect.TypeOf去验证
		// 这里调用Elem()可以理解为对一个指针的取值操作
		e := reflect.New(bean.beanType).Elem()
		// 根据BeanType去扫描Bean的所有字段的Tag
		_type := bean.beanType
		if numField := _type.NumField(); numField > 0 {
			for i := 0; i < numField; i++ {
				f := _type.Field(i)
				// findAllAnnotations会从Tag中解析出所有名字为"@"的pair
				// 比如`@:"Autowired" @:"GetMapping('/get')"`会解析成["Autowired", "GetMapping('/get')"]，一个字符串数组
				if annos, ok := bf.findAllAnnotations(f.Tag); ok {
					for _, anno := range annos {
						// 这里就简单的支持一下Autowired注解
						if anno == "Autowired" {
							switch f.Type.Kind() {
							case reflect.String:
								// 直接注入固定值
								e.Field(i).SetString("Luncert")
							default:
								fmt.Println("Unknown type")
							}
						}
					}
				}
			}
		}
		ins = e.Interface()
	}
	return
}

func (bf *BeanFactory) GetBeanDefinition(name string) (BeanDefinition, ok) {
	return bf.beans[name]
}

// Find all tag pair named with '@'
func (bf *BeanFactory) findAllAnnotations(tag reflect.StructTag) (ret []string, ok bool) {
	ret = make([]string, 0)

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if "@" == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			ok = true
			ret = append(ret, value)
		}
	}
	return
}
