package hashmap

type Dir struct {
	FileId    string
	Timestamp uint64
	ValueSize uint16
	ValuePos  uint64
}

var index map[string]Dir

func Init() {
	index = make(map[string]Dir)
}

func Upsert(key string, val Dir) {
	index[key] = val
}

func Get(key string) Dir {
	return index[key]
}

func Delete(key string) {
	delete(index, key)
}
