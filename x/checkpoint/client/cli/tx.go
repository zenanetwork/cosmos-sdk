package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/types"
)

// GetTxCmd는 checkpoint 모듈의 트랜잭션 명령어를 반환합니다.
func GetTxCmd() *cobra.Command {
	checkpointTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s 트랜잭션 명령어", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	checkpointTxCmd.AddCommand(
		NewCreateCheckpointCmd(),
		NewUpdateParamsCmd(),
	)

	return checkpointTxCmd
}

// NewCreateCheckpointCmd는 새로운 체크포인트를 생성하는 명령어를 반환합니다.
func NewCreateCheckpointCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-checkpoint [number] [start-block] [end-block] [root-hash] [proposer]",
		Short: "새로운 체크포인트를 생성합니다",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			number, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("체크포인트 번호를 파싱할 수 없습니다: %w", err)
			}

			startBlock, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("시작 블록을 파싱할 수 없습니다: %w", err)
			}

			endBlock, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("종료 블록을 파싱할 수 없습니다: %w", err)
			}

			rootHash := []byte(args[3])
			proposer := args[4]

			msg := types.NewMsgCreateCheckpoint(
				clientCtx.GetFromAddress().String(),
				number,
				startBlock,
				endBlock,
				rootHash,
				proposer,
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
