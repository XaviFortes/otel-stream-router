package routerprocessor

import "testing"

func TestBuildStreamName(t *testing.T) {
	tests := []struct {
		name          string
		defaultStream string
		namespace     string
		business      string
		environment   string
		nodepool      string
		expected      string
	}{
		{
			name:          "namespace business and environment",
			defaultStream: "default",
			namespace:     "payments",
			business:      "PGA",
			environment:   "p",
			nodepool:      "",
			expected:      "payments-PGA-p",
		},
		{
			name:          "fallback to default stream",
			defaultStream: "default",
			namespace:     "",
			business:      "unknown",
			environment:   "unknown",
			nodepool:      "",
			expected:      "default",
		},
		{
			name:          "skip unknown business",
			defaultStream: "default",
			namespace:     "payments",
			business:      "unknown",
			environment:   "p",
			nodepool:      "",
			expected:      "payments-p",
		},
		{
			name:          "include nodepool",
			defaultStream: "default",
			namespace:     "payments",
			business:      "PGA",
			environment:   "p",
			nodepool:      "np-1",
			expected:      "payments-PGA-p-np-1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildStreamName(tc.defaultStream, tc.namespace, tc.business, tc.environment, tc.nodepool)
			if got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}