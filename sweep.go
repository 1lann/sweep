package sweep

import (
	"bufio"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/tarm/serial"
)

// Possible errors returned.
var (
	ErrTimeout          = errors.New("sweep: timeout")
	ErrCommandMismatch  = errors.New("sweep: read response command mismatch")
	ErrInvalidParameter = errors.New("sweep: invalid parameter")
	ErrMotorChanging    = errors.New("sweep: motor changing")
	ErrInvalidResponse  = errors.New("sweep: invalid response")
	ErrMotorStationary  = errors.New("sweep: motor is stationary")
)

// Valid commands.
var (
	CmdMotorReady       = "MZ"
	CmdMotorSpeedAdjust = "MS"
	CmdDataStart        = "DS"
	CmdDataStop         = "DX"
	CmdMotorInfo        = "MI"
	CmdSampleRateAdjust = "LR"
	CmdSampleRateInfo   = "LI"
	CmdVersionInfo      = "IV"
	CmdDeviceInfo       = "ID"
	CmdResetDevice      = "RR"
)

// Valid sample rates.
const (
	Rate500  = 500
	Rate750  = 750
	Rate1000 = 1000
)

// Device represents a Scanse Sweep device.
type Device struct {
	reader *bufio.Reader
	serial *serial.Port
	path   string
}

// NewDevice returns a new scanse sweep controller.
func NewDevice(path string) (*Device, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name:        path,
		Baud:        115200,
		ReadTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Device{
		reader: bufio.NewReader(port),
		serial: port,
		path:   path,
	}, nil
}

// ExecuteCommand executes a command and stores the result in result.
func (d *Device) ExecuteCommand(cmd string, result interface{},
	args ...string) error {
	if err := d.WriteCommand(cmd, args...); err != nil {
		return err
	}

	if err := d.ReadDecode(result); err != nil {
		return err
	}

	value := reflect.ValueOf(result).Elem()
	cmdField, found := value.Type().FieldByName("Cmd")
	if found && cmdField.Index[0] == 0 {
		cmdResp, ok := value.Field(0).Interface().([2]byte)
		if ok && !cmpCmd(cmdResp, cmd) {
			log.Printf("got: %s, expected: %s\n", string(cmdResp[:]), cmd)
			return ErrCommandMismatch
		}
	}

	return nil
}

// WaitUntilMotorReady waits until the motor is ready.
func (d *Device) WaitUntilMotorReady() error {
	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()
	timeout := time.After(10 * time.Second)

	for {
		select {
		case <-t.C:
			ready, err := d.GetMotorReady()
			if err != nil {
				return err
			}

			if !ready {
				break
			}

			return nil
		case <-timeout:
			return ErrTimeout
		}
	}
}

// WriteCommand writes a command to the Sweep.
func (d *Device) WriteCommand(cmd string, args ...string) error {
	if len(args) > 0 {
		_, err := d.serial.Write([]byte(cmd + args[0] + "\n"))
		return err
	}

	_, err := d.serial.Write([]byte(cmd + "\n"))
	return err
}

func cmpCmd(cmd [2]byte, cmdStr string) bool {
	return cmd[0] == cmdStr[0] && cmd[1] == cmdStr[1]
}

// GetMotorReady returns whether or not the motor is ready.
func (d *Device) GetMotorReady() (bool, error) {
	var result ResponseMotorReady
	if err := d.ExecuteCommand(CmdMotorReady, &result); err != nil {
		return false, err
	}

	return result.MotorReady.Int() == 0, nil
}

// GetMotorSpeed returns the speed of the motor.
func (d *Device) GetMotorSpeed() (int, error) {
	var result ResponseMotorInfo
	err := d.ExecuteCommand(CmdMotorInfo, &result)

	return result.MotorSpeed.Int(), err
}

// SetMotorSpeed sets the speed of the motor in Hz.
func (d *Device) SetMotorSpeed(speed int) error {
	if speed < 0 || speed > 10 {
		return ErrInvalidParameter
	}

	var resultA ResponseParamA
	err := d.ExecuteCommand(CmdMotorSpeedAdjust, &resultA, NewInt2(speed).String())
	if err != nil {
		return err
	}

	var resultB ResponseParamB
	if err := d.ReadDecode(&resultB); err != nil {
		return err
	}

	switch resultB.CmdStatus.Int() {
	case 11:
		return ErrInvalidParameter
	case 12:
		return ErrMotorChanging
	}

	return nil
}

// GetSampleRate returns the sample rate of the sweep.
func (d *Device) GetSampleRate() (int, error) {
	var result ResponseSampleRate
	if err := d.ExecuteCommand(CmdSampleRateInfo, &result); err != nil {
		return 0, err
	}

	switch result.SampleRate.Int() {
	case 1:
		return Rate500, nil
	case 2:
		return Rate750, nil
	case 3:
		return Rate1000, nil
	}

	return 0, ErrInvalidResponse
}

// SetSampleRate sets the sample rate of the sweep.
func (d *Device) SetSampleRate(rate int) error {
	switch rate {
	case Rate500, Rate750, Rate1000:
	default:
		return ErrInvalidParameter
	}

	var resultA ResponseParamA
	if err := d.ExecuteCommand(CmdSampleRateAdjust, &resultA); err != nil {
		return err
	}

	var resultB ResponseParamB
	if err := d.ReadDecode(&resultB); err != nil {
		return err
	}

	switch resultB.CmdStatus.Int() {
	case 11:
		return ErrInvalidParameter
	default:
		return nil
	}
}

// Restart restarts the device by resetting it, then re-connecting to it.
func (d *Device) Restart() error {
	d.Reset()

	for i := 0; i < 2; i++ {
		_, err := d.GetMotorReady()
		if err != nil {
			continue
		}

		return nil
	}

	return ErrTimeout
}

// Reset sends the reset command to the sweep.
func (d *Device) Reset() error {
	return d.WriteCommand(CmdResetDevice)
}

// Close closes the underlying serial connection.
func (d *Device) Close() error {
	return d.serial.Close()
}
