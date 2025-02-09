package active

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type HasName interface {
	Name() string
}

type ActiveList[T HasName] struct {
	activeList *list.List
	mutex      *sync.RWMutex
}

func (a *ActiveList[T]) Insert(t T) *list.Element {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	fmt.Printf("Adding %T: %q\n", t, t.Name())
	el := a.activeList.PushBack(t)
	return el
}

func (a *ActiveList[T]) status() {
	if a.activeList.Len() > 0 {
		t := a.activeList.Front().Value.(T)
		fmt.Printf("Active %T [%d]:\n", t, a.activeList.Len())
		for e := a.activeList.Front(); e != nil; e = e.Next() {
			fmt.Printf("%v\n", e.Value.(T))
		}
	}
}

func (a *ActiveList[T]) Print() {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	a.status()
}

func (a *ActiveList[T]) Remove(e *list.Element) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.activeList.Len() == 0 {
		return
	}
	t := a.activeList.Remove(e)
	fmt.Printf("Removing %T: %q\n", t, t.(T).Name())
}

func (a *ActiveList[T]) Len() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.activeList.Len()
}

func (a *ActiveList[T]) DropFront() (T, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.activeList.Len() == 0 {
		return *new(T), errors.New("active list is empty")
	}
	rem := a.activeList.Remove(a.activeList.Front())
	return rem.(T), nil
}

func NewActiveList[T HasName]() *ActiveList[T] {
	return &ActiveList[T]{
		activeList: list.New(),
		mutex:      &sync.RWMutex{},
	}
}
