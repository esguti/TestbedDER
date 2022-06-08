module github.com/TRICERA-energy/sunspec

go 1.16

require (
	github.com/GoAethereal/cancel v0.0.1
	github.com/GoAethereal/modbus v0.0.3
)


replace (
	github.com/GoAethereal/cancel => ./modules/cancel
	github.com/GoAethereal/modbus => ./modules/modbus
)