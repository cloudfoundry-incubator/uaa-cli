package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "code.cloudfoundry.org/uaa-cli/fixtures"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("ListGroups", func() {

	var groupListResponse string

	BeforeEach(func() {
		cfg := uaa.NewConfigWithServerURL(server.URL())
		cfg.AddContext(uaa.NewContextWithToken("access_token"))
		config.WriteConfig(cfg)
		groupListResponse = fmt.Sprintf(PaginatedResponseTmpl, AdminGroupResponse, ReadGroupResponse)
	})

	It("executes SCIM queries based on flags", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=verified+eq+false&attributes=id%2CdisplayName&sortBy=displayName&sortOrder=descending&count=50&startIndex=100"),
			RespondWith(http.StatusOK, groupListResponse),
		))

		session := runCommand("list-groups",
			"--filter", "verified eq false",
			"--attributes", "id,displayName",
			"--sortBy", "displayName",
			"--sortOrder", "descending",
			"--count", "50",
			"--startIndex", "100")

		Eventually(session).Should(Exit(0))
	})
})
