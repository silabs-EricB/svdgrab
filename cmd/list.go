/*
Copyright © 2021 Eric Bissonnette <eric.bissonnette@silabs.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bissonex/svdgrab/packIndex"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func unique(sample []string) []string {
	var unique []string
sampleLoop:
	for _, v := range sample {
		for i, u := range unique {
			if v == u {
				unique[i] = v
				continue sampleLoop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve a listing of vendors",
	Long: `Use this option to find the name of vendor to use with the fetch command`,
	Run: func(cmd *cobra.Command, args []string) {
		vendors := []string{}
		if xmlBytes, err := getXML("https://www.keil.com/pack/index.pidx"); err != nil {
			log.Printf("Failed to get XML: %v", err)
		} else {
			result := packIndex.Index{}
			xml.Unmarshal(xmlBytes, &result)

			fmt.Println("Vendor ", result.Vendor)
			for _, s := range result.Pindex {
				for _, t := range s.Pdsc {
					//fmt.Printf(color.GreenString("🏬 Vendor %s\n"), t.VendorAttr)
					vendors = append(vendors, strings.ToLower(t.VendorAttr))
				}
			}
		}

		sort.Strings(vendors)
		vendors = unique(vendors)
		for _, vendor := range vendors {
			fmt.Printf(color.GreenString("🏬 %s\n"), vendor)
		}
		fmt.Println("list called")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
