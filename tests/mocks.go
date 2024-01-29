package protoconf

type MockProvider struct {
	ReadBytesCalled bool
	ReadCalled      bool
}

func (p *MockProvider) ReadBytes() ([]byte, error) {
	p.ReadBytesCalled = true
	return nil, nil
}

func (p *MockProvider) Read() (map[string]interface{}, error) {
	p.ReadCalled = true
	return make(map[string]interface{}), nil
}

type MockParser struct {
	UnmarshalCalled bool
}

func (p *MockParser) Unmarshal(_ []byte) (map[string]interface{}, error) {
	p.UnmarshalCalled = true
	return make(map[string]interface{}), nil
}
