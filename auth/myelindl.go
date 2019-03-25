package auth

import (
    "os"
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/filebrowser/filebrowser/v2/settings"
    "github.com/filebrowser/filebrowser/v2/users"
)

const MethodMyelindlAuth settings.AuthMethod = "myelindl"

type tokenCred struct {
    Username string `json:"username"`
    Token    string `json:"token"`
}

type MyelindlAuth struct {
    AuthUrl string `json:"authurl"`
    RdirectUrl string `json:"redirecturl"`
}

func (a MyelindlAuth) Auth(r *http.Request, sto *users.Storage, root string) (*users.User, error) {
    var cred tokenCred
    var username string
    var token string

    fmt.Fprintf(os.Stderr, "Auth request recevied\n")
    if r.Body == nil {
        queryValues := r.URL.Query()
        username = queryValues.Get("username")
        token = queryValues.Get("token")
        fmt.Fprintf(os.Stderr, "Auth request url mode u:%s, t:%s\n", username, token)
    } else {
        err := json.NewDecoder(r.Body).Decode(&cred)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Auth request body mode fail json deco fail\n")
            return nil, os.ErrPermission
        }
        username = cred.Username
        token = cred.Token
        fmt.Fprintf(os.Stderr, "Auth request body mode u:%s, t:%s\n", cred.Username, cred.Token)
        fmt.Fprintf(os.Stderr, "Auth request body mode u:%s, t:%s\n", username, token)
    }
    if username == "" {
        fmt.Fprintf(os.Stderr, "Auth fail username empty\n")
        return nil, os.ErrPermission
    }
    if token == "" {
        fmt.Fprintf(os.Stderr, "Auth fail token empty\n")
        return nil, os.ErrPermission
    }
    // check token status
    url := fmt.Sprintf("%s/api/auth/%s?token=%s", a.AuthUrl, username, token)
    client := &http.Client{}
    resp, err := client.Get(url)
    if err != nil {
        fmt.Fprintf(os.Stderr, "myelindl web auth fail %s", url)
	    return nil, os.ErrPermission
	}
    if resp.StatusCode != http.StatusOK {
        fmt.Fprintf(os.Stderr, "myelindl web auth fail status code not ok\n")
		return nil, os.ErrPermission
	}

    var data struct {
		Status string `json:"status"`
	}

    err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
        fmt.Fprintf(os.Stderr, "myelindl web auth result json decode fail")
		return nil, os.ErrPermission
	}

	if data.Status != "ok" {
		return nil, os.ErrPermission
	}

	u, err := sto.Get(root, username)
    if err == nil {
		return  u, nil
	}
    user := &users.User{
        Username: username,
        Password: username,
        Scope: "/home/" + username,
		Locale: "en",
		LockPassword: false,
		Perm: users.Permissions{
			Admin:    false,
			Execute:  true,
			Create:   true,
			Rename:   true,
			Modify:   true,
			Delete:   true,
			Share:    true,
			Download: true,
		},
	}

    err = sto.Save(user)
    if err != nil {
        return nil, err
    }
    return user, err
}


func (a MyelindlAuth) LoginPage() bool {
    return false
}
