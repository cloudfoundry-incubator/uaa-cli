package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
)

func ListGroupValidations(cfg uaa.Config) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	return nil
}

func ListGroupsCmd(gm uaa.GroupManager, printer cli.Printer, filter, sortBy, sortOrder, attributes string, startIndex, count int) error {
	group, err := gm.List(filter, sortBy, attributes, uaa.ScimSortOrder(sortOrder), startIndex, count)
	if err != nil {
		return err
	}
	return printer.Print(group)
}

var listGroupsCmd = &cobra.Command{
	Use:     "list-groups",
	Aliases: []string{"groups", "get-groups", "search-groups"},
	Short:   "Search and list groups with SCIM filters",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(ListGroupValidations(GetSavedConfig()), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		gm := uaa.GroupManager{GetHttpClient(), cfg}
		err := ListGroupsCmd(gm, cli.NewJsonPrinter(log), filter, sortBy, sortOrder, attributes, startIndex, count)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(listGroupsCmd)
	listGroupsCmd.Annotations = make(map[string]string)
	listGroupsCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"

	listGroupsCmd.Flags().StringVarP(&filter, "filter", "", "", `a SCIM filter, or query, e.g. 'userName eq "bob@example.com"'`)
	listGroupsCmd.Flags().StringVarP(&sortBy, "sortBy", "b", "", `the attribute to sort results by, e.g. "created" or "userName"`)
	listGroupsCmd.Flags().StringVarP(&sortOrder, "sortOrder", "o", "", `how results should be ordered. One of: [ascending, descending]`)
	listGroupsCmd.Flags().StringVarP(&attributes, "attributes", "a", "", `include only these comma-separated user attributes to improve query performance`)
	listGroupsCmd.Flags().IntVarP(&startIndex, "startIndex", "s", 1, `starting index of paginated results`)
	listGroupsCmd.Flags().IntVarP(&count, "count", "c", 100, `maximum number of results to return`)
}
