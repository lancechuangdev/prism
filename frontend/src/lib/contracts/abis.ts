import { parseAbi } from 'viem'

export const prismPoolAbi = parseAbi([
  'function depositLend(uint256 poolId, uint256 amount)',
  'function depositBorrow(uint256 poolId, uint256 amount)',
  'function refundExcessLend(uint256 poolId)',
  'function refundExcessCollateral(uint256 poolId)',
  'function claimLenderPosition(uint256 poolId)',
  'function claimBorrowerPositionAndLoan(uint256 poolId)',
  'function paused() view returns (bool)',
  'function userLendInfo(address user, uint256 poolId) view returns (uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)',
  'function userBorrowInfo(address user, uint256 poolId) view returns (uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)',
])

export const erc20Abi = parseAbi([
  'function allowance(address owner, address spender) view returns (uint256)',
  'function approve(address spender, uint256 amount) returns (bool)',
  'function balanceOf(address account) view returns (uint256)',
  'function decimals() view returns (uint8)',
  'function symbol() view returns (string)',
])
