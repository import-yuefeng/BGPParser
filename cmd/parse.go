// MIT License

// Copyright (c) 2019 Yuefeng Zhu

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	client "github.com/import-yuefeng/BGPParser/cmd/client"
	"github.com/spf13/cobra"
	// log "github.com/sirupsen/logrus"
)

var (
	raw  []string
	data []string
)

func init() {
	parseCmd.Flags().StringArrayVarP(&raw, "rawFilePath", "r", []string{}, "bgp raw data path")
	parseCmd.Flags().StringArrayVarP(&data, "filePath", "p", []string{}, "bgp data path")

	rootCmd.AddCommand(parseCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parser bgpData",
	Long:  `parser bgpData`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewClient(":2048")
		if len(raw) > 0 {
			c.AddRawParse(raw)
		} else if len(data) > 0 {
			c.AddBGPParse(data)
		}
	},
}
