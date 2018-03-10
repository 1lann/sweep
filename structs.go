package sweep

// ResponseHeader represents a command response header.
type ResponseHeader struct {
	Cmd       [2]byte
	CmdStatus Int2
	CmdSum    byte
}

// CommandParamPacket ...
type CommandParamPacket struct {
	Cmd      [2]byte
	CmdParam [2]byte
}

// ResponseParamA represents the first part of a response parameter header.
type ResponseParamA struct {
	Cmd      [2]byte
	CmdParam [2]byte
}

// ResponseParamB represents the second part of a response parameter header.
type ResponseParamB struct {
	CmdStatus Int2
	CmdSum    byte
}

// ResponseDevice ...
type ResponseDevice struct {
	Cmd        [2]byte
	BitRate    [6]byte
	LaserState byte
	Mode       byte
	Diagnostic byte
	MotorSpeed Int2
	SampleRate Int4
}

// ResponseVersion ...
type ResponseVersion struct {
	Cmd             [2]byte
	Model           [5]byte
	ProtocolMajor   byte
	ProtocolMinor   byte
	FirmwareMajor   byte
	FirmwareMinor   byte
	HardwareVersion byte
	SerialNum       [8]byte
}

// ResponseMotorReady ...
type ResponseMotorReady struct {
	Cmd        [2]byte
	MotorReady Int2
}

// ResponseMotorInfo ...
type ResponseMotorInfo struct {
	Cmd        [2]byte
	MotorSpeed Int2
}

// ResponseSampleRate ...
type ResponseSampleRate struct {
	Cmd        [2]byte
	SampleRate Int2
}

// Valid SyncFlags for ResponseScanPacket.
const (
	FlagSync = 1 << iota
	FlagCommunicationFail
)

// ResponseScanPacket represents a response scan packet.
type ResponseScanPacket struct {
	SyncFlags      byte
	Angle          uint16
	Distance       uint16
	SignalStrength byte
	Checksum       byte
}

// Checksum returns the checksum of the ResponseHeader
func (h *ResponseHeader) Checksum() byte {
	return ((h.Cmd[0] + h.Cmd[1]) & 0x3F) + 0x30
}

// Checksum returns the checksum of the ResponseParam.
func (p *ResponseParamA) Checksum() byte {
	return ((p.Cmd[0] + p.Cmd[1]) & 0x3F) + 0x30
}

// AngleDeg returns the angle of the scan in degrees.
func (p *ResponseScanPacket) AngleDeg() float64 {
	return float64(p.Angle) / 16
}
