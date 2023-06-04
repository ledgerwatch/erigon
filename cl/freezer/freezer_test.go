package freezer_test

import (
	"os"
	"testing"

	"github.com/ledgerwatch/erigon/cl/freezer"
	"github.com/stretchr/testify/assert"
)

func runBlobStoreTest(t *testing.T, b *freezer.BlobStore) {
	var err error
	// put bad item into obj
	err = b.Put(nil, "../../../test", "a", "b")
	assert.ErrorIs(t, err, os.ErrInvalid)
	// get bad item
	err = b.Put(nil, "../../../test", "a", "b")
	assert.ErrorIs(t, err, os.ErrInvalid)
	// put item into obj
	orig := []byte{1, 2, 3, 4}
	err = b.Put(orig, "test", "a", "b")
	assert.NoError(t, err)

	// get item from obj
	ans, err := b.Get("test", "a", "b")
	assert.NoError(t, err)
	assert.EqualValues(t, orig, ans)

	ans, err = b.Get("test", "b", "a")
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Nil(t, ans)
}

func runSidecarBlobStoreTest(t *testing.T, b *freezer.SidecarBlobStore) {
	var err error
	// put bad item into obj
	err = b.Put(nil, nil, "../../../test", "a", "b")
	assert.ErrorIs(t, err, os.ErrInvalid)
	// get bad item
	err = b.Put(nil, nil, "../../../test", "a", "b")
	assert.ErrorIs(t, err, os.ErrInvalid)
	// put item into obj
	orig := []byte{1, 2, 3, 4}
	orig2 := []byte{5, 6, 7, 8}
	err = b.Put(orig, orig2, "test", "a", "b")
	assert.NoError(t, err)

	// get item from obj
	ans, sidecar, err := b.Get("test", "a", "b")
	assert.NoError(t, err)
	assert.EqualValues(t, orig, ans)
	assert.EqualValues(t, orig2, sidecar)

	ans, sidecar, err = b.Get("test", "b", "a")
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Nil(t, ans)

	// put item without sidecar
	err = b.Put(orig2, nil, "test", "a", "c")
	assert.NoError(t, err)

	// get item from obj
	ans, sidecar, err = b.Get("test", "a", "c")
	assert.NoError(t, err)
	assert.EqualValues(t, orig2, ans)
	assert.Nil(t, sidecar)
}

func testFreezer(t *testing.T, fn func() (freezer.Freezer, func())) {
	t.Run("BlobStore", func(t *testing.T) {
		runBlobStoreTest(t, freezer.NewBlobStore(&freezer.InMemory{}))
	})
	t.Run("SidecarBlobStore", func(t *testing.T) {
		runSidecarBlobStoreTest(t, freezer.NewSidecarBlobStore(&freezer.InMemory{}))
	})
}

func TestMemoryStore(t *testing.T) {
	testFreezer(t, func() (freezer.Freezer, func()) {
		return &freezer.InMemory{}, nil
	})
}

func TestRootPathStore(t *testing.T) {
	testFreezer(t, func() (freezer.Freezer, func()) {
		return &freezer.RootPathOsFs{"test_output"}, func() { os.RemoveAll("test_output") }
	})
}
