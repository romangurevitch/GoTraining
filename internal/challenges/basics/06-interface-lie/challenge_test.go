package interfacelie

import "testing"

func TestRunAMLCheck_SkipReturnsNil(t *testing.T) {
	err := runAMLCheck(true)
	if err != nil {
		t.Fatalf(
			"runAMLCheck(skip=true) should return nil error, got: %v (%T)\n"+
				"  Hint: a typed nil (*AMLChecker)(nil) stored in an error interface is NOT nil.\n"+
				"  Fix: return bare nil instead of the typed variable.",
			err, err,
		)
	}
}

func TestRunAMLCheck_NoSkipReturnsNil(t *testing.T) {
	err := runAMLCheck(false)
	if err != nil {
		t.Fatalf("runAMLCheck(skip=false) should return nil for a passing check, got: %v", err)
	}
}
