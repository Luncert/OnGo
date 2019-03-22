package wind

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Annotation ...
type Annotation string

// AnnotationHandler ...
type AnnotationHandler func(name string, _type reflect.Type, value reflect.Value) error

// BeanDefinition ...
type BeanDefinition struct {
	Name      string
	beanType  reflect.Type
	beanValue interface{}
}

// BeanFactory ...
type BeanFactory struct {
	beans        map[string]BeanDefinition
	annoHandlers map[Annotation]AnnotationHandler
}

// CreateBeanFactory ...
func CreateBeanFactory() *BeanFactory {
	bf := &BeanFactory{
		beans:        make(map[string]BeanDefinition),
		annoHandlers: make(map[Annotation]AnnotationHandler),
	}
	bf.RegisterAnnotationHandler(Annotation("Autowired"),
		func(name string, _type reflect.Type, field reflect.Value) (err error) {
			var bean interface{}
			if bean, err = bf.GetBean(name); err == nil {
				switch _type.Kind() {
				case reflect.String:
					if b, ok := bean.(string); ok {
						field.SetString(b)
					}
				}
			}
			return
		})
	return bf
}

// RegisterAnnotationHandler ...
func (bf *BeanFactory) RegisterAnnotationHandler(anno Annotation, handler AnnotationHandler) (err error) {
	_, ok := bf.annoHandlers[anno]
	if ok {
		err = errors.New("annotation existed")
	} else {
		bf.annoHandlers[anno] = handler
	}
	return
}

// RegisterBean ...
func (bf *BeanFactory) RegisterBean(beanName string, any interface{}) (err error) {
	_type := reflect.TypeOf(any)
	if len(beanName) == 0 {
		beanName = strings.ToLower(_type.Name())
	} else {
		beanName = strings.ToLower(beanName)
	}

	kind := _type.Kind()
	if kind == reflect.Struct {
		bf.beans[beanName] = BeanDefinition{beanName, _type, nil}
	} else if kind == reflect.String {
		bf.beans[beanName] = BeanDefinition{beanName, _type, any}
	}
	return
}

// GetBean is case insensitive
func (bf *BeanFactory) GetBean(name string) (ins interface{}, err error) {
	name = strings.ToLower(name)

	if bean, ok := bf.beans[name]; ok {

		kind := bean.beanType.Kind()
		if kind == reflect.Struct {
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
							if handler, ok := bf.annoHandlers[Annotation(anno)]; ok {
								if err := handler(f.Name, f.Type, e.Field(i)); err != nil {
									return nil, err
								}
							} else {
								fmt.Println("No handler found for annotation", anno)
							}
						}
					}
				}
			}
			ins = e.Interface()
		} else if kind == reflect.String {
			ins = bean.beanValue
		}
	} else {
		err = fmt.Errorf("No bean named with \"%s\"", name)
	}
	return ins, err
}

// GetBeanDefinition ...
func (bf *BeanFactory) GetBeanDefinition(name string) (beanDef BeanDefinition, ok bool) {
	beanDef, ok = bf.beans[name]
	return
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
