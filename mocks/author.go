// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/zrma/uds-go/api"
	"golang.org/x/oauth2"
)

type Author struct {
	ConfigFromJSONStub        func([]byte, ...string) (*oauth2.Config, error)
	configFromJSONMutex       sync.RWMutex
	configFromJSONArgsForCall []struct {
		arg1 []byte
		arg2 []string
	}
	configFromJSONReturns struct {
		result1 *oauth2.Config
		result2 error
	}
	configFromJSONReturnsOnCall map[int]struct {
		result1 *oauth2.Config
		result2 error
	}
	GetTokenStub        func(*oauth2.Config, string, func() (string, error)) (*oauth2.Token, error)
	getTokenMutex       sync.RWMutex
	getTokenArgsForCall []struct {
		arg1 *oauth2.Config
		arg2 string
		arg3 func() (string, error)
	}
	getTokenReturns struct {
		result1 *oauth2.Token
		result2 error
	}
	getTokenReturnsOnCall map[int]struct {
		result1 *oauth2.Token
		result2 error
	}
	ReadFileStub        func(string) ([]byte, error)
	readFileMutex       sync.RWMutex
	readFileArgsForCall []struct {
		arg1 string
	}
	readFileReturns struct {
		result1 []byte
		result2 error
	}
	readFileReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Author) ConfigFromJSON(arg1 []byte, arg2 ...string) (*oauth2.Config, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.configFromJSONMutex.Lock()
	ret, specificReturn := fake.configFromJSONReturnsOnCall[len(fake.configFromJSONArgsForCall)]
	fake.configFromJSONArgsForCall = append(fake.configFromJSONArgsForCall, struct {
		arg1 []byte
		arg2 []string
	}{arg1Copy, arg2})
	fake.recordInvocation("ConfigFromJSON", []interface{}{arg1Copy, arg2})
	fake.configFromJSONMutex.Unlock()
	if fake.ConfigFromJSONStub != nil {
		return fake.ConfigFromJSONStub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.configFromJSONReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Author) ConfigFromJSONCallCount() int {
	fake.configFromJSONMutex.RLock()
	defer fake.configFromJSONMutex.RUnlock()
	return len(fake.configFromJSONArgsForCall)
}

func (fake *Author) ConfigFromJSONCalls(stub func([]byte, ...string) (*oauth2.Config, error)) {
	fake.configFromJSONMutex.Lock()
	defer fake.configFromJSONMutex.Unlock()
	fake.ConfigFromJSONStub = stub
}

func (fake *Author) ConfigFromJSONArgsForCall(i int) ([]byte, []string) {
	fake.configFromJSONMutex.RLock()
	defer fake.configFromJSONMutex.RUnlock()
	argsForCall := fake.configFromJSONArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Author) ConfigFromJSONReturns(result1 *oauth2.Config, result2 error) {
	fake.configFromJSONMutex.Lock()
	defer fake.configFromJSONMutex.Unlock()
	fake.ConfigFromJSONStub = nil
	fake.configFromJSONReturns = struct {
		result1 *oauth2.Config
		result2 error
	}{result1, result2}
}

func (fake *Author) ConfigFromJSONReturnsOnCall(i int, result1 *oauth2.Config, result2 error) {
	fake.configFromJSONMutex.Lock()
	defer fake.configFromJSONMutex.Unlock()
	fake.ConfigFromJSONStub = nil
	if fake.configFromJSONReturnsOnCall == nil {
		fake.configFromJSONReturnsOnCall = make(map[int]struct {
			result1 *oauth2.Config
			result2 error
		})
	}
	fake.configFromJSONReturnsOnCall[i] = struct {
		result1 *oauth2.Config
		result2 error
	}{result1, result2}
}

func (fake *Author) GetToken(arg1 *oauth2.Config, arg2 string, arg3 func() (string, error)) (*oauth2.Token, error) {
	fake.getTokenMutex.Lock()
	ret, specificReturn := fake.getTokenReturnsOnCall[len(fake.getTokenArgsForCall)]
	fake.getTokenArgsForCall = append(fake.getTokenArgsForCall, struct {
		arg1 *oauth2.Config
		arg2 string
		arg3 func() (string, error)
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetToken", []interface{}{arg1, arg2, arg3})
	fake.getTokenMutex.Unlock()
	if fake.GetTokenStub != nil {
		return fake.GetTokenStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getTokenReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Author) GetTokenCallCount() int {
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
	return len(fake.getTokenArgsForCall)
}

func (fake *Author) GetTokenCalls(stub func(*oauth2.Config, string, func() (string, error)) (*oauth2.Token, error)) {
	fake.getTokenMutex.Lock()
	defer fake.getTokenMutex.Unlock()
	fake.GetTokenStub = stub
}

func (fake *Author) GetTokenArgsForCall(i int) (*oauth2.Config, string, func() (string, error)) {
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
	argsForCall := fake.getTokenArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *Author) GetTokenReturns(result1 *oauth2.Token, result2 error) {
	fake.getTokenMutex.Lock()
	defer fake.getTokenMutex.Unlock()
	fake.GetTokenStub = nil
	fake.getTokenReturns = struct {
		result1 *oauth2.Token
		result2 error
	}{result1, result2}
}

func (fake *Author) GetTokenReturnsOnCall(i int, result1 *oauth2.Token, result2 error) {
	fake.getTokenMutex.Lock()
	defer fake.getTokenMutex.Unlock()
	fake.GetTokenStub = nil
	if fake.getTokenReturnsOnCall == nil {
		fake.getTokenReturnsOnCall = make(map[int]struct {
			result1 *oauth2.Token
			result2 error
		})
	}
	fake.getTokenReturnsOnCall[i] = struct {
		result1 *oauth2.Token
		result2 error
	}{result1, result2}
}

func (fake *Author) ReadFile(arg1 string) ([]byte, error) {
	fake.readFileMutex.Lock()
	ret, specificReturn := fake.readFileReturnsOnCall[len(fake.readFileArgsForCall)]
	fake.readFileArgsForCall = append(fake.readFileArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ReadFile", []interface{}{arg1})
	fake.readFileMutex.Unlock()
	if fake.ReadFileStub != nil {
		return fake.ReadFileStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.readFileReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Author) ReadFileCallCount() int {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	return len(fake.readFileArgsForCall)
}

func (fake *Author) ReadFileCalls(stub func(string) ([]byte, error)) {
	fake.readFileMutex.Lock()
	defer fake.readFileMutex.Unlock()
	fake.ReadFileStub = stub
}

func (fake *Author) ReadFileArgsForCall(i int) string {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	argsForCall := fake.readFileArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Author) ReadFileReturns(result1 []byte, result2 error) {
	fake.readFileMutex.Lock()
	defer fake.readFileMutex.Unlock()
	fake.ReadFileStub = nil
	fake.readFileReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Author) ReadFileReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.readFileMutex.Lock()
	defer fake.readFileMutex.Unlock()
	fake.ReadFileStub = nil
	if fake.readFileReturnsOnCall == nil {
		fake.readFileReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.readFileReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Author) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.configFromJSONMutex.RLock()
	defer fake.configFromJSONMutex.RUnlock()
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Author) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ api.Author = new(Author)
