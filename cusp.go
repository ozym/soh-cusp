package main

import (
	"encoding/xml"
	"os"
	"strconv"
	"time"
)

type Cusp struct {
	XMLName           xml.Name `xml:"report"`
	Title             string   `xml:"title"`
	Instrument        string   `xml:"instrument"`
	SystemDateAndTime string   `xml:"system_date_and_time"`
	ClockError        string   `xml:"clock_error"`
	DASTimeLock       string   `xml:"DAS_time_lock"`
	GPSState          string   `xml:"GPS_state"`
	GPSLossPeriod     string   `xml:"GPS_loss_period"`
	TimeSrcA          string   `xml:"Time_src_A"`
	TimeOffsetA       string   `xml:"Time_offset_A"`
	TimeJitterA       string   `xml:"Time_jitter_A"`
	TimeSuccessA      string   `xml:"Time_success_A"`
	TimeSrcB          string   `xml:"Time_src_B"`
	TimeSrcC          string   `xml:"Time_src_C"`
	InstSwVer         string   `xml:"inst_sw_ver"`
	SampBrdSer        string   `xml:"samp_brd_ser"`
	SampBrdSwVer      string   `xml:"samp_brd_sw_ver"`
	SysBattVolt       string   `xml:"sys_batt_volt"`
	SysBattVoltMax    string   `xml:"sys_batt_volt_max"`
	SysBattVoltMin    string   `xml:"sys_batt_volt_min"`
	SysTemp           string   `xml:"sys_temp"`
	SysTempMax        string   `xml:"sys_temp_max"`
	SysTempMin        string   `xml:"sys_temp_min"`
	SampRate          string   `xml:"samp_rate"`
	PreEvtLen         string   `xml:"pre_evt_len"`
	PostEvtLen        string   `xml:"post_evt_len"`
	CurrXNoise        string   `xml:"curr_X_noise"`
	CurrYNoise        string   `xml:"curr_Y_noise"`
	CurrZNoise        string   `xml:"curr_Z_noise"`
	CurrXStalta       string   `xml:"curr_X_stalta"`
	CurrYStalta       string   `xml:"curr_Y_stalta"`
	CurrZStalta       string   `xml:"curr_Z_stalta"`
	CurrXOffset       string   `xml:"curr_X_offset"`
	CurrYOffset       string   `xml:"curr_Y_offset"`
	CurrZOffset       string   `xml:"curr_Z_offset"`
	CurrXAdj          string   `xml:"curr_X_adj"`
	CurrYAdj          string   `xml:"curr_Y_adj"`
	CurrZAdj          string   `xml:"curr_Z_adj"`
	FreeSpace         string   `xml:"free_space"`
}

func DecodeCusp(file string) (*Cusp, error) {

	xmlFile, err := os.Open(file)
	if err != nil {
		return nil, nil
	}
	defer xmlFile.Close()

	d := xml.NewDecoder(xmlFile)
	d.CharsetReader = CharsetReader

	c := Cusp{}
	err = d.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Cusp) CheckIn() (time.Time, error) {
	t, err := time.Parse("Mon Jan 02 15:04:05 2006 ", c.SystemDateAndTime)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (c *Cusp) Temperature() (float64, error) {
	return strconv.ParseFloat(c.SysTemp, 32)
}
func (c *Cusp) Power() (float64, error) {
	return strconv.ParseFloat(c.SysBattVolt, 32)
}
func (c *Cusp) Storage() (int64, error) {
	return strconv.ParseInt(c.FreeSpace, 0, 32)
}
func (c *Cusp) SamplingRate() (float64, error) {
	return strconv.ParseFloat(c.SampRate, 32)
}
func (c *Cusp) PreEvent() (float64, error) {
	return strconv.ParseFloat(c.PreEvtLen, 32)
}
func (c *Cusp) PostEvent() (float64, error) {
	return strconv.ParseFloat(c.PostEvtLen, 32)
}
func (c *Cusp) ClockQuality() (float64, error) {
	return strconv.ParseFloat(c.TimeSuccessA, 32)
}
