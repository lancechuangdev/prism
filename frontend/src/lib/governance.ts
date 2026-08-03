import { getAddress, isAddress } from 'viem'

import type { ProposalOperation } from './api/types'

export const operationLabels: Record<string, string> = {
  add_owner: 'Add owner',
  remove_owner: 'Remove owner',
  replace_owner: 'Replace owner',
  change_threshold: 'Change threshold',
  create_pool: 'Create pool',
  settle_pool: 'Settle pool',
  repay_pool: 'Repay pool',
  liquidate_pool: 'Liquidate pool',
}

const uint = /^\d+$/
const positiveUint = /^[1-9]\d*$/

export function validateOperation(operation: ProposalOperation) {
  const errors: string[] = []
  const params = operation.params
  const addressFields =
    operation.type === 'create_pool'
      ? [
          'lendToken',
          'collateralToken',
          'lenderPositionToken',
          'borrowerPositionToken',
        ]
      : operation.type === 'replace_owner'
        ? ['old_owner', 'new_owner']
        : ['add_owner', 'remove_owner'].includes(operation.type)
          ? ['owner']
          : []
  for (const field of addressFields) {
    if (!isAddress(String(params[field] ?? '')))
      errors.push(`${field} must be a valid address.`)
  }
  const unsignedFields =
    operation.type === 'create_pool'
      ? [
          'settleTime',
          'maturityTime',
          'interestRate',
          'maxLendSupply',
          'collateralizationRatio',
          'liquidateRate',
        ]
      : ['settle_pool', 'repay_pool', 'liquidate_pool'].includes(operation.type)
        ? ['poolId']
        : []
  for (const field of unsignedFields) {
    if (!uint.test(String(params[field] ?? '')))
      errors.push(`${field} must be a non-negative decimal integer.`)
  }
  if (
    ['repay_pool', 'liquidate_pool'].includes(operation.type) &&
    !positiveUint.test(String(params.maxCollateralAmount ?? ''))
  )
    errors.push('maxCollateralAmount must be a positive decimal integer.')
  if (
    operation.type === 'change_threshold' &&
    (!Number.isInteger(Number(params.threshold)) ||
      Number(params.threshold) < 1)
  )
    errors.push('Threshold must be a positive whole number.')
  if (
    operation.type === 'create_pool' &&
    uint.test(String(params.maturityTime ?? '')) &&
    uint.test(String(params.settleTime ?? '')) &&
    BigInt(String(params.maturityTime)) <= BigInt(String(params.settleTime))
  )
    errors.push('Maturity must be after settlement.')
  return errors
}

export function operationSummary(operation: ProposalOperation) {
  const label = operationLabels[operation.type] ?? operation.type
  if (operation.type === 'replace_owner')
    return `${label}: ${operation.params.old_owner} → ${operation.params.new_owner}`
  if ('owner' in operation.params) return `${label}: ${operation.params.owner}`
  if ('poolId' in operation.params) return `${label} ${operation.params.poolId}`
  if ('threshold' in operation.params)
    return `${label} to ${operation.params.threshold} signature(s)`
  return label
}

export function normalizeOperation(
  operation: ProposalOperation,
): ProposalOperation {
  return {
    ...operation,
    params: Object.fromEntries(
      Object.entries(operation.params).map(([key, value]) => [
        key,
        isAddress(String(value)) ? getAddress(String(value)) : value,
      ]),
    ),
  }
}
