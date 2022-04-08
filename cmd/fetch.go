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
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/silabs-EricB/svdgrab/packIndex"
	"github.com/cavaliercoder/grab"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Svdindex struct {
	XMLName     xml.Name `xml:"item"`
	Game        string   `xml:"game,attr"`
	NextJp      string   `xml:"nextjp,attr"`
	NextDd      string   `xml:"nextdd,attr"`
	WinNums     string   `xml:"winnum,attr"`
	WinDd       string   `xml:"windd,attr"`
	Myflv       string   `xml:"myflv,attr"`
	WinnumNM    string   `xml:"winnumNM,attr"`
	Name        string   `xml:"name,attr"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubdate"`
	Guid        string   `xml:"guid"`
}

// tweaked from: https://stackoverflow.com/a/42718113/1170664
func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read body: %v", err)
	}

	return data, nil
}

func extractPack(archiveFilePath string, pathToSave string) {
	dst := pathToSave
	archive, err := zip.OpenReader(archiveFilePath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	_, file := filepath.Split(archiveFilePath)
	fmt.Println(color.GreenString("\t🛠️  unzipping "), file)
	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		if !strings.Contains(filePath, "SVD") {
			continue
		}

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println(color.RedString("invalid file path"))
			return
		}
		if f.FileInfo().IsDir() {
			// fmt.Println(color.BlueString("creating directory..."))
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		_, file := filepath.Split(filePath)
		fmt.Println(color.GreenString("\t\t⛏️  extracting "), file)

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Get all the latest SVD available from vendors",
	Long: `Retrieve the pack and extract the SVD.`,

	Run: func(cmd *cobra.Command, args []string) {
		vendor, _ := cmd.Flags().GetString("vendor")
		pathToSave, _ := cmd.Flags().GetString("path")
		count := 0

		err := os.MkdirAll(pathToSave, os.ModePerm)
		if err != nil {
			log.Printf("Failed to create the directory. %v\n", err)
			os.Exit(-1)
		}

		if xmlBytes, err := getXML("https://www.keil.com/pack/index.pidx"); err != nil {
			log.Printf("Failed to get XML: %v", err)
		} else {
			result := packIndex.Index{}
			xml.Unmarshal(xmlBytes, &result)

			fmt.Println("Vendor ", result.Vendor)
			for _, s := range result.Pindex {
				for _, t := range s.Pdsc {
					//fmt.Println(i, j, t.VendorAttr)
					if strings.EqualFold(strings.ToLower(t.VendorAttr), strings.ToLower(vendor)) && (t.DeprecatedAttr == "") {
						// https://www.silabs.com/documents/public/cmsis-packs/GeckoPlatform_EFR32MG24_DFP.3.2.0.pack
						filename := t.UrlAttr + t.VendorAttr + "." + t.NameAttr + "." + t.VersionAttr + ".pack"
						client := grab.NewClient()
						req, _ := grab.NewRequest(pathToSave, filename)

						// start download
						fmt.Printf(color.GreenString("🚚 Downloading %v...\n"), req.URL())
						resp := client.Do(req)
						if err := resp.Err(); err != nil {
							fmt.Printf(color.RedString("\t📜 HTTP status codes: %v\n"), resp.HTTPResponse.Status)
						} else {
							fmt.Printf(color.BlueString("\t📜 HTTP status codes: %v\n"), resp.HTTPResponse.Status)
						}

						// start UI loop
						t := time.NewTicker(500 * time.Millisecond)
						defer t.Stop()

					Loop:
						for {
							select {
							case <-t.C:
								fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
									resp.BytesComplete(),
									resp.Size,
									100*resp.Progress())

							case <-resp.Done:
								// download is complete
								break Loop
							}
						}

						// check for errors
						if err := resp.Err(); err != nil {
							//fmt.Fprintf(os.Stderr, color.RedString("Download failed: %v\n"), err)
							//os.Exit(1)
							continue
						} else {
							fmt.Printf(color.GreenString("\t📦 saved to ./%v \n"), resp.Filename)
							extractPack(resp.Filename, pathToSave)
							fmt.Printf("\n")
							count++
						}
					}
				}
			}
		}

		fmt.Printf("✨ Done, %d packs from vendor %s stored in folder [%s].\n", count, color.CyanString(vendor), color.CyanString(filepath.Join(pathToSave, "SVD")))
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
