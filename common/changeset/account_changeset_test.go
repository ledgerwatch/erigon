package changeset

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon/common"
	"github.com/stretchr/testify/assert"
)

func TestEncodingAccount(t *testing.T) {
	bkt := kv.AccountChangeSet
	m := Mapper[bkt]

	ch := m.New()
	// empty StorageChangeSset first
	err := m.Encode(1, ch, func(k, v []byte) error {
		return fmt.Errorf("must never call")
	})
	assert.NoError(t, err)

	vals := [][]byte{
		common.FromHex("f7f6db1eb17c6d582078e0ffdd0c"),
		common.FromHex("b1e9b5c16355eede662031dd621d08faf4ea"),
		common.FromHex("862cf52b74f1cea41ddd8ffa4b3e7c7790"),
	}
	numOfElements := 3
	for i := 0; i < numOfElements; i++ {
		address := common.HexToAddress(fmt.Sprintf("0xBe828AD8B538D1D691891F6c725dEdc5989abBc%d", i))
		err2 := ch.Add(address[:], vals[i])
		if err2 != nil {
			t.Fatal(err)
		}
	}

	ch2 := m.New()
	err = m.Encode(1, ch, func(k, v []byte) error {
		var err error
		_, k, v, err = m.Decode(k, v)
		if err != nil {
			return err
		}
		return ch2.Add(k, v)
	})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ch, ch2) {
		fmt.Println("ch", len(ch.Changes), "ch2", len(ch2.Changes))
		for i, v := range ch.Changes {
			fmt.Println("Line ", i)

			if !bytes.Equal(v.Key, ch2.Changes[i].Key) || !bytes.Equal(v.Value, ch2.Changes[i].Value) {
				fmt.Println("Diff ", i)
				fmt.Println("k1", common.Bytes2Hex(v.Key), len(v.Key))
				fmt.Println("k2", common.Bytes2Hex(ch2.Changes[i].Key))
				fmt.Println("v1", common.Bytes2Hex(v.Value))
				fmt.Println("v2", common.Bytes2Hex(ch2.Changes[i].Value))
			}
		}
		fmt.Printf("%+v %+v\n", ch, ch2)
		t.Fatal("not equal")
	}
}
