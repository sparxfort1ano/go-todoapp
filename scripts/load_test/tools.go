package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func seedUsers(base string, n int) ([]int, error) {
	ids := make([]int, n)
	for i := range n {
		resp, body, err := doJSON(http.MethodPost, base+"/users", CreateUserRequest{
			FullName:    fmt.Sprintf("LT-User-%d-%s", i+1, randomString(4)),
			PhoneNumber: randomPhone(),
		})
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("seed POST /users → %d: %s", resp.StatusCode, body)
		}
		var u UserResponse
		json.Unmarshal(body, &u)
		ids[i] = u.ID
	}

	return ids, nil
}

func seedTasks(base string, userIDs []int, perUser int) (map[int][]int, error) {
	m := make(map[int][]int, len(userIDs))
	for _, uid := range userIDs {
		for j := range perUser {
			resp, body, err := doJSON(http.MethodPost, base+"/tasks", CreateTaskRequest{
				Title:        fmt.Sprintf("Seed task %d-%d", uid, j),
				Description:  fmt.Sprintf("Seeded task for user %d", uid),
				AuthorUserID: uid,
			})
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusCreated {
				return nil, fmt.Errorf("seed POST /tasks → %d: %s", resp.StatusCode, body)
			}
			var t TaskResponse
			json.Unmarshal(body, &t)
			m[uid] = append(m[uid], t.ID)
		}
	}

	return m, nil
}

func warmUpCache(base string, userIDs []int) {
	for _, uid := range userIDs {
		doJSON(http.MethodGet, fmt.Sprintf("%s/tasks?user_id=%d&limit=10", base, uid), nil)
	}
}

func checkServer(base string) error {
	fmt.Print("	server check... ")
	_, _, err := doJSON(http.MethodGet, base+"/tasks?limit=1", nil)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return err
	}

	fmt.Println("OK")
	return nil
}
