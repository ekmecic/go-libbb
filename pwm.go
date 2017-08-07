package libbeaglebone

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// ExportState is state in which the pin is in, either exported or unexported.
type PWMState int

const (
	// Exported means the pin will be made ready to use.
	Enabled PWMState = iota
	// UnExported means the pin will be unavailable for use.
	// This is the default state when the pin is first created.
	Disabled
)

type PWM struct {
	pwmChipNum uint8
	pwmNum     uint8
	period     uint32
	dutyCycle  uint32
	state      PWMState
}

func NewPWM(pwmChipNum uint8, pwmNum uint8) *PWM {
	pwm := new(PWM)
	pwm.pwmChipNum = pwmChipNum
	pwm.pwmNum = pwmNum
	pwm.dutyCycle = 0
	pwm.state = Disabled
	return pwm
}

func (pwm *PWM) SetExportState(es ExportState) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d", pwm.pwmChipNum, pwm.pwmNum)
	if _, err := os.Stat(path); os.IsNotExist(err) && es == Exported {
		// Try to export if the GPIO isn't already exported
		file, err := os.OpenFile(fmt.Sprintf("/sys/class/pwm/pwmchip%d/export", pwm.pwmChipNum), os.O_WRONLY|os.O_SYNC, 0666)
		_, err = file.Write([]byte(strconv.Itoa(int(pwm.pwmNum))))
		if err != nil {
			return err
		}
	} else if _, err := os.Stat(path); err == nil && es == UnExported {
		// Try to unexport if the GPIO is already exported
		file, err := os.OpenFile(fmt.Sprintf("/sys/class/pwm/pwmchip%d/unexport", pwm.pwmChipNum), os.O_WRONLY|os.O_SYNC, 0666)
		_, err = file.Write([]byte(strconv.Itoa(int(pwm.pwmNum))))
		if err != nil {
			return err
		}
	} else {
		// User either tried to export an exported pin or unexport an unexported pin.
		return errors.New("Unable to export or unexport GPIO pin")
	}
	return nil
}

func (pwm *PWM) SetPeriod(period uint32) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/period", pwm.pwmChipNum, pwm.pwmNum)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)
	_, err = file.Write([]byte(strconv.Itoa(int(period))))
	if err != nil {
		return err
	}
	pwm.period = period
	return nil
}

func (pwm *PWM) SetState(ps PWMState) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/enable", pwm.pwmChipNum, pwm.pwmNum)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)
	var ps_string string
	if ps == Enabled {
		ps_string = "1"
	} else {
		ps_string = "0"
	}
	_, err = file.Write([]byte(ps_string))
	if err != nil {
		return err
	}
	pwm.state = ps
	return nil
}

func (pwm *PWM) Write(percentage float32) error {
	newDutyCycle := uint32((percentage / 100.0) * float32(pwm.period))
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/duty_cycle", pwm.pwmChipNum, pwm.pwmNum)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)
	_, err = file.Write([]byte(strconv.Itoa(int(newDutyCycle))))
	if err != nil {
		return err
	}
	pwm.dutyCycle = newDutyCycle
	return nil
}

func (pwm *PWM) SetDutyCycle(dutyCycle uint32) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/duty_cycle", pwm.pwmChipNum, pwm.pwmNum)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)
	_, err = file.Write([]byte(strconv.Itoa(int(dutyCycle))))
	if err != nil {
		return err
	}
	pwm.dutyCycle = dutyCycle
	return nil
}
