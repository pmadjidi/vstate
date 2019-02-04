package main

import (
	"encoding/json"
	"fmt"
	"gotest.tools/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"vehicles/vstate"
)



func Test001(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	close(app.start) //Obs only once for the entire test suit!!

	tests := []struct {
		name string
		req    *http.Request
	}{
		{name: "1: testing listv", req: newreq("GET", ts.URL+"/admin/listv", nil)},
		{name: "2: testing newv", req: newreq("GET", ts.URL+"/admin/newv", nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, 200, resp.StatusCode, "OK response is expected")
		})
	}
}




func Test002(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing claim user", req: newreq("GET", ts.URL+"/user/claim/1GcsahF1mmbeX2y4uCgf96HISba" , nil),expect: 200,state: vstate.Riding},
		{name: "2: testing claim admin", req: newreq("GET", ts.URL+"/admin/claim/1GcsahF1mmbeX2y4uCgf96HISbb", nil),expect:200,state: vstate.Riding},
		{name: "3: testing claim hunter", req: newreq("GET", ts.URL+"/hunter/claim/1GcsahF1mmbeX2y4uCgf96HISbc", nil),expect:200,state: vstate.Riding},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			var v Vehicle
			json.NewDecoder(resp.Body).Decode(&v)
			assert.Equal(t,v.State,tt.state,"Veachle should be in state riding")
		})
	}
}


func Test003(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{   // allready claimed
		{name: "1: testing claim user", req: newreq("GET", ts.URL+"/user/claim/1GcsahF1mmbeX2y4uCgf96HISba" , nil),expect: http.StatusForbidden,state: vstate.Nothing},
		{name: "2: testing claim admin", req: newreq("GET", ts.URL+"/admin/claim/1GcsahF1mmbeX2y4uCgf96HISbb", nil),expect: http.StatusForbidden,state: vstate.Nothing},
		{name: "3: testing claim hunter", req: newreq("GET", ts.URL+"/hunter/claim/1GcsahF1mmbeX2y4uCgf96HISbc", nil),expect: http.StatusForbidden,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}


func Test004(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing admin setstate ready", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/ready" , nil),expect: 200,state: vstate.Ready},
		{name: "2: testing admin setstate Battery_low ", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISbb/Battery_low", nil),expect:200,state: vstate.Battery_low},
		{name: "3: testing admin setstate Collected", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISbc/Collected", nil),expect:200,state: vstate.Collected},
		{name: "4: testing admin setstate Unknown", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISbd/Unknown", nil),expect:200,state: vstate.Unknown},
		{name: "5: testing admin setstate Bounty", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISbe/Bounty", nil),expect:200,state: vstate.Bounty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
			var v Vehicle
			json.NewDecoder(resp.Body).Decode(&v)
			assert.Equal(t,v.State,tt.state,"Veachle should be in state ready")
		})
	}
}



func Test005(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing setstate user", req: newreq("GET", ts.URL+"/user/setstate/1GcsahF1mmbeX2y4uCgf96HISba/ready" , nil),expect: 404,state: vstate.Nothing},
		{name: "2: testing setstate user", req: newreq("GET", ts.URL+"/user/setstate/1GcsahF1mmbeX2y4uCgf96HISbb/Battery_low", nil),expect:404,state: vstate.Nothing},
		{name: "3: testing setstate user", req: newreq("GET", ts.URL+"/user/setstate/1GcsahF1mmbeX2y4uCgf96HISbc/Collected", nil),expect:404,state: vstate.Nothing},
		{name: "4: testing setstate user", req: newreq("GET", ts.URL+"/user/setstate/1GcsahF1mmbeX2y4uCgf96HISbd/Unknown", nil),expect:404,state: vstate.Nothing},
		{name: "5: testing setstate user", req: newreq("GET", ts.URL+"/user/setstate/1GcsahF1mmbeX2y4uCgf96HISbe/Bounty", nil),expect:404,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}

func Test006(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing setstate hunter", req: newreq("GET", ts.URL+"/hunter/setstate/1GcsahF1mmbeX2y4uCgf96HISba/ready" , nil),expect: 404,state: vstate.Ready},
		{name: "2: testing setstate hunter", req: newreq("GET", ts.URL+"/hunter/setstate/1GcsahF1mmbeX2y4uCgf96HISbb/Battery_low", nil),expect:404,state: vstate.Nothing},
		{name: "3: testing setstate hunter", req: newreq("GET", ts.URL+"/hunter/setstate/1GcsahF1mmbeX2y4uCgf96HISbc/Collected", nil),expect:404,state: vstate.Nothing},
		{name: "4: testing setstate hunter", req: newreq("GET", ts.URL+"/hunter/setstate/1GcsahF1mmbeX2y4uCgf96HISbd/Unknown", nil),expect:404,state: vstate.Nothing},
		{name: "5: testing setstate hunter", req: newreq("GET", ts.URL+"/hunter/setstate/1GcsahF1mmbeX2y4uCgf96HISbe/Bounty", nil),expect:404,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}


func Test007(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing claim user", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/Battery_low" , nil),expect: 200,state: vstate.Battery_low},
		{name: "2: testing claim hunter battery_low", req: newreq("GET", ts.URL+"/hunter/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:200,state: vstate.Bounty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
			var v Vehicle
			json.NewDecoder(resp.Body).Decode(&v)
			assert.Equal(t,v.State,tt.state,"Veachle should be in state ready")
		})
	}
}


