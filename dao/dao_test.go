package dao

import (
	"os"
	"testing"
)

var d *Dao

func TestMain(m *testing.M) {
	d = New("test.db")
	os.Exit(m.Run())
}

func Test_Sync(t *testing.T) {
	d.Sync2()
}
