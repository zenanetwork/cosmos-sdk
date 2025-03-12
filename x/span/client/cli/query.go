package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/span/types"
)

// GetQueryCmd는 span 모듈의 쿼리 명령어를 반환합니다.
func GetQueryCmd() *cobra.Command {
	spanQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s 쿼리 명령어", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	spanQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQuerySpan(),
		GetCmdQueryCurrentSpan(),
		GetCmdQuerySpans(),
		GetCmdQuerySpanByHeight(),
	)

	return spanQueryCmd
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

// GetCmdQuerySpan은 특정 ID의 스팬을 조회하는 명령어를 반환합니다.
func GetCmdQuerySpan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "span [span-id]",
		Short: "특정 ID의 스팬을 조회합니다",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			spanID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("스팬 ID를 파싱할 수 없습니다: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Span(cmd.Context(), &types.QuerySpanRequest{SpanId: spanID})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Span)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCurrentSpan은 현재 활성화된 스팬을 조회하는 명령어를 반환합니다.
func GetCmdQueryCurrentSpan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-span",
		Short: "현재 활성화된 스팬을 조회합니다",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CurrentSpan(cmd.Context(), &types.QueryCurrentSpanRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Span)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQuerySpans은 모든 스팬을 조회하는 명령어를 반환합니다.
func GetCmdQuerySpans() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spans",
		Short: "모든 스팬을 조회합니다",
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
			res, err := queryClient.Spans(cmd.Context(), &types.QuerySpansRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "spans")

	return cmd
}

// GetCmdQuerySpanByHeight은 특정 블록 높이에 해당하는 스팬을 조회하는 명령어를 반환합니다.
func GetCmdQuerySpanByHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "span-by-height [height]",
		Short: "특정 블록 높이에 해당하는 스팬을 조회합니다",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			height, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("블록 높이를 파싱할 수 없습니다: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SpanByHeight(cmd.Context(), &types.QuerySpanByHeightRequest{Height: height})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Span)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
