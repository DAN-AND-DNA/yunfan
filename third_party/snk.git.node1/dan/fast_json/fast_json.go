package fast_json

import (
	"github.com/segmentio/encoding/json"
	"sync"
	"unsafe"
)

var (
	is_inited   = false
	buffer_pool sync.Pool
)

type Json_pool struct {

	// hide
	buffer_pool *sync.Pool
}

type Json_buffer struct {
	size int
	data []byte
}

func (this *Json_buffer) clean() {
	if len(this.data) > this.size {
		this.data = make([]byte, this.size)
	}

	this.data = this.data[:0]
}

func (this *Json_buffer) Data() []byte {
	return this.data
}

func (this *Json_buffer) String() string {
	return *(*string)(unsafe.Pointer(&this.data))
}

func (this *Json_buffer) String_simple() string {
	return string(this.data)
}

func (this *Json_buffer) Len() int {
	return len(this.data)
}

func New(raw_pre_alloced ...int) *Json_pool {
	jp := &Json_pool{}
	pre_alloced := 500
	if len(raw_pre_alloced) == 0 || raw_pre_alloced[0] <= 0 {
	} else {
		pre_alloced = raw_pre_alloced[0]
	}

	jp.buffer_pool = &sync.Pool{
		New: func() interface{} {
			buffer := &Json_buffer{
				size: pre_alloced,
				data: make([]byte, pre_alloced),
			}
			return buffer
		},
	}

	return jp
}

func (this *Json_pool) Alloc_buffer() *Json_buffer {
	return this.buffer_pool.Get().(*Json_buffer)
}

func (this *Json_pool) Release_buffer(ptr *Json_buffer) {
	ptr.clean()
	this.buffer_pool.Put(ptr)
}

func Marshal(ptr_buffer *Json_buffer, src interface{}) error {
	ptr_buffer.data = ptr_buffer.data[:0]
	byte_buffer, err := json.Append(ptr_buffer.data, src, json.TrustRawMessage)
	if err != nil {
		return err
	}
	ptr_buffer.data = byte_buffer
	return nil
}

func Marshal_simple(src interface{}) ([]byte, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Unmarshal_simple(byte_json []byte, to interface{}) error {
	err := json.Unmarshal(byte_json, to)
	if err != nil {
		return err
	}

	return nil
}

func Unmarshal(byte_json []byte, to interface{}) error {
	b, err := json.Parse(byte_json, to, json.ZeroCopy)
	_ = b
	if err != nil {
		return err
	}

	return nil
}
