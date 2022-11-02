package hal

type ChipMode int
type Parity byte
type OnMessageCb func([]byte, error)

const (
	ParityNone  Parity = 'N'
	ParityOdd   Parity = 'O'
	ParityEven  Parity = 'E'
	ParityMark  Parity = 'M' // parity bit is always 1
	ParitySpace Parity = 'S' // parity bit is always 0
)

const (
	ModeNormal ChipMode = iota
	ModeWakeUp
	ModePowerSave
	ModeSleep
)

type HWHandler interface {
	ReadSerial() ([]byte, error)
	WriteSerial(msg []byte) error
	StageSerialPortConfig(baudRate int, parityBit Parity)
	SetMode(mode ChipMode) error
	GetMode() (ChipMode, error)
	RegisterOnMessageCb(OnMessageCb) error
}
