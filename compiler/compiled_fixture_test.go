package compiler

import "sync"

var (
	m6Once       sync.Once
	m6Cached     M6Result
	m6CachedErr  error
	m8Once       sync.Once
	m8Cached     M8Result
	m8CachedErr  error
	m9Once       sync.Once
	m9Cached     M9Result
	m9CachedErr  error
	m10Once      sync.Once
	m10Cached    M10Result
	m10CachedErr error
	m11Once      sync.Once
	m11Cached    M11Result
	m11CachedErr error
)

func testM6() (M6Result, error) {
	m6Once.Do(func() { m6Cached, m6CachedErr = CompileM6() })
	return m6Cached, m6CachedErr
}

func testM8() (M8Result, error) {
	m8Once.Do(func() { m8Cached, m8CachedErr = CompileM8() })
	return m8Cached, m8CachedErr
}

func testM9() (M9Result, error) {
	m9Once.Do(func() { m9Cached, m9CachedErr = CompileM9() })
	return m9Cached, m9CachedErr
}

func testM10() (M10Result, error) {
	m10Once.Do(func() { m10Cached, m10CachedErr = CompileM10() })
	return m10Cached, m10CachedErr
}

func testM11() (M11Result, error) {
	m11Once.Do(func() { m11Cached, m11CachedErr = CompileM11() })
	return m11Cached, m11CachedErr
}
