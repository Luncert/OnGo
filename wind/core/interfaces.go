package core

import (
	"log"
)

// BeanScope defines the scope of a Bean
type BeanScope int

const (
	// ScopeSingleton is a Scope identifier for the standard singleton scope
	ScopeSingleton = iota
	// ScopePrototype is a Scope identifier for the standard prototype scope
	ScopePrototype
)

type beanDefinition interface {
	setScope(scope BeanScope)
	getScope() BeanScope
	setLazyInit(lazyInit bool)
	isLazyInit() bool
	setDependsOn(dependsOn []string)
	getDependsOn() []string
	setPrimary(primary bool)
	isPrimary() bool
	setFactoryBeanName(factoryBeanName string)
	getFactoryBeanName() string
	setFactoryMethodName(factoryMethodName string)
	getFactoryMethodName() string
	setInitMethodName(initMethodName string)
	getInitMethodName() string
	setDestroyMethodName(destroyMethodName string)
	getDestroyMethodName() string
	setDescription(description string)
	getDescription() string
	isSingleton() bool
	isPrototype() bool
}

type beanDefinitionRegistry interface {
	registerBeanDefinition(beanName string, beanDef beanDefinition) error
	removeBeanDefinition(beanName string) error
	getBeanDefinition(beanName string) error
	containsBeanDefinition(beanName string) bool
	getBeanDefinitionNames() []string
	getBeanDefinitionCount() int
	isBeanNameInUse(beanName string) bool
}

type defaultListableBeanFactory struct {
	allowBeanDefOverriding bool
	beanDefs               map[string]beanDefinition
	alreadyCreated         []string // Names of beans that have already been created at least once
}

func (factory *defaultListableBeanFactory) registerBeanDefinition(beanName string, beanDef beanDefinition) error {
	existingDef, ok := factory.beanDefs[beanName]
	if ok {
		if !factory.allowBeanDefOverriding {
			return &BeanDefinitionOverrideError{beanName, beanDef, existingDef}
		} else if beanDef != existingDef {
			// reflect.DeepEqual(beanDef, existingDef)
			log.Printf("Overriding bean definition for bean %s with a different definition: replacing [%v] with [%v]",
				beanName, existingDef, beanDef)
		} else {
			log.Printf("Overriding bean definition for bean %s with a equivalent definition: replacing [%v] with [%v]",
				beanName, existingDef, beanDef)
		}
	}
	factory.beanDefs[beanName] = beanDef
	return nil
}

func (factory *defaultListableBeanFactory) removeBeanDefinition(beanName string) error {

}
