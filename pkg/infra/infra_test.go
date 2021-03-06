package infra

import (
	"context"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"

	tcErr "github.com/pddg/tiny-cluster/pkg/errors"
)

type testFixture interface {
	prepare(ctx context.Context, t *testing.T, client *clientv3.Client)
	clean(ctx context.Context, t *testing.T, client *clientv3.Client)
}

const testEtcdEndpointsKey = "TC_ETCD_ENDPOINTS"

func getTestEndpoints(t *testing.T) []string {
	t.Helper()
	urlsStr := os.Getenv(testEtcdEndpointsKey)
	if len(urlsStr) == 0 {
		t.Errorf("%s does not specified", testEtcdEndpointsKey)
	}
	return strings.Split(urlsStr, ",")
}

func getTestClient(t *testing.T) *clientv3.Client {
	t.Helper()
	urls := getTestEndpoints(t)
	client, err := clientv3.NewFromURLs(urls)
	if err != nil {
		t.Errorf("failed to get client due to %v", err)
	}
	return client
}

func setUpTest(ctx context.Context, t *testing.T, client *clientv3.Client, fixtures testFixture) {
	t.Helper()
	if fixtures == nil {
		return
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	fixtures.prepare(timeoutCtx, t, client)
}

func tearDownTest(ctx context.Context, t *testing.T, client *clientv3.Client, fixtures testFixture) {
	t.Helper()
	if fixtures == nil {
		return
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	fixtures.clean(timeoutCtx, t, client)
}

type testFixtureImpl map[string]string

func (f *testFixtureImpl) prepare(ctx context.Context, t *testing.T, client *clientv3.Client) {
	t.Helper()
	for k, v := range *f {
		_, err := client.Put(ctx, k, v)
		if err != nil {
			t.Errorf("Failed to put value due to %v", err)
		}
	}
}

func (f *testFixtureImpl) clean(ctx context.Context, t *testing.T, client *clientv3.Client) {
	t.Helper()
	for k := range *f {
		_, err := client.Delete(ctx, k)
		if err != nil {
			t.Errorf("Failed to delete value due to %v", err)
		}
	}
}

func Test_doGet(t *testing.T) {
	testCases := map[string]struct {
		fixtures    *testFixtureImpl
		key         string
		expected    string
		expectedErr error
	}{
		"get exists value": {
			fixtures:    &testFixtureImpl{"key": "value"},
			key:         "key",
			expected:    "value",
			expectedErr: nil,
		},
		"not found": {
			fixtures:    &testFixtureImpl{},
			key:         "notfound",
			expectedErr: tcErr.ErrNotFound,
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
			actualValue := string(actual)
			if actualValue != tc.expected {
				t.Errorf("Obtained value is invalid. Expected: %#v, Actual: %#v", tc.expected, actualValue)
			}
		})
	}
}

func Test_doGetWithRev(t *testing.T) {
	testCases := map[string]struct {
		fixtures         *testFixtureImpl
		key              string
		expected         string
		expectedRevision int64
		expectedErr      error
	}{
		"get exists value with revision": {
			fixtures:         &testFixtureImpl{"key": "value"},
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
			defer tearDownTest(ctx, t, client, tc.fixtures)
			_, initRev, err := doGetWithRev(ctx, client, tc.key)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
			for i = 0; i < tc.expectedRevision; i++ {
				setUpTest(ctx, t, client, tc.fixtures)
			}
			actual, actualRev, err := doGetWithRev(ctx, client, tc.key)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
			if (actualRev - initRev) != tc.expectedRevision {
				t.Errorf("Obtained revision is invalid. Expected: %d, Actual: %d", tc.expectedRevision, actualRev)
			}
			actualValue := string(actual)
			if actualValue != tc.expected {
				t.Errorf("Obtained value is invalid. Expected: %#v, Actual: %#v", tc.expected, actualValue)
			}
		})
	}
}

