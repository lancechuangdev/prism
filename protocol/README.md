# prism smart contracts

Fixed-rate lending protocol distilled from previous DeFi experience.

Lending Pool Transitions:
```
FUNDING ──settle successfully──> ACTIVE ──repay──> REPAID
    │                                │
    │                                └─liquidate─> LIQUIDATED
    │
    └──settlement fails────────> CANCELLED
```