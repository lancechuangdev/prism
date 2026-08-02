import { parseAbi } from 'viem'

export const prismPoolAbi = parseAbi([
  'function depositLend(uint256 poolId, uint256 amount)',
  'function depositBorrow(uint256 poolId, uint256 amount)',
  'function refundExcessLend(uint256 poolId)',
  'function refundExcessCollateral(uint256 poolId)',
  'function claimLenderPosition(uint256 poolId)',
  'function claimBorrowerPositionAndLoan(uint256 poolId)',
  'function globalPaused() view returns (bool)',
  'function minLendAmount() view returns (uint256)',
  'function minBorrowAmount() view returns (uint256)',
  'function getPoolState(uint256 poolId) view returns (uint8)',
  'function poolCount() view returns (uint256)',
  'function userLendInfo(address user, uint256 poolId) view returns (uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)',
  'function userBorrowInfo(address user, uint256 poolId) view returns (uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)',
  'event DepositLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)',
  'event DepositBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)',
  'event RefundLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)',
  'event RefundBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)',
  'event ClaimLend(address indexed lender, uint256 indexed poolId, address indexed spToken, uint256 spAmount)',
  'event ClaimBorrow(address indexed borrower, uint256 indexed poolId, address indexed jpToken, uint256 jpAmount, uint256 loanAmount)',
  'event WithdrawLend(address indexed lender, uint256 indexed poolId, uint256 spAmount, uint256 lendAmount)',
  'event WithdrawBorrow(address indexed borrower, uint256 indexed poolId, uint256 jpAmount, uint256 collateralAmount)',
])

export const erc20Abi = parseAbi([
  'function allowance(address owner, address spender) view returns (uint256)',
  'function approve(address spender, uint256 amount) returns (bool)',
  'function balanceOf(address account) view returns (uint256)',
  'function decimals() view returns (uint8)',
  'function symbol() view returns (string)',
])
