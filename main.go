package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	table.SetHeader([]string{"ID", "IMAGE", "STATUS", "PORTS", "NAMES"})
	table.SetBorder(false)
	table.SetTablePadding("\t")

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	var stoppedContainers []string

	for _, c := range containers {
		portStr := ""
		for _, p := range c.Ports {
			if p.PublicPort != 0 {
				portStr += fmt.Sprintf("%d:%d ", p.PublicPort, p.PrivatePort)
			}
		}

		displayStatus := c.Status
		if strings.HasPrefix(c.Status, "Up") {
			displayStatus = green(c.Status)
		} else if strings.HasPrefix(c.Status, "Exited") {
			displayStatus = red(c.Status)
			stoppedContainers = append(stoppedContainers, c.ID) // Запоминаем ID для удаления
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

	fmt.Println("\n🚀 Orbit Docker Inspector v1.2")
	fmt.Println("--------------------------------------------------")
	table.Render()

	// --- ЛОГИКА ОЧИСТКИ ---
	if len(stoppedContainers) > 0 {
		fmt.Printf("\n🧹 Найдено остановленных контейнеров: %d. Удалить их? (y/n): ", len(stoppedContainers))
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" {
			for _, id := range stoppedContainers {
				err := cli.ContainerRemove(context.Background(), id, container.RemoveOptions{})
				if err != nil {
					fmt.Printf("❌ Не удалось удалить %s: %v\n", id[:12], err)
				} else {
					fmt.Printf("✅ Удален: %s\n", id[:12])
				}
			}
			fmt.Println("✨ Очистка завершена!")
		} else {
			fmt.Println("Отмена очистки.")
		}
	} else {
		fmt.Println("\n✅ Все чисто! Остановленных контейнеров не найдено.")
	}
}