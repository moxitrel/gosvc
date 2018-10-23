/*

NewDispatch bufferSize:
	Set    x cb: "add handler for x.Key"
	Call   x   : "sched cb(x)"

*/
package svc

import (
	"fmt"
	"github.com/moxitrel/golib"
	"reflect"
	"sync"
)

// Cache the check result
var validateMapKeyCache = new(sync.Map)

// Check whether <keyType> is a valid key type of map
func ValidateMapKey(keyType reflect.Type) (v bool) {
	v = true

	if keyType != nil {
		switch keyType.Kind() {
		case reflect.Invalid, reflect.Func, reflect.Slice, reflect.Map:
			// shoudn't be a function, slice or map
			v = false
		case reflect.Struct: // shoudn't be a struct contains function, slice or map field
			// fetch result from cache if exists
			if anyV, ok := validateMapKeyCache.Load(keyType); ok {
				v = anyV.(bool)
			} else {
				for i := 0; i < keyType.NumField(); i++ {
					if ValidateMapKey(keyType.Field(i).Type) == false {
						v = false
						break
					}
				}
				validateMapKeyCache.Store(keyType, v)
			}
		}
	}

	return
}

/* Usage:

// 1. define a new type derive Dispatch
//
type MyHandlerService struct {
	Dispatch
}

// 2. define method Call() with arg has key() attribute
//
func (o MyHandlerService) Call(arg T) {
	o.Dispatch.Call(arg.Key(), arg)
}

*/
// Process arg by handlers[arg.key()] in goroutine pool
type Dispatch struct {
	*Func
	*Pool
	handlers *sync.Map
}

func NewDispatch(bufferSize uint, poolMin uint) (v *Dispatch) {
	v = new(Dispatch)
	v.handlers = new(sync.Map)
	v.Pool = NewPool(func(anyFunArg interface{}) {
		funArg := anyFunArg.([2]interface{})
		fun := funArg[0].(func(interface{}))
		arg := funArg[1]
		fun(arg)
	}).WithCount(poolMin, defaultPoolMax)
	v.Func = NewFuncWithSize(bufferSize, v.Pool.Apply)
	return
}

func (o *Dispatch) Stop() {
	if o.state == RUNNING {
		o.Func.Stop()
		o.Pool.WithTime(o.Pool.delay, 0)
	}
}

// key: key's type shoudn't be function, slice, map or struct contains function, slice or map field
// fun: nil, delete the handler according to key
func (o *Dispatch) Set(key interface{}, fun func(arg interface{})) {
	// assert invalid key type
	//if ValidateMapKey(reflect.TypeOf(key)) == false {
	//	golib.Panic(fmt.Sprintf("%t isn't a valid map key type!\n", key))
	//}
	if fun == nil {
		// delete handler
		o.handlers.Delete(key)
	} else {
		o.handlers.Store(key, fun)
	}
}

func (o *Dispatch) Get(key interface{}) func(interface{}) {
	v, _ := o.handlers.Load(key)
	if v == nil {
		return nil
	}
	return v.(func(interface{}))
}

func (o *Dispatch) Apply(key interface{}, arg interface{}) {
	//if ValidateMapKey(reflect.TypeOf(key)) == false {
	//	golib.Warn(fmt.Sprintf("%t isn't a valid map key type!\n", key))
	//	return
	//}
	fun, _ := o.handlers.Load(key)
	if fun == nil {
		golib.Warn(fmt.Sprintf("%v, handler doesn't exist!\n", key))
		return
	}
	o.Func.Apply([2]interface{}{fun, arg})
}