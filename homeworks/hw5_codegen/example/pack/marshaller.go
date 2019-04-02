package main

import "encoding/binary"
import "bytes"

func (in *User) Unpack(data []byte) error {
	r := bytes.NewReader(data)

	// ID
	var IDRaw uint32
	binary.Read(r, binary.LittleEndian, &IDRaw)
	in.ID = int(IDRaw)

	// Login
	var LoginLenRaw uint32
	binary.Read(r, binary.LittleEndian, &LoginLenRaw)
	LoginRaw := make([]byte, LoginLenRaw)
	binary.Read(r, binary.LittleEndian, &LoginRaw)
	in.Login = string(LoginRaw)

	// Flags
	var FlagsRaw uint32
	binary.Read(r, binary.LittleEndian, &FlagsRaw)
	in.Flags = int(FlagsRaw)
	return nil
}
