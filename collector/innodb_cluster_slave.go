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

const innodbClusterSlaveQuery = `
	SELECT COUNT(*) FROM performance_schema.replication_group_members WHERE 
		MEMBER_ROLE = "SECONDARY" AND MEMBER_STATE = 'ONLINE';
	`

// Metric descriptors.
var (
	InnodbClusterSlaveDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, InnodbClusterSchema, "slave"),
		"display slave num", nil, nil,
	)
)

// ScrapeReplicationGroupMemberStats collects from `performance_schema.replication_group_member_stats`.
type ScrapeInnodbClusterSlave struct{}

// Name of the Scraper. Should be unique.
func (ScrapeInnodbClusterSlave) Name() string {
	return InnodbClusterSchema + ".slave"
}

// Help describes the role of the Scraper.
func (ScrapeInnodbClusterSlave) Help() string {
	return "Collect metrics from performance_schema.replication_group_members"
}

// Version of MySQL from which scraper is available.
func (ScrapeInnodbClusterSlave) Version() float64 {
	return 8.0
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeInnodbClusterSlave) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	innodbClusterSlaveRows, err := db.QueryContext(ctx, innodbClusterSlaveQuery)
	if err != nil {
		return err
	}
	defer innodbClusterSlaveRows.Close()

	var (
		slaveCount uint64
	)

	for innodbClusterSlaveRows.Next() {
		if err := innodbClusterSlaveRows.Scan(
			&slaveCount,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			InnodbClusterSlaveDesc, prometheus.GaugeValue, float64(slaveCount),
		)
	}
	return nil
}

// check interface
var _ Scraper = ScrapeInnodbClusterSlave{}
