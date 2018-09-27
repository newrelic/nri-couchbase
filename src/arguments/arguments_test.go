package arguments

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateArguments(t *testing.T) {
	testCases := []struct{
		args ArgumentList
		expectedError bool
	}{
		{
			ArgumentList{
				Username: "testuser",
				Password: "testpass",
				Hostname: "testhost",
				Port: 8901,
				QueryPort: 8093,
				UseSSL: true,
				CABundleDir: "test-dir",
				CABundleFile: "test-file",
			},
			false,
		},
		{
			ArgumentList{
				Username: "",
				Password: "testpass",
				Hostname: "testhost",
				Port: 8091,
				QueryPort: 8093,
			},
			true,
		},
		{
			ArgumentList{
				Username: "testuser",
				Password: "",
				Hostname: "testhost",
				Port: 8091,
				QueryPort: 8093,
			},
			true,
		},
		{
			ArgumentList{
				Username: "testuser",
				Password: "testpass",
				Hostname: "",
				Port: 8091,
				QueryPort: 8093,
			},
			true,
		},
		{
			ArgumentList{
				Username: "testuser",
				Password: "testpass",
				Hostname: "testhost",
				Port: 7676767,
				QueryPort: 8093,
			},
			true,
		},
		{
			ArgumentList{
				Username: "testuser",
				Password: "testpass",
				Hostname: "testhost",
				Port: 8091,
				QueryPort: 8181818,
			},
			true,
		},
		{
			ArgumentList{
				Username: "testuser",
				Password: "testpass",
				Hostname: "testhost",
				Port: 8091,
				QueryPort: 8093,
				UseSSL: true,
				CABundleDir: "",
				CABundleFile: "",
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.args.Validate()
		assert.Equal(t, tc.expectedError, err != nil)
	}
}