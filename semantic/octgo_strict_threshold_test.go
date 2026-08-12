package semantic

import "testing"

func TestStrictlyAboveOwnsEqualityInRemainder(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name             string
		value, threshold int
		want             bool
	}{
		{name: "above", value: 2, threshold: 1, want: true},
		{name: "equal", value: 1, threshold: 1, want: false},
		{name: "below", value: 0, threshold: 1, want: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := StrictlyAbove(test.value, test.threshold); got != test.want {
				t.Fatalf("StrictlyAbove(%d, %d) = %t, want %t", test.value, test.threshold, got, test.want)
			}
		})
	}
}
