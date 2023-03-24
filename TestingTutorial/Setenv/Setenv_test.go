package setenv_test

import (
 	"os/exec"
    "testing"
	"path\filepath"
}

//テスト終了とともに環境変数が基に戻される
func TestCreateProfile(t *testing.T)  {
	t.Setenv("DATABASE_URL")
	err := doSomething()
	if err!= nil {
      t.Fatalf("cannot do something: %v", err)
    }
}

