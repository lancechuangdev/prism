import { useState } from 'react'
import type { Hash } from 'viem'

import { config } from '../config/env'
import { prismPoolAbi } from '../lib/contracts/abis'
import { publicClient } from '../lib/contracts/client'
import {
  borrowerClaimAvailable,
  borrowerRedemption,
  borrowerRefundAvailable,
  lenderClaimAvailable,
  lenderRedemption,
  lenderRefundAvailable,
  type UserPoolPosition,
} from '../lib/portfolio'
import { useWallet } from '../wallet/WalletProvider'
import { Button } from './Button'

type LifecycleAction =
  | 'refund-lend'
  | 'refund-borrow'
  | 'claim-lend'
  | 'claim-borrow'
  | 'redeem-lend'
  | 'redeem-borrow'
type ActionState = {
  stage: 'idle' | 'wallet' | 'confirming' | 'success' | 'error'
  action?: LifecycleAction
  hash?: Hash
  message?: string
}

const actionConfig: Record<
  LifecycleAction,
  {
    functionName:
      | 'refundExcessLend'
      | 'refundExcessCollateral'
      | 'claimLenderPosition'
      | 'claimBorrowerPositionAndLoan'
      | 'withdrawLend'
      | 'withdrawBorrow'
    positionSide?: 'lender' | 'borrower'
    pending: string
    success: string
  }
> = {
  'refund-lend': {
    functionName: 'refundExcessLend',
    pending: 'Refund unused lending funds',
    success: 'Unused lending funds refunded.',
  },
  'refund-borrow': {
    functionName: 'refundExcessCollateral',
    pending: 'Refund unused collateral',
    success: 'Unused collateral refunded.',
  },
  'claim-lend': {
    functionName: 'claimLenderPosition',
    pending: 'Claim lender position',
    success: 'Lender position tokens claimed.',
  },
  'claim-borrow': {
    functionName: 'claimBorrowerPositionAndLoan',
    pending: 'Claim borrower position and loan',
    success: 'Borrower position tokens and loan claimed.',
  },
  'redeem-lend': {
    functionName: 'withdrawLend',
    positionSide: 'lender',
    pending: 'Redeem lender position',
    success: 'Lender position redeemed for available proceeds.',
  },
  'redeem-borrow': {
    functionName: 'withdrawBorrow',
    positionSide: 'borrower',
    pending: 'Redeem borrower position',
    success: 'Borrower position redeemed for available collateral.',
  },
}

function errorMessage(cause: unknown) {
  if (
    typeof cause === 'object' &&
    cause &&
    'shortMessage' in cause &&
    typeof cause.shortMessage === 'string'
  )
    return cause.shortMessage
  return cause instanceof Error
    ? cause.message
    : 'The transaction could not be completed.'
}

export function LifecycleActions({
  position,
  paused,
  onConfirmed,
}: {
  position: UserPoolPosition
  paused?: boolean
  onConfirmed: () => void
}) {
  const wallet = useWallet()
  const [state, setState] = useState<ActionState>({ stage: 'idle' })
  const actions: LifecycleAction[] = []
  if (lenderRefundAvailable(position) > 0n) actions.push('refund-lend')
  if (borrowerRefundAvailable(position) > 0n) actions.push('refund-borrow')
  if (lenderClaimAvailable(position) > 0n) actions.push('claim-lend')
  if (borrowerClaimAvailable(position).position > 0n)
    actions.push('claim-borrow')
  if (lenderRedemption(position) > 0n) actions.push('redeem-lend')
  if (borrowerRedemption(position) > 0n) actions.push('redeem-borrow')

  async function execute(action: LifecycleAction) {
    if (!wallet.account || !wallet.walletClient || paused) return
    const selected = actionConfig[action]
    setState({ stage: 'wallet', action })
    try {
      const simulation = await publicClient.simulateContract({
        account: wallet.account,
        address: config.contracts.pool,
        abi: prismPoolAbi,
        functionName: selected.functionName,
        args: selected.positionSide
          ? [
              BigInt(position.pool.index),
              position[selected.positionSide].positionBalance,
            ]
          : [BigInt(position.pool.index)],
      })
      const hash = await wallet.walletClient.writeContract(simulation.request)
      setState({ stage: 'confirming', action, hash })
      const receipt = await publicClient.waitForTransactionReceipt({
        hash,
        confirmations: 1,
        onReplaced: (replacement) =>
          setState({
            stage: 'confirming',
            action,
            hash: replacement.transaction.hash,
            message: 'Tracking the replacement transaction…',
          }),
      })
      if (receipt.status !== 'success')
        throw new Error('The transaction reverted onchain.')
      setState({
        stage: 'success',
        action,
        hash: receipt.transactionHash,
        message: selected.success,
      })
      onConfirmed()
    } catch (cause) {
      setState({ stage: 'error', action, message: errorMessage(cause) })
    }
  }

  if (position.liveState === '4')
    return (
      <p className="form-error">
        This pool was cancelled, but the deployed contract has no cancelled-pool
        refund function. Funds cannot be recovered through the frontend until
        the protocol adds and audits that capability.
      </p>
    )
  if (actions.length === 0)
    return (
      <p className="position-caught-up">
        No claim or refund is currently available.
      </p>
    )
  const busy = state.stage === 'wallet' || state.stage === 'confirming'

  return (
    <div className="lifecycle-actions">
      <div>
        {actions.map((action) => (
          <Button
            key={action}
            variant="secondary"
            disabled={busy || paused || wallet.status !== 'connected'}
            onClick={() => void execute(action)}
          >
            {busy && state.action === action
              ? state.stage === 'wallet'
                ? 'Confirm in wallet…'
                : 'Confirming…'
              : actionConfig[action].pending}
          </Button>
        ))}
      </div>
      {paused && (
        <p className="form-error">The Prism pool contract is paused.</p>
      )}
      {!position.redemptionSafe &&
        (position.lender.positionBalance > 0n ||
          position.borrower.positionBalance > 0n) && (
          <p className="form-error">
            Redemption is disabled because position-token isolation cannot be
            proven across every live pool. The index may be incomplete or this
            deployment may reuse a position-token contract.
          </p>
        )}
      {state.message && (
        <div
          className={`transaction-message transaction-message--${state.stage}`}
          role="status"
        >
          <p>{state.message}</p>
          {state.hash && config.chain.explorerUrl && (
            <a
              href={`${config.chain.explorerUrl}/tx/${state.hash}`}
              target="_blank"
              rel="noreferrer"
            >
              View transaction ↗
            </a>
          )}
        </div>
      )}
    </div>
  )
}
