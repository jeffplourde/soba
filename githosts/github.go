package githosts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	gitHubCallSize = 100
)

type githubHost struct {
	Provider string
	APIURL   string
}

type Edge struct {
	Node struct {
		Name          string
		NameWithOwner string
		URL           string `json:"Url"`
		SSHURL        string `json:"sshUrl"`
	}
	Cursor string
}

type githubQueryNamesResponse struct {
	Data struct {
		Viewer struct {
			Repositories struct {
				Edges    []Edge
				PageInfo struct {
					EndCursor   string
					HasNextPage bool
				}
			}
		}
	}
}

func (provider githubHost) describeRepos() describeReposOutput {
	logger.Println("listing GitHub repositories")
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	var repos []repository
	reqBody := "{\"query\": \"query { viewer { repositories(first:" + strconv.Itoa(gitHubCallSize) + ") { edges { node { name nameWithOwner url sshUrl } cursor } pageInfo { endCursor hasNextPage }} } }\""
	for {
		mJSON := reqBody
		contentReader := bytes.NewReader([]byte(mJSON))
		req, _ := http.NewRequest("POST", "https://api.github.com/graphql", contentReader)
		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", os.Getenv("GITHUB_TOKEN")))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Accept", "application/json; charset=utf-8")

		resp, _ := client.Do(req)
		bodyB, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(bytes.Replace(bodyB, []byte("\r"), []byte("\r\n"), -1))
		var respObj githubQueryNamesResponse
		if err := json.Unmarshal([]byte(bodyStr), &respObj); err != nil {
			logger.Fatal(err)
		}

		for _, repo := range respObj.Data.Viewer.Repositories.Edges {
			repos = append(repos, repository{
				Name:          repo.Node.Name,
				SSHUrl:        repo.Node.SSHURL,
				HTTPSUrl:      repo.Node.URL,
				NameWithOwner: repo.Node.NameWithOwner,
				Domain:        "github.com",
			})
		}
		if !respObj.Data.Viewer.Repositories.PageInfo.HasNextPage {
			break
		} else {
			reqBody = "{\"query\": \"query($first:Int $after:String){ viewer { repositories(first:$first after:$after) { edges { node { name nameWithOwner url sshUrl } cursor } pageInfo { endCursor hasNextPage }} } }\", \"variables\":{\"first\":" + strconv.Itoa(gitHubCallSize) + ",\"after\":\"" + respObj.Data.Viewer.Repositories.PageInfo.EndCursor + "\"} }"
		}
	}

	return describeReposOutput{
		Repos: repos,
	}
}

func (provider githubHost) getAPIURL() string {
	return provider.APIURL
}

func (provider githubHost) Backup(backupDIR string) {
	output := provider.describeRepos()
	for _, repo := range output.Repos {
		firstPos := strings.Index(repo.HTTPSUrl, "//")
		repo.URLWithToken = repo.HTTPSUrl[:firstPos+2] + os.Getenv("GITHUB_TOKEN") + "@" + repo.HTTPSUrl[firstPos+2:]
		processBackup(repo, backupDIR)
	}
}