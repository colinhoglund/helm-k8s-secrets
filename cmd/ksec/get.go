package main

import (
	"encoding/json"
	"fmt"

	"github.com/colinhoglund/ksec/pkg/models"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [secret]",
	Short: "Get values from a Secret",
	Args:  cobra.ExactArgs(1),
	RunE:  getCommand,
}

func getCommand(cmd *cobra.Command, args []string) error {
	secret, err := secretsClient.Get(args[0])
	if err != nil {
		return err
	}

	var lines []string
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if verbose {
		for key, value := range secret.Data {
			rawAnnotation := secret.Annotations[fmt.Sprintf("ksec.io/%s", key)]

			var jsonAnnotation []byte
			if rawAnnotation == "" {
				jsonAnnotation = []byte(`{"updatedBy": "", "lastUpdated": ""}`)
			} else {
				jsonAnnotation = []byte(rawAnnotation)
			}

			annotation := models.KeyAnnotation{}
			if err := json.Unmarshal(jsonAnnotation, &annotation); err != nil {
				return err
			}
			lines = append(lines, fmt.Sprintf("Key:\t%s", key))
			lines = append(lines, fmt.Sprintf("Value:\t%s", value))
			lines = append(lines, fmt.Sprintf("User:\t%s", annotation.UpdatedBy))
			lines = append(lines, fmt.Sprintf("Updated:\t%s\n", annotation.LastUpdated))
		}
	} else {
		lines = append(lines, "KEY\tVALUE")
		for key, value := range secret.Data {
			lines = append(lines, fmt.Sprintf("%s\t%s", key, value))
		}
	}
	outputTabular(lines)
	return nil
}
