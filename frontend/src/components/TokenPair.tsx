import type { TokenSnapshot } from '../lib/api/types'

function TokenMark({ token }: { token: TokenSnapshot }) {
  return (
    <span className="token-mark" title={token.symbol}>
      {token.symbol.slice(0, 2).toUpperCase()}
    </span>
  )
}

export function TokenPair({
  lend,
  collateral,
}: {
  lend: TokenSnapshot
  collateral: TokenSnapshot
}) {
  return (
    <div className="token-pair">
      <span className="token-pair__marks">
        <TokenMark token={lend} />
        <TokenMark token={collateral} />
      </span>
      <span>
        <strong>{lend.symbol}</strong>
        <small>against {collateral.symbol}</small>
      </span>
    </div>
  )
}
