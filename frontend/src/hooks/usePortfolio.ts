import { useCallback, useEffect, useState } from 'react'
import { getAddress, type Address } from 'viem'

import { config } from '../config/env'
import { erc20Abi, prismPoolAbi } from '../lib/contracts/abis'
import { publicClient } from '../lib/contracts/client'
import type {
  ActivityKind,
  PortfolioActivity,
  UserPoolPosition,
} from '../lib/portfolio'
import type { PoolRecord } from '../lib/pools'

type PortfolioState = {
  positions: UserPoolPosition[]
  activity: PortfolioActivity[]
  loading: boolean
  ignoredIndexedPools: number
  error?: string
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
    : 'Wallet positions could not be loaded.'
}

export function usePortfolio(pools: PoolRecord[], account?: Address) {
  const [state, setState] = useState<PortfolioState>({
    positions: [],
    activity: [],
    loading: false,
    ignoredIndexedPools: 0,
  })

  const load = useCallback(async () => {
    if (!account) {
      setState({
        positions: [],
        activity: [],
        loading: false,
        ignoredIndexedPools: 0,
      })
      return
    }
    setState((current) => ({ ...current, loading: true, error: undefined }))
    try {
      const [positionResult, activity] = await Promise.all([
        readPositions(pools, account),
        readActivity(account),
      ])
      setState({
        positions: positionResult.positions,
        ignoredIndexedPools: positionResult.ignoredIndexedPools,
        activity,
        loading: false,
      })
    } catch (cause) {
      setState((current) => ({
        ...current,
        loading: false,
        error: errorMessage(cause),
      }))
    }
  }, [account, pools])

  useEffect(() => {
    queueMicrotask(() => void load())
  }, [load])

  return { ...state, refresh: load }
}

async function readPositions(pools: PoolRecord[], account: Address) {
  const poolCount = await publicClient.readContract({
    address: config.contracts.pool,
    abi: prismPoolAbi,
    functionName: 'poolCount',
  })
  const livePools = pools.filter((pool) => BigInt(pool.index) < poolCount)
  const positions = await Promise.all(
    livePools.map(async (pool): Promise<UserPoolPosition> => {
      const [
        lender,
        borrower,
        lenderPositionBalance,
        borrowerPositionBalance,
        liveState,
      ] = await Promise.all([
        publicClient.readContract({
          address: config.contracts.pool,
          abi: prismPoolAbi,
          functionName: 'userLendInfo',
          args: [account, BigInt(pool.index)],
        }),
        publicClient.readContract({
          address: config.contracts.pool,
          abi: prismPoolAbi,
          functionName: 'userBorrowInfo',
          args: [account, BigInt(pool.index)],
        }),
        publicClient.readContract({
          address: getAddress(pool.base.lenderPositionToken),
          abi: erc20Abi,
          functionName: 'balanceOf',
          args: [account],
        }),
        publicClient.readContract({
          address: getAddress(pool.base.borrowerPositionToken),
          abi: erc20Abi,
          functionName: 'balanceOf',
          args: [account],
        }),
        publicClient.readContract({
          address: config.contracts.pool,
          abi: prismPoolAbi,
          functionName: 'getPoolState',
          args: [BigInt(pool.index)],
        }),
      ])
      return {
        pool,
        liveState: String(liveState) as UserPoolPosition['liveState'],
        lender: {
          stakeAmount: lender[0],
          refundAmount: lender[1],
          hasRefunded: lender[2],
          hasClaimed: lender[3],
          positionBalance: lenderPositionBalance,
        },
        borrower: {
          stakeAmount: borrower[0],
          refundAmount: borrower[1],
          hasRefunded: borrower[2],
          hasClaimed: borrower[3],
          positionBalance: borrowerPositionBalance,
        },
      }
    }),
  )
  return {
    positions: positions.filter(
      (position) =>
        position.lender.stakeAmount > 0n ||
        position.borrower.stakeAmount > 0n ||
        position.lender.positionBalance > 0n ||
        position.borrower.positionBalance > 0n,
    ),
    ignoredIndexedPools: pools.length - livePools.length,
  }
}

async function readActivity(account: Address) {
  const common = {
    address: config.contracts.pool,
    abi: prismPoolAbi,
    fromBlock: config.contracts.deploymentBlock,
  } as const
  const [
    depositLend,
    depositBorrow,
    refundLend,
    refundBorrow,
    claimLend,
    claimBorrow,
    withdrawLend,
    withdrawBorrow,
  ] = await Promise.all([
    publicClient.getContractEvents({
      ...common,
      eventName: 'DepositLend',
      args: { lender: account },
    }),
    publicClient.getContractEvents({
      ...common,
      eventName: 'DepositBorrow',
      args: { borrower: account },
    }),
    publicClient.getContractEvents({
      ...common,
      eventName: 'RefundLend',
      args: { lender: account },
    }),
    publicClient.getContractEvents({
      ...common,
      eventName: 'RefundBorrow',
      args: { borrower: account },
    }),
    publicClient.getContractEvents({
      ...common,
      eventName: 'ClaimLend',
      args: { lender: account },
    }),
    publicClient.getContractEvents({
      ...common,
      eventName: 'ClaimBorrow',
      args: { borrower: account },
    }),
    publicClient.getContractEvents({
      ...common,
      eventName: 'WithdrawLend',
      args: { lender: account },
    }),
    publicClient.getContractEvents({
      ...common,
      eventName: 'WithdrawBorrow',
      args: { borrower: account },
    }),
  ])

  const activity: PortfolioActivity[] = []
  const add = (
    kind: ActivityKind,
    logs: readonly {
      args: Record<string, unknown>
      transactionHash: `0x${string}` | null
      blockNumber: bigint | null
    }[],
    amountKey: string,
  ) => {
    for (const log of logs) {
      if (!log.transactionHash || log.blockNumber === null) continue
      activity.push({
        kind,
        poolIndex: Number(log.args.poolId),
        amount: BigInt(log.args[amountKey] as bigint),
        token: log.args.token
          ? getAddress(log.args.token as string)
          : undefined,
        transactionHash: log.transactionHash,
        blockNumber: log.blockNumber,
      })
    }
  }
  add('deposit-lend', depositLend, 'amount')
  add('deposit-borrow', depositBorrow, 'amount')
  add('refund-lend', refundLend, 'amount')
  add('refund-borrow', refundBorrow, 'amount')
  add('claim-lend', claimLend, 'spAmount')
  add('claim-borrow', claimBorrow, 'jpAmount')
  add('withdraw-lend', withdrawLend, 'lendAmount')
  add('withdraw-borrow', withdrawBorrow, 'collateralAmount')
  return activity.sort((left, right) =>
    left.blockNumber === right.blockNumber
      ? 0
      : left.blockNumber > right.blockNumber
        ? -1
        : 1,
  )
}
