package compiler

import "sync"

// buildFixturePipeline evaluates the expensive M7 quadrature once, then feeds
// each immutable milestone value into the next stage. M7 is the one stage that
// extends its input graph, so it receives a graph clone and cannot mutate M6.
var (
	baseOnce sync.Once
	baseErr  error
	baseM6   M6Result
	baseM7   M7Result
	baseM8   M8Result
	baseM9   M9Result
	baseM10  M10Result
	baseM11  M11Result
	baseM12  M12Result
	baseM13  M13Result
	baseM14A M14AResult
	baseM15  M15Result

	m16Once      sync.Once
	m16Cached    M16Result
	m16CachedErr error
	m17Once      sync.Once
	m17Cached    M17Result
	m17CachedErr error
)

func buildFixturePipeline() {
	baseOnce.Do(func() {
		if baseM6, baseErr = CompileM6(); baseErr != nil {
			return
		}
		m6ForM7 := baseM6
		m6ForM7.M5.Graph = baseM6.M5.Graph.Clone()
		if baseM7, baseErr = compileM7FromM6(m6ForM7, M7Options{}); baseErr != nil {
			return
		}
		if baseM8, baseErr = compileM8FromM7(baseM7); baseErr != nil {
			return
		}
		if baseM9, baseErr = compileM9FromM8(baseM8); baseErr != nil {
			return
		}
		if baseM10, baseErr = compileM10FromM9(baseM9); baseErr != nil {
			return
		}
		if baseM11, baseErr = compileM11FromM10(baseM10); baseErr != nil {
			return
		}
		if baseM12, baseErr = compileM12FromM11(baseM11); baseErr != nil {
			return
		}
		if baseM13, baseErr = compileM13FromM12(baseM12); baseErr != nil {
			return
		}
		if baseM14A, baseErr = compileM14AFromM13(baseM13); baseErr != nil {
			return
		}
		baseM15, baseErr = compileM15FromM14A(baseM14A)
	})
}

func testM6() (M6Result, error)     { buildFixturePipeline(); return baseM6, baseErr }
func testM7() (M7Result, error)     { buildFixturePipeline(); return baseM7, baseErr }
func testM8() (M8Result, error)     { buildFixturePipeline(); return baseM8, baseErr }
func testM9() (M9Result, error)     { buildFixturePipeline(); return baseM9, baseErr }
func testM10() (M10Result, error)   { buildFixturePipeline(); return baseM10, baseErr }
func testM11() (M11Result, error)   { buildFixturePipeline(); return baseM11, baseErr }
func testM12() (M12Result, error)   { buildFixturePipeline(); return baseM12, baseErr }
func testM13() (M13Result, error)   { buildFixturePipeline(); return baseM13, baseErr }
func testM14A() (M14AResult, error) { buildFixturePipeline(); return baseM14A, baseErr }
func testM15() (M15Result, error)   { buildFixturePipeline(); return baseM15, baseErr }

func testM16() (M16Result, error) {
	m16Once.Do(func() {
		buildFixturePipeline()
		if baseErr != nil {
			m16CachedErr = baseErr
			return
		}
		m16Cached, m16CachedErr = compileM16FromM15(baseM15, false)
	})
	return m16Cached, m16CachedErr
}

func testM17() (M17Result, error) {
	m17Once.Do(func() {
		m16, err := testM16()
		if err != nil {
			m17CachedErr = err
			return
		}
		m17Cached, m17CachedErr = compileM17FromM16(m16, false)
	})
	return m17Cached, m17CachedErr
}
