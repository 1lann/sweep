package sweep

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
)

// Scan represents a full scan.
type Scan []*ScanSample

// ScanSample represents a scan result.
type ScanSample struct {
	Angle          float64 `json:"a"` // Angle in degrees
	Distance       int     `json:"d"` // Distance in cm
	SignalStrength byte    `json:"s"` // Signal strength in god knows what
}

// StartScan starts the scan, and returns a channel that is closed when
// StopScan is called.
func (d *Device) StartScan() (<-chan Scan, error) {
	var result ResponseHeader
	err := d.ExecuteCommand(CmdDataStart, &result)
	if err != nil {
		return nil, err
	}

	switch result.CmdStatus.Int() {
	case 12:
		return nil, ErrMotorChanging
	case 13:
		return nil, ErrMotorStationary
	}

	results := make(chan Scan, 100)
	go func() {
		buffer := make(Scan, 0, 500)

		for {
			buf := make([]byte, 2)
			_, err := io.ReadFull(d.reader, buf)
			if err != nil {
				log.Println("sweep: error during read, likely fatal:", err)
				continue
			}

			if string(buf) == CmdDataStop {
				d.reader.ReadBytes('\n')
				close(results)
				return
			}

			var scanRes ResponseScanPacket
			err = rawReadDecode(io.MultiReader(bytes.NewReader(buf),
				d.reader), &scanRes)
			if err != nil {
				log.Println("sweep: error during scan:", err)
				continue
			}

			if scanRes.SyncFlags&FlagSync != 0 && len(buffer) > 0 {
				results <- buffer
				buffer = make(Scan, 0, 500)
			}

			if scanRes.Distance <= 2 {
				continue
			}

			buffer = append(buffer, &ScanSample{
				Angle:          scanRes.AngleDeg(),
				Distance:       int(scanRes.Distance),
				SignalStrength: scanRes.SignalStrength,
			})
		}
	}()

	return results, nil
}

// StopScan stops an ongoing scan.
func (d *Device) StopScan() error {
	return d.WriteCommand(CmdDataStop)
}

// Drain drains the serial buffer.
func (d *Device) Drain() {
	io.Copy(ioutil.Discard, d.reader)
}
