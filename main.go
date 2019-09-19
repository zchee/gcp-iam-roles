// Copyright 2019 The gcp-iam-roles Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/olekukonko/tablewriter"
)

func main() {
	c := colly.NewCollector(colly.Async(true))

	var roles [][]string

	c.OnHTML("div#predefined-roles", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, e *colly.HTMLElement) {
			e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
				role := make([]string, 4)
				e.ForEach("td.role-name code", func(_ int, e *colly.HTMLElement) {
					role[0] = e.Text
				})
				title := e.ChildText("td.role-title")
				if idx := strings.Index(title, "Beta"); idx > -1 {
					title = title[:idx] + "\x1b[31m" + title[idx:idx+4] + "\x1b[0m" + title[idx+4:]
				}
				if idx := strings.Index(title, "Alpha"); idx > -1 {
					title = title[:idx] + "\x1b[31;1;4m" + title[idx:idx+5] + "\x1b[0m" + title[idx+5:]
				}
				role[1] = strings.TrimSuffix(title, "\n")
				role[2] = e.ChildText("td.role-description")
				role[3] = e.ChildText("td.role-permissions")

				if role[0] != "" {
					roles = append(roles, role)
				}
			})
		})
	})

	c.Visit("https://cloud.google.com/iam/docs/understanding-roles")
	c.Wait()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetHeader([]string{"Role", "Title", "Description", "Permissions"})

	for _, v := range roles {
		table.Append(v)
	}

	table.Render()
}
