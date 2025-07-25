// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"sync"

	"github.com/bborbe/kv"
)

type DB struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	RemoveStub        func() error
	removeMutex       sync.RWMutex
	removeArgsForCall []struct {
	}
	removeReturns struct {
		result1 error
	}
	removeReturnsOnCall map[int]struct {
		result1 error
	}
	SyncStub        func() error
	syncMutex       sync.RWMutex
	syncArgsForCall []struct {
	}
	syncReturns struct {
		result1 error
	}
	syncReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateStub        func(context.Context, func(ctx context.Context, tx kv.Tx) error) error
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 context.Context
		arg2 func(ctx context.Context, tx kv.Tx) error
	}
	updateReturns struct {
		result1 error
	}
	updateReturnsOnCall map[int]struct {
		result1 error
	}
	ViewStub        func(context.Context, func(ctx context.Context, tx kv.Tx) error) error
	viewMutex       sync.RWMutex
	viewArgsForCall []struct {
		arg1 context.Context
		arg2 func(ctx context.Context, tx kv.Tx) error
	}
	viewReturns struct {
		result1 error
	}
	viewReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DB) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	stub := fake.CloseStub
	fakeReturns := fake.closeReturns
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DB) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *DB) CloseCalls(stub func() error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *DB) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *DB) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DB) Remove() error {
	fake.removeMutex.Lock()
	ret, specificReturn := fake.removeReturnsOnCall[len(fake.removeArgsForCall)]
	fake.removeArgsForCall = append(fake.removeArgsForCall, struct {
	}{})
	stub := fake.RemoveStub
	fakeReturns := fake.removeReturns
	fake.recordInvocation("Remove", []interface{}{})
	fake.removeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DB) RemoveCallCount() int {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return len(fake.removeArgsForCall)
}

func (fake *DB) RemoveCalls(stub func() error) {
	fake.removeMutex.Lock()
	defer fake.removeMutex.Unlock()
	fake.RemoveStub = stub
}

func (fake *DB) RemoveReturns(result1 error) {
	fake.removeMutex.Lock()
	defer fake.removeMutex.Unlock()
	fake.RemoveStub = nil
	fake.removeReturns = struct {
		result1 error
	}{result1}
}

func (fake *DB) RemoveReturnsOnCall(i int, result1 error) {
	fake.removeMutex.Lock()
	defer fake.removeMutex.Unlock()
	fake.RemoveStub = nil
	if fake.removeReturnsOnCall == nil {
		fake.removeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DB) Sync() error {
	fake.syncMutex.Lock()
	ret, specificReturn := fake.syncReturnsOnCall[len(fake.syncArgsForCall)]
	fake.syncArgsForCall = append(fake.syncArgsForCall, struct {
	}{})
	stub := fake.SyncStub
	fakeReturns := fake.syncReturns
	fake.recordInvocation("Sync", []interface{}{})
	fake.syncMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DB) SyncCallCount() int {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	return len(fake.syncArgsForCall)
}

func (fake *DB) SyncCalls(stub func() error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = stub
}

func (fake *DB) SyncReturns(result1 error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = nil
	fake.syncReturns = struct {
		result1 error
	}{result1}
}

func (fake *DB) SyncReturnsOnCall(i int, result1 error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = nil
	if fake.syncReturnsOnCall == nil {
		fake.syncReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.syncReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DB) Update(arg1 context.Context, arg2 func(ctx context.Context, tx kv.Tx) error) error {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 context.Context
		arg2 func(ctx context.Context, tx kv.Tx) error
	}{arg1, arg2})
	stub := fake.UpdateStub
	fakeReturns := fake.updateReturns
	fake.recordInvocation("Update", []interface{}{arg1, arg2})
	fake.updateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DB) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *DB) UpdateCalls(stub func(context.Context, func(ctx context.Context, tx kv.Tx) error) error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *DB) UpdateArgsForCall(i int) (context.Context, func(ctx context.Context, tx kv.Tx) error) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *DB) UpdateReturns(result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 error
	}{result1}
}

func (fake *DB) UpdateReturnsOnCall(i int, result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DB) View(arg1 context.Context, arg2 func(ctx context.Context, tx kv.Tx) error) error {
	fake.viewMutex.Lock()
	ret, specificReturn := fake.viewReturnsOnCall[len(fake.viewArgsForCall)]
	fake.viewArgsForCall = append(fake.viewArgsForCall, struct {
		arg1 context.Context
		arg2 func(ctx context.Context, tx kv.Tx) error
	}{arg1, arg2})
	stub := fake.ViewStub
	fakeReturns := fake.viewReturns
	fake.recordInvocation("View", []interface{}{arg1, arg2})
	fake.viewMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *DB) ViewCallCount() int {
	fake.viewMutex.RLock()
	defer fake.viewMutex.RUnlock()
	return len(fake.viewArgsForCall)
}

func (fake *DB) ViewCalls(stub func(context.Context, func(ctx context.Context, tx kv.Tx) error) error) {
	fake.viewMutex.Lock()
	defer fake.viewMutex.Unlock()
	fake.ViewStub = stub
}

func (fake *DB) ViewArgsForCall(i int) (context.Context, func(ctx context.Context, tx kv.Tx) error) {
	fake.viewMutex.RLock()
	defer fake.viewMutex.RUnlock()
	argsForCall := fake.viewArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *DB) ViewReturns(result1 error) {
	fake.viewMutex.Lock()
	defer fake.viewMutex.Unlock()
	fake.ViewStub = nil
	fake.viewReturns = struct {
		result1 error
	}{result1}
}

func (fake *DB) ViewReturnsOnCall(i int, result1 error) {
	fake.viewMutex.Lock()
	defer fake.viewMutex.Unlock()
	fake.ViewStub = nil
	if fake.viewReturnsOnCall == nil {
		fake.viewReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.viewReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DB) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DB) recordInvocation(key string, args []interface{}) {
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

var _ kv.DB = new(DB)
