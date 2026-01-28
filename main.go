package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/olekukonko/tablewriter"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Создаем таблицу
	table := tablewriter.NewWriter(os.Stdout)
	
	// Если SetHeader не работает, мы просто выведем данные строками
	fmt.Println("\n🚀 Orbit Docker Inspector v1.0")
	fmt.Println("ID\t\tIMAGE\t\tSTATUS")
	fmt.Println("--------------------------------------------------")

	for _, c := range containers {
		table.Append([]string{
			c.ID[:12],
			c.Image,
			c.Status,
		})
	}

	table.Render()
}