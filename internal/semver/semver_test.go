package semver

import "testing"

func TestValid(t *testing.T) {
	for _, value := range []string{"0.1.0", "1.2.3-rc.1+build.001", "1.0.0-0", "999999999999999999999.0.0"} {
		if !Valid(value) {
			t.Errorf("expected valid version: %s", value)
		}
	}
}

func TestInvalid(t *testing.T) {
	for _, value := range []string{"", "v1.0.0", "1.2", "01.2.3", "1.2.3-01", "1.2.3-", "1.2.3+", "1.2.3-a..b", "1.2.3\n", "1.2.3-ä"} {
		if Valid(value) {
			t.Errorf("expected invalid version: %q", value)
		}
	}
}

func TestPrecedence(t *testing.T) {
	versions := []string{"1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.beta", "1.0.0-beta", "1.0.0-beta.2", "1.0.0-beta.11", "1.0.0-rc.1", "1.0.0", "1.0.1", "1.10.0", "2.0.0", "999999999999999999999.0.0"}
	for i := 1; i < len(versions); i++ {
		if Compare(versions[i-1], versions[i]) >= 0 {
			t.Errorf("expected %s < %s", versions[i-1], versions[i])
		}
	}
}

func TestBuildMetadataDoesNotAffectPrecedence(t *testing.T) {
	if Compare("1.0.0+one", "1.0.0+two") != 0 {
		t.Fatal("build metadata must not affect precedence")
	}
}
