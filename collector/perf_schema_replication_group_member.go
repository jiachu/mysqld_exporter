// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const perfReplicationGroupMemeberQuery = `
	SELECT MEMBER_ID, MEMBER_HOST, MEMBER_PORT, MEMBER_STATE, MEMBER_ROLE
	  FROM performance_schema.replication_group_members
	`

// Metric descriptors.
var (
	performanceSchemaReplicationGroupMemberDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "members"),
		"mysql replication group members",
		[]string{"member_id", "member_host", "member_port", "member_state", "member_role"}, nil,
	)
)

// ScrapeReplicationGroupMemberStats collects from `performance_schema.replication_group_member_stats`.
type ScrapePerfReplicationGroupMember struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfReplicationGroupMember) Name() string {
	return performanceSchema + ".replication_group_member"
}

// Help describes the role of the Scraper.
func (ScrapePerfReplicationGroupMember) Help() string {
	return "Collect metrics from performance_schema.replication_group_member"
}

// Version of MySQL from which scraper is available.
func (ScrapePerfReplicationGroupMember) Version() float64 {
	return 8.0
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfReplicationGroupMember) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	perfReplicationGroupMemeberRows, err := db.QueryContext(ctx, perfReplicationGroupMemeberQuery)
	if err != nil {
		return err
	}
	defer perfReplicationGroupMemeberRows.Close()

	var (
		memberId, memberHost, memberPort, memberState, memberRole string
	)

	for perfReplicationGroupMemeberRows.Next() {
		if err := perfReplicationGroupMemeberRows.Scan(
			&memberId, &memberHost, &memberPort,
			&memberState, &memberRole,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaReplicationGroupMemberDesc, prometheus.GaugeValue, 1 ,
			memberId, memberHost, memberPort, memberState, memberRole,
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapePerfReplicationGroupMember{}
