package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	kubernetesclient "github.com/PaulWBrassard/devops-skills-assessment/go-rollout-restart/pkg/kubernetesclient"
)

func main() {
	search := flag.String("search", "", "The substring to search for in Deployment names")
	flag.Parse()

	if *search == "" {
		fmt.Println("Error: The 'search' flag is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	kubernetesclient, err := kubernetesclient.NewKubernetesClient()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Searching for deployments with %s in their name...\n", *search)
	deployments, err := kubernetesclient.ListDeployments(*search)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Found %d deployment(s) matching the search criteria\n", len(deployments))
	if len(deployments) > 0 {
		for _, deployment := range deployments {
			fmt.Printf("- %s\n", deployment.Name)
		}
		fmt.Print("Do you want to rollout restart these deployments? (y/n): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "y" {
			for _, deployment := range deployments {
				err = kubernetesclient.RolloutRestartDeployment(deployment)
				if err != nil {
					fmt.Printf("Error rolling out restart for deployment %s: %v\n", deployment.Name, err)
				} else {
					fmt.Printf("Rollout restart initiated for deployment %s\n", deployment.Name)
				}
			}
		} else {
			fmt.Println("Operation cancelled.")
		}
	}
}
