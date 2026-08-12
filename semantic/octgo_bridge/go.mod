module github.com/yuechen-li-dev/Riemann/semantic/octgo_bridge

go 1.25.0

require (
	github.com/yuechen-li-dev/Riemann v0.0.0
	github.com/yuechen-li-dev/oct v1.0.1-0.20260812180913-0fdca7a3365e
)

replace github.com/yuechen-li-dev/Riemann => ../..

replace github.com/yuechen-li-dev/oct => ../../../oct
