package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/cloudwatch"
	"github.com/K-Phoen/grabana/timeseries"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: go run main.go http://grafana-host:3000 api-key-string-here\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[1], grabana.WithAPIToken(os.Args[2]))

	// create the folder holding the dashboard for the service
	folder, err := client.FindOrCreateFolder(ctx, "Test Folder")
	if err != nil {
		fmt.Printf("Could not find or create folder: %s\n", err)
		os.Exit(1)
	}

	builder, err := dashboard.New(
		"Awesome dashboard test",
		dashboard.UID("test-dashboard-alerts"),
		dashboard.AutoRefresh("30s"),
		dashboard.Time("now-30m", "now"),
		dashboard.Tags([]string{"generated"}),
		dashboard.Row(
			"Cloudwatch",
			row.WithTimeSeries(
				"test",
				timeseries.DataSource("AWS CloudWatch"),
				timeseries.WithCloudwatchTarget(
					cloudwatch.QueryParams{
						Dimensions: map[string]string{
							"DBClusterIdentifier": "RDS-Cluster-Name",
						},
						MetricName: "CPUUtilization",
						Statistics: []string{"Maximum"},
						Namespace:  "AWS/RDS",
						Period:     "30",
						Region:     "default",
					},
				),
			),
		),
	)
	if err != nil {
		fmt.Printf("Could not build dashboard: %s\n", err)
		os.Exit(1)
	}

	dash, err := client.UpsertDashboard(ctx, folder, builder)

	if err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("The deed is done:\n%s\n", os.Args[1]+dash.URL)
}

func float64Ptr(input float64) *float64 {
	return &input
}
