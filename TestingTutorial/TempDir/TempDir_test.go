package tempdir_test

import (
	"path/filepath"
	"testing"
)

// 終了すると自動でテンポラリディレクトリが削除される
func TestCreateProfile(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.json")
	got, err := CreateProfile(filename) //定義していない
	if err != nil {
		t.Fatal(err)
	}
	want := true
	if got != want {
		t.Fatalf("want %v, but %v", want, got)
	}
}
