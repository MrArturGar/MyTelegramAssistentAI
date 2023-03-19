package modules

import (
	"MyTelegramAssistentAI/src/modules/dalle"
	"errors"
	"sync"
)

type Module interface {
	Execute(*string) *string
}

type NamedObject struct {
	// Структура, содержащая объект и его имя.
	Name   string
	Object Module
}

type Pool struct {
	// Объектный пул.
	pool map[string]chan *NamedObject
	// Мьютекс для защиты объектного пула от одновременного доступа из разных горутин.
	mutex sync.Mutex
}

var instance *Pool
var once sync.Once

func GetInstance() *Pool {
	once.Do(func() {
		instance = createPool(10)
	})
	return instance
}

func createPool(maxSize int) *Pool {
	return &Pool{
		pool: make(map[string]chan *NamedObject),
	}
}

func (p *Pool) Get(name string) (Module, error) {
	// Получение объекта из пула.
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ch, ok := p.pool[name]
	if !ok {
		// Канала с таким именем нет, создаем новый.
		ch = make(chan *NamedObject)
		p.pool[name] = ch
	}

	select {
	case namedObject := <-ch:
		// Объект есть в пуле, просто извлекаем его.
		return namedObject.Object, nil
	default:
		// Объекта в пуле нет, создаем новый.
		var module Module

		switch name {
		case "dall-e":
			module = dalle.GetInstance()
		default:
			return nil, errors.New("The module named " + name + " does not exist.")
		}
		return module, nil
	}
}

func (p *Pool) Put(name string, object Module) error {
	// Возвращение объекта в пул по имени.
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ch, ok := p.pool[name]
	if !ok {
		return errors.New("no such pool for the name")
	}

	select {
	case ch <- &NamedObject{Name: name, Object: object}:
		// Объект успешно возвращен в пул.
	default:
		// Пул заполнен, объект не помещается в него, просто отбрасываем его.
	}

	return nil
}
