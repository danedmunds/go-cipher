package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"unicode/utf8"
)

type RailFencerEncipher struct {
	index        int
	buffer       []byte
	rails        []io.ReadWriteCloser
	switchToRead func(io.ReadWriteCloser) (io.ReadWriteCloser, error)
	cleanUp      func(io.ReadWriteCloser) error
}

func NewRailFencerEncipher(rails []io.ReadWriteCloser) *RailFencerEncipher {
	return &RailFencerEncipher{
		rails:  rails,
		index:  0,
		buffer: make([]byte, 10, 10),
	}
}

func (rfe *RailFencerEncipher) Write(r rune) error {
	railIndex := selectRail(len(rfe.rails), rfe.index)
	rail := rfe.rails[railIndex]
	size := utf8.EncodeRune(rfe.buffer, r)
	_, err := rail.Write(rfe.buffer[0:size])
	rfe.index++
	return err
}

func selectRail(numRails, index int) int {
	var mod int
	if numRails == 1 {
		mod = 1
	} else {
		mod = (numRails - 1) * 2
	}

	return mod/2 - int(math.Abs(float64(index%mod-mod/2)))
}

func (rfe *RailFencerEncipher) DumpTo(writer io.Writer) (err error) {
	for _, reader := range rfe.rails {
		if rfe.switchToRead != nil {
			reader, err = rfe.switchToRead(reader)
			if err != nil {
				break
			}
		}

		defer reader.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			break
		}
	}

	return
}

func (rfe *RailFencerEncipher) CleanUp() (err error) {
	if rfe.cleanUp != nil {
		for _, reader := range rfe.rails {
			err = rfe.cleanUp(reader)
		}
	}

	return
}

func RailFence(numRails int) Cipher {
	rails := make([]io.ReadWriteCloser, 0, numRails)
	for i := 0; i < numRails; i++ {
		file, err := ioutil.TempFile(".", fmt.Sprintf("cipher_rail_%d_", i))
		rails = append(rails, file)
		if err != nil {
			panic(err)
		}
	}

	encipher := NewRailFencerEncipher(rails)
	encipher.switchToRead = func(in io.ReadWriteCloser) (out io.ReadWriteCloser, err error) {
		file := in.(*os.File)
		err = file.Close()
		if err != nil {
			return nil, err
		}

		return os.Open(file.Name())
	}
	encipher.cleanUp = func(rwc io.ReadWriteCloser) (err error) {
		file := rwc.(*os.File)
		return os.Remove(file.Name())
	}
	return NewFullMessageCipher(encipher)
}