func Test_doGetAll(t *testing.T) {
	testCases := map[string]struct {
		fixtures    *testFixtureImpl
		key         string
		expected    []string
		expectedErr error
	}{
		"get an item": {
			fixtures:    &testFixtureImpl{"key": "value"},
			key:         "key",
			expected:    []string{"value"},
			expectedErr: nil,
		},
		"get multiple item": {
			fixtures:    &testFixtureImpl{"key": "value", "key/key1": "value1", "key/key2": "value2"},
			key:         "key",
			expected:    []string{"value", "value1", "value2"},
			expectedErr: nil,
		},
		"not found": {
			fixtures:    &testFixtureImpl{},
			key:         "key",
			expected:    []string(nil),
			expectedErr: nil,
		},
	}
	for tn, tc := range testCases {
		ctx := context.Background()
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			values, err := doGetAll(ctx, client, tc.key)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
			var stringValues []string
			for _, v := range values {
				stringValues = append(stringValues, string(v))
			}
			sort.Strings(stringValues)
			sort.Strings(tc.expected)
			if !reflect.DeepEqual(stringValues, tc.expected) {
				t.Errorf("Obtained value is invalid. Expected: %#v, Actual: %#v", tc.expected, stringValues)
			}
		})
	}
}

func Test_doCreate(t *testing.T) {
	testCases := map[string]struct {
		fixtures    *testFixtureImpl
		key         string
		expected    string
		expectedErr error
	}{
		"create new item": {
			fixtures:    &testFixtureImpl{},
			key:         "key",
			expected:    "value",
			expectedErr: nil,
		},
		"create exists item": {
			fixtures:    &testFixtureImpl{"key": "value"},
			key:         "key",
			expected:    "value",
			expectedErr: tcErr.ErrAlreadyExists,
		},
	}
	for tn, tc := range testCases {
		ctx := context.Background()
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			err := doCreate(ctx, client, tc.key, tc.expected)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
		})
	}
}

func Test_doUpdate(t *testing.T) {
	testCases := map[string]struct {
		fixtures        *testFixtureImpl
		key             string
		delBeforeUpdate bool
		putBeforeUpdate bool
		expected        string
		expectedErr     error
	}{
		"update successfully": {
			fixtures:        &testFixtureImpl{"key": "value"},
			key:             "key",
			delBeforeUpdate: false,
			putBeforeUpdate: false,
			expected:        "value",
			expectedErr:     nil,
		},
		"update retry": {
			fixtures:        &testFixtureImpl{"key": "value"},
			key:             "key",
			delBeforeUpdate: false,
			putBeforeUpdate: true,
			expected:        "value",
			expectedErr:     nil,
		},
		"delete before update": {
			fixtures:        &testFixtureImpl{"key": "value"},
			key:             "key",
			delBeforeUpdate: true,
			putBeforeUpdate: false,
			expected:        "value",
			expectedErr:     tcErr.ErrNotFound,
		},
	}
	for tn, tc := range testCases {
		ctx := context.Background()
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			resp, err := client.Get(ctx, tc.key)
			if err != nil {
				t.Errorf("Failed to get the value due to %v", err)
			}
			rev := resp.Kvs[0].ModRevision
			if tc.delBeforeUpdate {
				if _, err := client.Delete(ctx, tc.key); err != nil {
					t.Errorf("Failed to delete the key due to %v", err)
				}
			}
			if tc.putBeforeUpdate {
				if _, err := client.Put(ctx, tc.key, tc.expected); err != nil {
					t.Errorf("Failed to update the key due to %v", err)
				}
			}
			err = doUpdate(ctx, client, rev, tc.key, tc.expected)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
		})
	}
}

func Test_doDelete(t *testing.T) {
	testCases := map[string]struct {
		fixtures    *testFixtureImpl
		key         string
		expectedErr error
	}{
		"delete exist item": {
			fixtures:    &testFixtureImpl{"key": "value"},
			key:         "key",
			expectedErr: nil,
		},
		"delete item which does not exist": {
			fixtures:    &testFixtureImpl{},
			key:         "key",
			expectedErr: tcErr.ErrNotFound,
		},
	}
	for tn, tc := range testCases {
		ctx := context.Background()
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			err := doDelete(ctx, client, tc.key)
			if err != tc.expectedErr {
				t.Errorf("Error type is invalid. Expected: %v, Actual: %v", tc.expectedErr, err)
			}
		})
	}
}
