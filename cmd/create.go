package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/daichirata/hammer/internal/hammer"
)

var (
	createExample = `
* Create database and apply local schema (faster than running database creation and schema apply separately)
  hammer create spanner://projects/projectId/instances/instanceId/databases/databaseName /path/to/file

* Copy database
  hammer create spanner://projects/projectId/instances/instanceId/databases/databaseName1 spanner://projects/projectId/instances/instanceId/databases/databaseName2`

	createCmd = &cobra.Command{
		Use:     "create DATABASE SOURCE",
		Short:   "Create database and apply schema",
		Example: createExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("must specify 2 argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			databaseURI := args[0]
			sourceURI := args[1]

			if hammer.Scheme(databaseURI) != "spanner" {
				return fmt.Errorf("DATABASE must be a spanner URI")
			}
			database, err := hammer.NewSpannerSource(ctx, databaseURI)
			if err != nil {
				return err
			}
			source, err := hammer.NewSource(ctx, sourceURI)
			if err != nil {
				return err
			}

			ddl, err := source.DDL(ctx)
			if err != nil {
				return err
			}
			return database.Create(ctx, ddl)
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
}
