package main

import (
	"encoding/json"
	"os"
)

// Config for client
type Config struct {
	LocalAddr    string `json:"localaddr"`
	RemoteAddr   string `json:"remoteaddr"`
	Key          string `json:"key"`
	Crypt        string `json:"crypt" ui:"select:aes,aes-128,aes-192,salsa20,blowfish,twofish,cast5,3des,tea,xtea,xor,none"`
	Mode         string `json:"mode" ui:"select:fast,fast2,fast3,normal"`
	Conn         int    `json:"conn"`
	AutoExpire   int    `json:"autoexpire"`
	MTU          int    `json:"mtu"`
	SndWnd       int    `json:"sndwnd"`
	RcvWnd       int    `json:"rcvwnd"`
	DataShard    int    `json:"datashard"`
	ParityShard  int    `json:"parityshard"`
	DSCP         int    `json:"dscp"`
	NoComp       bool   `json:"nocomp"`
	AckNodelay   bool   `json:"acknodelay"`
	NoDelay      int    `json:"nodelay"`
	Interval     int    `json:"interval"`
	Resend       int    `json:"resend"`
	NoCongestion int    `json:"nc"`
	SockBuf      int    `json:"sockbuf"`
	KeepAlive    int    `json:"keepalive"`
	Log          string `json:"log"`
	SnmpLog      string `json:"snmplog"`
	SnmpPeriod   int    `json:"snmpperiod"`
}

func parseJSONConfig(config *Config, path string) error {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(config)
}

func saveJsonConfig(config *Config, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644) // For read access.
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "  ")
	return encoder.Encode(config)
}