func Test008(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing hunter bounty", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/Bounty" , nil),expect: 200,state: vstate.Bounty},
		{name: "2: testing hunt hunter collected", req: newreq("GET", ts.URL+"/hunter/hunt/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:200,state: vstate.Collected},
		{name: "3: testing hunt hunter Dropped", req: newreq("GET", ts.URL+"/hunter/hunt/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:200,state: vstate.Dropped},
		{name: "4: testing hunt hunter Ready", req: newreq("GET", ts.URL+"/hunter/hunt/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:200,state: vstate.Ready},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
			var v Vehicle
			json.NewDecoder(resp.Body).Decode(&v)
			assert.Equal(t,v.State,tt.state,"Veachle should be in state ready")
		})
	}
}

func Test009(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing user hunt", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/Bounty" , nil),expect: 200,state: vstate.Bounty},
		{name: "2: testing user hunter", req: newreq("GET", ts.URL+"/user/hunt/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:404,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}



func Test010(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing admin hunt", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/Bounty" , nil),expect: 200,state: vstate.Bounty},
		{name: "2: testing admin hunt", req: newreq("GET", ts.URL+"/admin/hunt/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:404,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}






func Test011(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing unknown claim", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/Unknown" , nil),expect: 200,state: vstate.Battery_low},
		{name: "2: testing claim user", req: newreq("GET", ts.URL+"/user/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
		{name: "3: testing claim hunter", req: newreq("GET", ts.URL+"/hunter/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
		{name: "4: testing claim admin", req: newreq("GET", ts.URL+"/admin/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}

func Test012(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing Service_mode claim", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/Service_mode" , nil),expect: 200,state: vstate.Service_mode},
		{name: "2: testing claim user", req: newreq("GET", ts.URL+"/user/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
		{name: "3: testing claim hunter", req: newreq("GET", ts.URL+"/hunter/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
		{name: "4: testing claim admin", req: newreq("GET", ts.URL+"/admin/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}

func Test013(t *testing.T) {
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		fmt.Printf("Requesting for URl: %s\n",url)
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		req    *http.Request
		expect int
		state vstate.State
	}{
		{name: "1: testing Terminated claim", req: newreq("GET", ts.URL+"/admin/setstate/1GcsahF1mmbeX2y4uCgf96HISba/Terminated" , nil),expect: 200,state: vstate.Terminated},
		{name: "2: testing claim user", req: newreq("GET", ts.URL+"/user/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
		{name: "3: testing claim hunter", req: newreq("GET", ts.URL+"/hunter/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
		{name: "4: testing claim admin", req: newreq("GET", ts.URL+"/admin/claim/1GcsahF1mmbeX2y4uCgf96HISba", nil),expect:403,state: vstate.Nothing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, resp.StatusCode, "OK response is expected")
			fmt.Print(resp.Body)
		})
	}
}
