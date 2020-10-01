package jsii

import (
	"sync"
)

var (
	clientInstance      *client
	clientInstanceMutex sync.Mutex
	once                sync.Once
)

// getClient returns a singleton client instance, initializing one the first
// time it is called.
func getClient() *client {
	once.Do(func() {
		// Locking early to be safe with a concurrent Close execution
		clientInstanceMutex.Lock()
		defer clientInstanceMutex.Unlock()

		client, err := newClient()
		if err != nil {
			panic(err)
		}

		clientInstance = client
	})

	return clientInstance
}

// Close finalizes the runtime process, signalling the end of the execution to
// the jsii kernel process, and waiting for graceful termination. The best
// practice is to defer call thins at the beginning of the "main" function.
//
// If a jsii client is used *after* Close was called, a new jsii kernel process
// will be initialized, and Close should be called again to correctly finalize
// that, too. This behavior is intended for use in unit/integration tests.
func Close() {
	// Locking early to be safe with a concurrent getClient execution
	clientInstanceMutex.Lock()
	defer clientInstanceMutex.Unlock()

	// Reset the "once" so a new client would get initialized next time around
	once = sync.Once{}

	if clientInstance != nil {
		// Close the client & reset it
		clientInstance.close()
		clientInstance = nil
	}
}