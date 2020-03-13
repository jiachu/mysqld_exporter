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

const innodbClusterPrimaryQuery = `
	SELECT COUNT(*) FROM performance_schema.replication_group_members WHERE 
		member_id = (SELECT variable_value FROM performance_schema.global_status WHERE variable_name = 'group_replication_primary_member') AND MEMBER_STATE = 'ONLINE';
	`

// Metric descriptors.
var (
	InnodbClusterPrimaryDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, InnodbClusterSchema, "primary"),
		"display primary num", nil, nil,
	)
)

// ScrapeReplicationGroupMemberStats collects from `performance_schema.replication_group_member_stats`.
type ScrapeInnodbClusterPrimary struct{}

// Name of the Scraper. Should be unique.
func (ScrapeInnodbClusterPrimary) Name() string {
	return InnodbClusterSchema + ".primary"
}

// Help describes the role of the Scraper.
func (ScrapeInnodbClusterPrimary) Help() string {
	return "Collect metrics from performance_schema.replication_group_member and global variables"
}

// Version of MySQL from which scraper is available.
func (ScrapeInnodbClusterPrimary) Version() float64 {
	return 8.0
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeInnodbClusterPrimary) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	innodbClusterPrimaryRows, err := db.QueryContext(ctx, innodbClusterPrimaryQuery)
	if err != nil {
		return err
	}
	defer innodbClusterPrimaryRows.Close()

	var (
		primaryCount uint64
	)

	for innodbClusterPrimaryRows.Next() {
		if err := innodbClusterPrimaryRows.Scan(
			&primaryCount,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			InnodbClusterPrimaryDesc, prometheus.CounterValue, float64(primaryCount)
		),
	}
	return nil
}

// check interface
var _ Scraper = ScrapeInnodbClusterPrimary{}
