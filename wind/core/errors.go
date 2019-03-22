package core

import "fmt"

// BeanDefinitionOverrideError ...
type BeanDefinitionOverrideError struct {
	beanName    string
	beanDef     beanDefinition
	existingDef beanDefinition
}

func (e *BeanDefinitionOverrideError) Error() string {
	return fmt.Sprintf("BeanDefinitionOverrideError bean name: %s, bean def: %v, existing def: %v",
		e.beanName, e.beanDef, e.existingDef)
}
