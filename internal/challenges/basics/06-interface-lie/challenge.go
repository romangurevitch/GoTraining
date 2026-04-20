package interfacelie

import "fmt"

// AMLChecker performs anti-money-laundering checks.
// It implements the error interface so it can be returned as an error.
type AMLChecker struct {
	threshold int64
}

// Error implements the error interface.
func (c *AMLChecker) Error() string {
	return fmt.Sprintf("AML check failed: threshold %d exceeded", c.threshold)
}

// check simulates an AML check (always passes in this stub).
func (c *AMLChecker) check() error {
	return nil // real implementation would inspect transactions
}

// runAMLCheck runs an AML check unless skip is true.
//
// BUG: when skip is true, this returns a typed nil (*AMLChecker)(nil) as an error interface.
// The interface is (type=*AMLChecker, value=nil) — which is NOT == nil.
//
// TODO: make runAMLCheck return a genuinely nil error when skip is true.
func runAMLCheck(skip bool) error {
	if skip {
		return nil
	}
	checker := &AMLChecker{threshold: 1000000}
	return checker.check()
}

// Explain shows what the interface actually holds — useful for debugging.
func Explain(err error) string {
	return fmt.Sprintf("err=%v, type=%T, isNil=%v", err, err, err == nil)
}
