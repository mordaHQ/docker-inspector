package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color" // Мы скачали это раньше через зависимости tablewriter
	"github.com/olekukonko/tablewriter"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("❌ Ошибка: %v", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Fatalf("❌ Ошибка Docker: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "IMAGE", "STATUS", "PORTS", "NAMES"})
	table.SetBorder(false)
	table.SetTablePadding("\t")

	// Настраиваем цвета
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	for _, c := range containers {
		// Собираем порты
		portStr := ""
		for _, p := range c.Ports {
			if p.PublicPort != 0 {
				portStr += fmt.Sprintf("%d:%d ", p.PublicPort, p.PrivatePort)
			}
		}

		// Раскрашиваем статус
		displayStatus := c.Status
		if strings.HasPrefix(c.Status, "Up") {
			displayStatus = green(c.Status)
		} else if strings.HasPrefix(c.Status, "Exited") {
			displayStatus = red(c.Status)
		} else {
			displayStatus = yellow(c.Status)
		}

		table.Append([]string{
			c.ID[:12],
			c.Image,
			displayStatus,
			portStr,
			fmt.Sprintf("%v", c.Names),
		})
	}

	fmt.Println("\n🚀 Orbit Docker Inspector v1.1")
	fmt.Println("--------------------------------------------------")
	table.Render()
}