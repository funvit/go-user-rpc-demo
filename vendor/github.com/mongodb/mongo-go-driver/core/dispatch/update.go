// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package dispatch

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/options"

	"github.com/mongodb/mongo-go-driver/core/command"
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/result"
	"github.com/mongodb/mongo-go-driver/core/session"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/core/uuid"
	"github.com/mongodb/mongo-go-driver/core/writeconcern"
)

// Update handles the full cycle dispatch and execution of an update command against the provided
// topology.
func Update(
	ctx context.Context,
	cmd command.Update,
	topo *topology.Topology,
	selector description.ServerSelector,
	clientID uuid.UUID,
	pool *session.Pool,
	retryWrite bool,
	opts ...*options.UpdateOptions,
) (result.Update, error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return result.Update{}, err
	}

	// If no explicit session and deployment supports sessions, start implicit session.
	if cmd.Session == nil && topo.SupportsSessions() {
		cmd.Session, err = session.NewClientSession(pool, clientID, session.Implicit)
		if err != nil {
			return result.Update{}, err
		}
		defer cmd.Session.EndSession()
	}

	updateOpts := options.MergeUpdateOptions(opts...)

	if updateOpts.ArrayFilters != nil {
		if ss.Description().WireVersion.Max < 6 {
			return result.Update{}, ErrArrayFilters
		}
		arr, err := updateOpts.ArrayFilters.ToArray()
		if err != nil {
			return result.Update{}, err
		}
		cmd.Opts = append(cmd.Opts, bson.EC.Array("arrayFilters", arr))
	}
	if updateOpts.BypassDocumentValidation != nil && ss.Description().WireVersion.Includes(4) {
		cmd.Opts = append(cmd.Opts, bson.EC.Boolean("bypassDocumentValidation", *updateOpts.BypassDocumentValidation))
	}
	if updateOpts.Collation != nil {
		if ss.Description().WireVersion.Max < 5 {
			return result.Update{}, ErrCollation
		}
		cmd.Opts = append(cmd.Opts, bson.EC.SubDocument("collation", updateOpts.Collation.ToDocument()))
	}
	if updateOpts.Upsert != nil {
		cmd.Opts = append(cmd.Opts, bson.EC.Boolean("upsert", *updateOpts.Upsert))
	}

	// Execute in a single trip if retry writes not supported, or retry not enabled
	if !retrySupported(topo, ss.Description(), cmd.Session, cmd.WriteConcern) || !retryWrite {
		if cmd.Session != nil {
			cmd.Session.RetryWrite = false // explicitly set to false to prevent encoding transaction number
		}
		return update(ctx, cmd, ss, nil)
	}

	cmd.Session.RetryWrite = retryWrite
	cmd.Session.IncrementTxnNumber()

	res, originalErr := update(ctx, cmd, ss, nil)

	// Retry if appropriate
	if cerr, ok := originalErr.(command.Error); ok && cerr.Retryable() ||
		res.WriteConcernError != nil && command.IsWriteConcernErrorRetryable(res.WriteConcernError) {
		ss, err := topo.SelectServer(ctx, selector)

		// Return original error if server selection fails or new server does not support retryable writes
		if err != nil || !retrySupported(topo, ss.Description(), cmd.Session, cmd.WriteConcern) {
			return res, originalErr
		}

		return update(ctx, cmd, ss, cerr)
	}
	return res, originalErr

}

func update(
	ctx context.Context,
	cmd command.Update,
	ss *topology.SelectedServer,
	oldErr error,
) (result.Update, error) {
	desc := ss.Description()

	conn, err := ss.Connection(ctx)
	if err != nil {
		if oldErr != nil {
			return result.Update{}, oldErr
		}
		return result.Update{}, err
	}

	if !writeconcern.AckWrite(cmd.WriteConcern) {
		go func() {
			defer func() { _ = recover() }()
			defer conn.Close()

			_, _ = cmd.RoundTrip(ctx, desc, conn)
		}()

		return result.Update{}, command.ErrUnacknowledgedWrite
	}
	defer conn.Close()

	return cmd.RoundTrip(ctx, desc, conn)
}