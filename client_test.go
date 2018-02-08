// +build unittest

// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// test client for ots2
package goots

import (
	"fmt"
	. "github.com/GiterLab/goots/otstype"
	"github.com/GiterLab/goots/urllib"
	"net"
	"net/http"
	"os"
	"testing"
)

func newWithRetry(tableName string) (*OTSClient, error) {
	host := os.Getenv("OTS_TEST_HOST")
	region := os.Getenv("OTS_TEST_REGION")
	instanceName := os.Getenv("OTS_TEST_INSTANCE")
	endpoint := fmt.Sprintf("https://%s.%s.%s", instanceName, region, host)
	akID := os.Getenv("OTS_TEST_ACCESS_KEY_ID")
	akSec := os.Getenv("OTS_TEST_ACCESS_KEY_SECRET")
	o, err := NewWithRetryPolicy(endpoint, akID, akSec, instanceName, nil)
	o.DeleteTable(tableName)
	if err != nil {
		return nil, err
	}
	tableMeta := &OTSTableMeta{
		TableName: tableName,
		SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
			{K: "gid", V: "INTEGER"},
			{K: "uid", V: "INTEGER"},
		},
	}
	reservedThroughput := &OTSReservedThroughput{
		OTSCapacityUnit{0, 0},
	}
	if e := o.CreateTable(tableMeta, reservedThroughput); e != nil {
		return o, e
	}
	return o, nil
}

// Test_NewWithRetry_OnNetError should succ...
func Test_NewWithRetry_OnNetError(t *testing.T) {
	tname := "myTestTable"
	client, err := newWithRetry(tname)
	if err != nil {
		t.Logf("create %s failed due to %s", tname, err)
		t.Fail()
		return
	}
	defer func() {
		err := client.DeleteTable(tname)
		if err != nil {
			t.Errorf("delete table `%s` fail due to %s", tname, err)
		}
	}()

	urllib.MockResp.Reset()
	urllib.MockResp.MockFunc = func(m *urllib.MockResponse, b *urllib.HttpRequest) (*http.Response, error) {
		if m.Count < 2 {
			return m.MockError(new(net.DNSError))
		}
		return b.GetResponse()
	}
	tblList, err := client.ListTable()
	if urllib.MockResp.Count != 3 {
		t.Errorf("number of request excpected %d but %d", 3, urllib.MockResp.Count)
	}
	if err != (*OTSError)(nil) {
		t.Errorf("list table fail due to %s", err)
		return
	}
	foundTbl := false
	for _, name := range tblList.TableNames {
		if name == tname {
			foundTbl = true
			break
		}
	}
	if !foundTbl {
		t.Errorf("not found %s on ListTable", tname)
	}
}

// Test_NewWithRetry_OnUnkownError should succ...
func Test_NewWithRetry_OnUnkownError(t *testing.T) {
	tname := "myTestTable"
	client, err := newWithRetry(tname)
	if err != nil {
		t.Logf("create %s failed due to %s", tname, err)
		t.Fail()
		return
	}
	defer func() {
		err := client.DeleteTable(tname)
		if err != nil {
			t.Errorf("delete table `%s` fail due to %s", tname, err)
		}
	}()

	urllib.MockResp.Reset()
	urllib.MockResp.MockFunc = func(m *urllib.MockResponse, b *urllib.HttpRequest) (*http.Response, error) {
		if m.Count < 1 {
			return m.MockError(fmt.Errorf("unkown error by mock"))
		}
		return b.GetResponse()
	}
	_, err = client.ListTable()
	if urllib.MockResp.Count != 1 {
		t.Errorf("number of request excpected %d but %d", 1, urllib.MockResp.Count)
	}
	if err == (*OTSError)(nil) {
		t.Errorf("list table should fail")
		return
	}
}

