// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/zrma/uds-go/api"
	"golang.org/x/oauth2"
)

type Author struct {
	GetTokenStub        func(*oauth2.Config) *oauth2.Token
	getTokenMutex       sync.RWMutex
	getTokenArgsForCall []struct {
		arg1 *oauth2.Config
	}
	getTokenReturns struct {
		result1 *oauth2.Token
	}
	getTokenReturnsOnCall map[int]struct {
		result1 *oauth2.Token
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Author) GetToken(arg1 *oauth2.Config) *oauth2.Token {
	fake.getTokenMutex.Lock()
	ret, specificReturn := fake.getTokenReturnsOnCall[len(fake.getTokenArgsForCall)]
	fake.getTokenArgsForCall = append(fake.getTokenArgsForCall, struct {
		arg1 *oauth2.Config
	}{arg1})
	fake.recordInvocation("GetToken", []interface{}{arg1})
	fake.getTokenMutex.Unlock()
	if fake.GetTokenStub != nil {
		return fake.GetTokenStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getTokenReturns
	return fakeReturns.result1
}

func (fake *Author) GetTokenCallCount() int {
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
	return len(fake.getTokenArgsForCall)
}

func (fake *Author) GetTokenCalls(stub func(*oauth2.Config) *oauth2.Token) {
	fake.getTokenMutex.Lock()
	defer fake.getTokenMutex.Unlock()
	fake.GetTokenStub = stub
}

func (fake *Author) GetTokenArgsForCall(i int) *oauth2.Config {
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
	argsForCall := fake.getTokenArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Author) GetTokenReturns(result1 *oauth2.Token) {
	fake.getTokenMutex.Lock()
	defer fake.getTokenMutex.Unlock()
	fake.GetTokenStub = nil
	fake.getTokenReturns = struct {
		result1 *oauth2.Token
	}{result1}
}

func (fake *Author) GetTokenReturnsOnCall(i int, result1 *oauth2.Token) {
	fake.getTokenMutex.Lock()
	defer fake.getTokenMutex.Unlock()
	fake.GetTokenStub = nil
	if fake.getTokenReturnsOnCall == nil {
		fake.getTokenReturnsOnCall = make(map[int]struct {
			result1 *oauth2.Token
		})
	}
	fake.getTokenReturnsOnCall[i] = struct {
		result1 *oauth2.Token
	}{result1}
}

func (fake *Author) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
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
