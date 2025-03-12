package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/types"
)

// GetQueryCmd는 checkpoint 모듈의 쿼리 명령어를 반환합니다.
func GetQueryCmd() *cobra.Command {
	checkpointQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s 쿼리 명령어", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	checkpointQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryCheckpoint(),
		GetCmdQueryLatestCheckpoint(),
		GetCmdQueryCheckpoints(),
		GetCmdQueryCheckpointCount(),
	)

	return checkpointQueryCmd
}

// GetCmdQueryParams는 모듈 파라미터를 조회하는 명령어를 반환합니다.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "모듈 파라미터를 조회합니다",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCheckpoint는 특정 번호의 체크포인트를 조회하는 명령어를 반환합니다.
func GetCmdQueryCheckpoint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint [number]",
		Short: "특정 번호의 체크포인트를 조회합니다",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			number, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("체크포인트 번호를 파싱할 수 없습니다: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Checkpoint(cmd.Context(), &types.QueryCheckpointRequest{Number: number})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Checkpoint)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryLatestCheckpoint는 최신 체크포인트를 조회하는 명령어를 반환합니다.
func GetCmdQueryLatestCheckpoint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-checkpoint",
		Short: "최신 체크포인트를 조회합니다",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.LatestCheckpoint(cmd.Context(), &types.QueryLatestCheckpointRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Checkpoint)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCheckpoints는 모든 체크포인트를 조회하는 명령어를 반환합니다.
func GetCmdQueryCheckpoints() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoints",
		Short: "모든 체크포인트를 조회합니다",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Checkpoints(cmd.Context(), &types.QueryCheckpointsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "checkpoints")

	return cmd
}

// GetCmdQueryCheckpointCount는 체크포인트 수를 조회하는 명령어를 반환합니다.
func GetCmdQueryCheckpointCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-count",
		Short: "체크포인트 수를 조회합니다",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CheckpointCount(cmd.Context(), &types.QueryCheckpointCountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
