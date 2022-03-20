/*
	Copyright © 2022 Mikhail Batsin

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
	associated documentation files (the “Software”), to deal in the Software without restriction, including
	without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the
	following conditions:

	The above copyright notice and this permission notice shall be included in all copies or substantial
	portions of the Software.

	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
	LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO
	EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
	IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
	USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package gomemdb

import (
	"compress/zlib"
	"encoding/gob"
	"io"
	"os"
)

const (
	__GMDB_VERSION = "v0.1"
	__GMDB_LICENSE = "GPL-3.0"
)

// GoMemoryDatabase structure.
// Filename is name of file where to write db.
// KeyPairs are values.
// Compress. Set it to true if you want zlib compression.
type GoMemDb struct {
	Filename          string
	KeyPairs          []KeyPair
	Compress          bool
	ZlibCompressLevel int
}

// adds to keypairs new value
func (gmdb *GoMemDb) Add(key string, pair interface{}) {
	gmdb.KeyPairs = append(gmdb.KeyPairs, NewKP(key, pair))
}

// sets Compress to true
func (gmdb *GoMemDb) NeedCompress() {
	gmdb.Compress = true
}

// saves db state to file
func (gmdb *GoMemDb) Save() error {
	buf, err := os.OpenFile(gmdb.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	if gmdb.Compress {
		var out io.Writer

		out, err := zlib.NewWriterLevel(buf, gmdb.ZlibCompressLevel)

		if err != nil {
			return err
		}

		enc := gob.NewEncoder(out)
		err = enc.Encode(gmdb)

		if err != nil {
			return err
		}

		if c, ok := out.(io.Closer); ok {
			err = c.Close()

			if err != nil {
				return err
			}

			err = buf.Close()

			if err != nil {
				return err
			}
		}

		return nil
	}

	enc := gob.NewEncoder(buf)
	enc.Encode(gmdb)

	return nil
}

// opens saved db state to file and returns it
func Open(filename string, compressed bool) (GoMemDb, error) {
	var db GoMemDb
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)

	if compressed {
		rc, err := zlib.NewReader(file)

		if err != nil {
			return db, err
		}

		enc := gob.NewDecoder(rc)
		err = enc.Decode(&db)

		if err != nil {
			rc.Close()

			return GoMemDb{}, err
		}

		err = rc.Close()

		if err != nil {
			return db, err
		}

		err = file.Close()

		return db, err
	}

	enc := gob.NewDecoder(file)
	err = enc.Decode(&db)

	if err != nil {
		return db, err
	}

	err = file.Close()

	return db, err
}

// finds value by key
func (gmdb *GoMemDb) Get(key string) interface{} {
	if idx := gmdb.keyExists(key); idx != -1 {
		return gmdb.KeyPairs[idx]
	}

	return nil
}

// sets new value for key
func (gmdb *GoMemDb) Set(key string, value interface{}) {
	if idx := gmdb.keyExists(key); idx != -1 {
		gmdb.KeyPairs[idx].Pair = value
	}
}

// checks if key exists and returns it
func (gmdb *GoMemDb) keyExists(key string) int {
	for idx, kp := range gmdb.KeyPairs {
		if kp.Key == key {
			return idx
		}
	}

	return -1
}

// deletes key pair
func (gmdb *GoMemDb) Delete(key string) {
	if idx := gmdb.keyExists(key); idx != 1 {
		gmdb.KeyPairs = append(gmdb.KeyPairs[:idx], gmdb.KeyPairs[:idx+1]...)
	}
}

// returns map representation of all keypairs
func (gmdb *GoMemDb) Representate() map[string]interface{} {
	var repr map[string]interface{}

	for _, kp := range gmdb.KeyPairs {
		repr[kp.Key] = kp.Pair
	}

	return repr
}

// creates database instance
func NewDb(filename string) GoMemDb {
	return GoMemDb{Filename: filename, KeyPairs: []KeyPair{}}
}

// returns gomemdb version and it's license
func Info() (string, string) {
	return __GMDB_VERSION, __GMDB_LICENSE
}