func Test_NewWithRetry_BodyError(t *testing.T) {
	tname := "myTestTable"
	client, err := newWithRetry(tname)
	if err != nil {
		t.Logf("create %s failed due to %s", tname, err)
		t.Fail()
		return
	}
	defer func() {
		err := client.DeleteTable(tname)
		if err != nil {
			t.Errorf("delete table `%s` fail due to %s", tname, err)
		}
	}()

	expectedCount := 3

	urllib.MockResp.Reset()
	urllib.MockResp.MockFunc = func(m *urllib.MockResponse, b *urllib.HttpRequest) (*http.Response, error) {
		if m.Count < expectedCount-1 {
			return m.MockReadBodyError()
		}
		return b.GetResponse()
	}
	_, err = client.ListTable()
	if urllib.MockResp.Count != expectedCount {
		t.Errorf("number of request excpected %d but %d", expectedCount, urllib.MockResp.Count)
	}
	if err != (*OTSError)(nil) {
		t.Errorf("list table fail due to %s", err)
		return
	}
}

func Test_NewWithRetry_NoBody(t *testing.T) {
	tname := "myTestTable"
	client, err := newWithRetry(tname)
	if err != nil {
		t.Logf("create %s failed due to %s", tname, err)
		t.Fail()
		return
	}
	defer func() {
		err := client.DeleteTable(tname)
		if err != nil {
			t.Errorf("delete table `%s` fail due to %s", tname, err)
		}
	}()

	expectedCount := 3

	urllib.MockResp.Reset()
	urllib.MockResp.MockFunc = func(m *urllib.MockResponse, b *urllib.HttpRequest) (*http.Response, error) {
		if m.Count < expectedCount-1 {
			return m.MockNoBody()
		}
		return b.GetResponse()
	}
	_, err = client.ListTable()
	if urllib.MockResp.Count != expectedCount {
		t.Errorf("number of request excpected %d but %d", expectedCount, urllib.MockResp.Count)
	}
	if err != (*OTSError)(nil) {
		t.Errorf("list table fail due to %s", err)
		return
	}
}

func Test_NewWithRetry_5xxError(t *testing.T) {
	OTSErrorPanicMode = false
	tname := "myTestTable"
	client, err := newWithRetry(tname)
	if err != nil {
		t.Logf("create %s failed due to %s", tname, err)
		t.Fail()
		return
	}
	defer func() {
		err := client.DeleteTable(tname)
		if err != nil {
			t.Errorf("delete table `%s` fail due to %s", tname, err)
		}
	}()

	expectedCount := 3

	urllib.MockResp.Reset()
	urllib.MockResp.MockFunc = func(m *urllib.MockResponse, b *urllib.HttpRequest) (*http.Response, error) {
		if m.Count < expectedCount-1 {
			return m.Mock5xxError()
		}
		return b.GetResponse()
	}
	_, err = client.ListTable()
	if urllib.MockResp.Count != expectedCount {
		t.Errorf("number of request excpected %d but %d", expectedCount, urllib.MockResp.Count)
	}
	if err != (*OTSError)(nil) {
		t.Errorf("list table fail due to %s", err)
		return
	}
}

func Test_New(t *testing.T) {
	o, err := New("http://127.0.0.1:8800", "OTSMultiUser177_accessid", "OTSMultiUser177_accesskey", "TestInstance177")
	if err != nil {
		t.Fail()
	}
	t.Log(o)

	o, err = New("http://127.0.0.1:8800", "OTSMultiUser177_accessid", "OTSMultiUser177_accesskey", "TestInstance177",
		60, 60, "ots2-client-test", "utf-8")
	if err != nil {
		t.Fail()
	}
	t.Log(o)
	// t.Fail()
}

func Test_Set(t *testing.T) {
	o, err := New("http://127.0.0.1:8800", "OTSMultiUser177_accessid", "OTSMultiUser177_accesskey", "TestInstance177")
	if err != nil {
		t.Fail()
	}
	t.Log(o)

	o = o.Set(DictString{
		"EndPoint": "http://127.0.0.1:8888",
		// "NotExist": 123,
	})

	t.Log(o)
}
