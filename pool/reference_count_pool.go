package pool

//https://www.akshaydeo.com/blog/2017/12/23/How-did-I-improve-latency-by-700-percent-using-syncPool/

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Interface following reference countable interface
// We have provided inbuilt embeddable implementation of the reference countable pool
// This interface just provides the extensibility for the implementation
type ReferenceCountable interface {
	// Method to set the current instance
	SetInstance(i interface{})
	// Method to increment the reference count
	IncrementReferenceCount()
	// Method to decrement reference count
	DecrementReferenceCount()
}

// Struct representing reference
// This struct is supposed to be embedded inside the object to be pooled
// Along with that incrementing and decrementing the references is highly important specifically around routines
type ReferenceCounter struct {
	count       *uint32                 `sql:"-" json:"-" yaml:"-"`
	destination *sync.Pool              `sql:"-" json:"-" yaml:"-"`
	released    *uint32                 `sql:"-" json:"-" yaml:"-"`
	instance    interface{}             `sql:"-" json:"-" yaml:"-"`
	reset       func(interface{}) error `sql:"-" json:"-" yaml:"-"`
	id          uint32                  `sql:"-" json:"-" yaml:"-"`
}

// Method to increment a reference
func (r ReferenceCounter) IncrementReferenceCount() {
	atomic.AddUint32(r.count, 1)
}

// Method to decrement a reference
// If the reference count goes to zero, the object is put back inside the pool
func (r ReferenceCounter) DecrementReferenceCount() {
	if atomic.LoadUint32(r.count) == 0 {
		panic("this should not happen =>" + reflect.TypeOf(r.instance).String())
	}
	if atomic.AddUint32(r.count, ^uint32(0)) == 0 {
		atomic.AddUint32(r.released, 1)
		if err := r.reset(r.instance); err != nil {
			panic("error while resetting an instance => " + err.Error())
		}
		r.destination.Put(r.instance)
		r.instance = nil
	}
}

// Method to set the current instance
func (r *ReferenceCounter) SetInstance(i interface{}) {
	r.instance = i
}

// Struct representing the pool
type referenceCountedPool struct {
	pool       *sync.Pool
	factory    func() ReferenceCountable
	returned   uint32
	allocated  uint32
	referenced uint32
}

// Method to create a new pool
func NewReferenceCountedPool(factory func(referenceCounter ReferenceCounter) ReferenceCountable, reset func(interface{}) error) *referenceCountedPool {
	p := new(referenceCountedPool)
	p.pool = new(sync.Pool)
	p.pool.New = func() interface{} {
		// Incrementing allocated count
		atomic.AddUint32(&p.allocated, 1)
		c := factory(ReferenceCounter{
			count:       new(uint32),
			destination: p.pool,
			released:    &p.returned,
			reset:       reset,
			id:          p.allocated,
		})
		return c
	}
	return p
}

// Method to get new object
func (p *referenceCountedPool) Get() ReferenceCountable {
	c := p.pool.Get().(ReferenceCountable)
	c.SetInstance(c)
	atomic.AddUint32(&p.referenced, 1)
	c.IncrementReferenceCount()
	return c
}

// Method to return reference counted pool stats
func (p *referenceCountedPool) Stats() map[string]interface{} {
	return map[string]interface{}{"allocated": p.allocated, "referenced": p.referenced, "returned": p.returned}
}
