package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/cosmos/cosmos-sdk/x/span/types"
)

// GetTxCmd는 span 모듈의 트랜잭션 명령어를 반환합니다.
func GetTxCmd() *cobra.Command {
	spanTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s 트랜잭션 명령어", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	spanTxCmd.AddCommand(
		NewCreateSpanCmd(),
		NewUpdateParamsCmd(),
	)

	return spanTxCmd
}

// NewCreateSpanCmd는 새로운 스팬을 생성하는 명령어를 반환합니다.
func NewCreateSpanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-span [id] [start-block] [end-block] [validator-set-file] [selected-producers] [chain-id]",
		Short: "새로운 스팬을 생성합니다",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("스팬 ID를 파싱할 수 없습니다: %w", err)
			}

			startBlock, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("시작 블록을 파싱할 수 없습니다: %w", err)
			}

			endBlock, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("종료 블록을 파싱할 수 없습니다: %w", err)
			}

			// 검증자 세트 파일 읽기
			validatorSetBytes, err := os.ReadFile(args[3])
			if err != nil {
				return fmt.Errorf("검증자 세트 파일을 읽을 수 없습니다: %w", err)
			}

			var validatorSet []*types.Validator
			if err := json.Unmarshal(validatorSetBytes, &validatorSet); err != nil {
				return fmt.Errorf("검증자 세트를 파싱할 수 없습니다: %w", err)
			}

			// 선택된 생산자 주소 파싱
			selectedProducers := strings.Split(args[4], ",")

			chainID := args[5]

			msg := types.NewMsgCreateSpan(
				clientCtx.GetFromAddress().String(),
				id,
				startBlock,
				endBlock,
				validatorSet,
				selectedProducers,
				chainID,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUpdateParamsCmd는 모듈 파라미터를 업데이트하는 명령어를 반환합니다.
func NewUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [params-file]",
		Short: "모듈 파라미터를 업데이트합니다",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// 파라미터 파일 읽기
			paramsBytes, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("파라미터 파일을 읽을 수 없습니다: %w", err)
			}

			var params types.Params
			if err := json.Unmarshal(paramsBytes, &params); err != nil {
				return fmt.Errorf("파라미터를 파싱할 수 없습니다: %w", err)
			}

			msg := types.NewMsgUpdateParams(
				clientCtx.GetFromAddress().String(),
				params,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSubmitUpdateSpanParamsProposalTxCmd는 스팬 파라미터 업데이트 제안을 제출하는 명령어를 반환합니다.
func NewSubmitUpdateSpanParamsProposalTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-span-params [proposal-file]",
		Short: "스팬 파라미터 업데이트 제안을 제출합니다",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// 제안 파일 읽기
			proposalBytes, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("제안 파일을 읽을 수 없습니다: %w", err)
			}

			var proposal types.UpdateSpanParamsProposal
			if err := clientCtx.Codec.UnmarshalJSON(proposalBytes, &proposal); err != nil {
				return fmt.Errorf("제안을 파싱할 수 없습니다: %w", err)
			}

			from := clientCtx.GetFromAddress()
			content := &proposal

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := cli.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
