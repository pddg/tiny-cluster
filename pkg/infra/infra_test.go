package infra

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const testEtcdEndpointsKey = "TC_ETCD_ENDPOINTS"

type testFixtures map[string]string

func (f *testFixtures) Prepare(t *testing.T, ctx context.Context, client *clientv3.Client) {
	t.Helper()
	for k, v := range *f {
		_, err := client.Put(ctx, k, v)
		if err != nil {
			t.Errorf("Failed to put value due to %v", err)
		}
	}
}

func (f *testFixtures) Clean(t *testing.T, ctx context.Context, client *clientv3.Client) error {
	for k, _ := range *f {
		_, err := client.Delete(ctx, k)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTestClient(t *testing.T) *clientv3.Client {
	t.Helper()
	urlsStr := os.Getenv(testEtcdEndpointsKey)
	if len(urlsStr) == 0 {
		t.Errorf("%s does not specified", testEtcdEndpointsKey)
	}
	urls := strings.Split(urlsStr, ",")
	client, err := clientv3.NewFromURLs(urls)
	if err != nil {
		t.Errorf("failed to get client due to %v", err)
	}
	return client
}

func setUpTest(ctx context.Context, t *testing.T, client *clientv3.Client, fixtures testFixtures) {
	t.Helper()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	fixtures.Prepare(t, timeoutCtx, client)
}

func tearDownTest(ctx context.Context, t *testing.T, client *clientv3.Client, fixtures testFixtures) {
	t.Helper()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	fixtures.Clean(t, timeoutCtx, client)
}

func Test_doGet(t *testing.T) {
	testCases := map[string]struct {
		fixtures    testFixtures
		key         string
		expected    string
		expectedErr error
	}{
		"get exists value": {
			fixtures:    testFixtures{"key": "value"},
			key:         "key",
			expected:    "value",
			expectedErr: nil,
		},
		"not found": {
			key:         "notfound",
			expectedErr: ErrNotFound,
		},
	}
	for tn, tc := range testCases {
		ctx := context.Background()
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			actual, err := doGet(context.TODO(), client, tc.key)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
			if actual != tc.expected {
				t.Errorf("Obtained value is invalid. Expected: %s, Actual: %s", tc.expected, actual)
			}
		})
	}
}

func Test_doGetWithRev(t *testing.T) {
	testCases := map[string]struct {
		fixtures         testFixtures
		key              string
		expected         string
		expectedRevision int64
		expectedErr      error
	}{
		"get exists value with revision": {
			fixtures:         testFixtures{"key": "value"},
			key:              "key",
			expected:         "value",
			expectedRevision: 3,
			expectedErr:      nil,
		},
	}
	for tn, tc := range testCases {
		ctx := context.Background()
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			var i int64
			setUpTest(ctx, t, client, tc.fixtures)
			_, initRev, err := doGetWithRev(ctx, client, tc.key)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
			for i = 0; i < tc.expectedRevision; i++ {
				setUpTest(ctx, t, client, tc.fixtures)
			}
			actualValue, actualRev, err := doGetWithRev(ctx, client, tc.key)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
			if (actualRev - initRev) != tc.expectedRevision {
				t.Errorf("Obtained revision is invalid. Expected: %d, Actual: %d", tc.expectedRevision, actualRev)
			}
			if actualValue != tc.expected {
				t.Errorf("Obtained value is invalid. Expected: %s, Actual: %s", tc.expected, actualValue)
			}
			tearDownTest(ctx, t, client, tc.fixtures)
		})
	}
}

func Test_doCreate(t *testing.T) {
	testCases := map[string]struct {
		fixtures    testFixtures
		key         string
		expected    string
		expectedErr error
	}{
		"create new item": {
			key:         "key",
			expected:    "value",
			expectedErr: nil,
		},
		"create exists item": {
			fixtures:    testFixtures{"key": "value"},
			key:         "key",
			expected:    "value",
			expectedErr: ErrAlreadyExists,
		},
	}
	for tn, tc := range testCases {
		ctx := context.Background()
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			err := doCreate(ctx, client, tc.key, tc.expected)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
			tearDownTest(ctx, t, client, tc.fixtures)
		})
	}
}
