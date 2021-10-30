package cache

//TODO use redis/memory
var propertyMap map[string]string

type PropertyCAO struct {
}

func NewPropertyCAO() *PropertyCAO {
	return &PropertyCAO{}
}

func (this *PropertyCAO) Find(_key string) (string, bool, error) {
	value, ok := propertyMap[_key]
	return value, ok, nil
}

func (this *PropertyCAO) Save(_key string, _value string) error {
	propertyMap[_key] = _value
	return nil
}

func (this *PropertyCAO) Delete(_key string) error {
	delete(propertyMap, _key)
	return nil
}
