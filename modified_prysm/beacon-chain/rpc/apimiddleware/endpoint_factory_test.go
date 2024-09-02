package apimiddleware_test

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/rpc/apimiddleware"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

func TestBeaconEndpointFactory_AllPathsRegistered(t *testing.T) {
	f := &apimiddleware.BeaconEndpointFactory{}

	for _, p := range f.Paths() {
		_, err := f.Create(p)
		require.NoError(t, err, "failed to register %s", p)
	}
}
