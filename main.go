package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
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
	table.SetHeader([]string{"ID", "IMAGE", "STATUS", "IP ADDRESS", "RAM USAGE", "NAMES"})
	table.SetBorder(false)
	table.SetTablePadding("\t")

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	var stoppedContainers []string

	for _, c := range containers {
		// 1. Получаем IP
		ip := "-"
		for _, net := range c.NetworkSettings.Networks {
			if net.IPAddress != "" {
				ip = net.IPAddress
			}
		}

		// 2. Получаем RAM (только для запущенных)
		ramUsage := "-"
		if strings.HasPrefix(c.Status, "Up") {
			stats, err := cli.ContainerStats(context.Background(), c.ID, false)
			if err == nil {
				var v types.StatsJSON
				json.NewDecoder(stats.Body).Decode(&v)
				stats.Body.Close()
				// Переводим байты в Мегабайты
				ramUsage = fmt.Sprintf("%.2f MB", float64(v.MemoryStats.Usage)/1024/1024)
			}
		}

		displayStatus := c.Status
		if strings.HasPrefix(c.Status, "Up") {
			displayStatus = green(c.Status)
		} else {
			displayStatus = red(c.Status)
			stoppedContainers = append(stoppedContainers, c.ID)
		}

		table.Append([]string{
			c.ID[:12],
			c.Image,
			displayStatus,
			ip,
			ramUsage,
			fmt.Sprintf("%v", c.Names),
		})
	}

	fmt.Println("\n🚀 Orbit Docker Inspector v1.3 | Professional Edition")
	table.Render()

	if len(stoppedContainers) > 0 {
		fmt.Printf("\n🧹 Найдено остановленных контейнеров: %d. Удалить их? (y/n): ", len(stoppedContainers))
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) == "y" {
			for _, id := range stoppedContainers {
				cli.ContainerRemove(context.Background(), id, container.RemoveOptions{})
				fmt.Printf("✅ Удален: %s\n", id[:12])
			}
		}
	}
}