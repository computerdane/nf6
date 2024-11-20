package server_secure

import (
	"context"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerSecure) AddMachine(ctx context.Context, in *nf6.AddMachineRequest) (*nf6.AddMachineReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if err := lib.ValidateHostName(in.GetHostName()); err != nil {
		return nil, err
	}
	if err := lib.ValidateWireguardKey(in.GetWgPublicKey()); err != nil {
		return nil, err
	}
	if err := lib.ValidateIpv6Address(in.GetAddrIpv6()); err != nil {
		return nil, err
	}

	query := "select count(*) from machine where account_id = @account_id and host_name = @host_name"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"host_name":  in.GetHostName(),
	}
	machineExists := 0
	if err := s.db.QueryRow(ctx, query, args).Scan(&machineExists); err != nil {
		return nil, err
	}
	if machineExists != 0 {
		return nil, status.Error(codes.AlreadyExists, "machine already exists")
	}

	args = pgx.NamedArgs{
		"account_id":    accountId,
		"host_name":     in.GetHostName(),
		"wg_public_key": in.GetWgPublicKey(),
		"addr_ipv6":     in.GetAddrIpv6(),
	}
	query = "insert into machine (account_id, host_name, wg_public_key, addr_ipv6) values (@account_id, @host_name, @wg_public_key, @addr_ipv6)"
	if _, err := s.db.Exec(ctx, query, args); err != nil {
		return nil, err
	}

	return &nf6.AddMachineReply{Success: true}, nil
}

func (s *ServerSecure) ListMachines(ctx context.Context, in *nf6.ListMachinesRequest) (*nf6.ListMachinesReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, "select host_name from machine where account_id = $1", accountId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no machines found")
	}

	reply := &nf6.ListMachinesReply{Names: []string{}}

	for rows.Next() {
		var machineName = ""
		err := rows.Scan(&machineName)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal server error")
		}
		reply.Names = append(reply.Names, machineName)
	}

	return reply, nil
}

func (s *ServerSecure) RenameMachine(ctx context.Context, in *nf6.RenameMachineRequest) (*nf6.RenameMachineReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if err := lib.ValidateHostName(in.GetNewName()); err != nil {
		return nil, err
	}

	query := "select id from machine where account_id = @account_id and host_name = @host_name"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"host_name":  in.GetOldName(),
	}
	machineId := 0
	err = s.db.QueryRow(ctx, query, args).Scan(&machineId)
	if machineId == 0 || err != nil {
		return nil, status.Error(codes.NotFound, "machine not found")
	}

	query = "update machine set host_name = @host_name where id = @id"
	args = pgx.NamedArgs{
		"host_name": in.GetNewName(),
		"id":        machineId,
	}
	if _, err := s.db.Exec(ctx, query, args); err != nil {
		return nil, err
	}

	return &nf6.RenameMachineReply{Success: true}, nil
}
