package about

import "testing"

func TestVersion_GetCurrent(t *testing.T) {
	v := NewVersion()

	if v.GetCurrent() != version {
		t.Error("Version should match version variable in package")
	}
}
