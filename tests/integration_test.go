package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func getToken(t *testing.T, role string) string {
	body := []byte(fmt.Sprintf(`{"role":"%s"}`, role))
	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal("login failed:", err)
	}
	defer resp.Body.Close()

	var res map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&res)

	return res["token"]
}

func authPost(t *testing.T, token, path string, body any) *http.Response {
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+path, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s failed: %v", path, err)
	}
	return resp
}

func TestFullFlow(t *testing.T) {
	moderatorToken := getToken(t, "moderator")

	pvzReq := map[string]any{"city": "Казань"}
	resp := authPost(t, moderatorToken, "/pvz", pvzReq)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("create pvz failed: %s", string(body))
	}
	var pvz map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&pvz)
	pvzID := pvz["id"].(string)
	resp.Body.Close()

	employeeToken := getToken(t, "employee")

	receptionReq := map[string]any{"pvzId": pvzID}
	resp = authPost(t, employeeToken, "/receptions", receptionReq)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("create reception failed")
	}
	resp.Body.Close()

	for i := 0; i < 50; i++ {
		productReq := map[string]any{
			"type":  "электроника",
			"pvzId": pvzID,
		}
		resp = authPost(t, employeeToken, "/products", productReq)
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Fatalf("add product %d failed", i+1)
		}
		resp.Body.Close()
	}

	resp = authPost(t, employeeToken, "/pvz/"+pvzID+"/close_last_reception", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("close reception failed")
	}
	resp.Body.Close()
}
