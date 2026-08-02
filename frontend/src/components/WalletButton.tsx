import { formatAddress } from '../lib/format'
import { useWallet } from '../wallet/WalletProvider'
import { Button } from './Button'

export function WalletButton() {
  const wallet = useWallet()

  if (wallet.status === 'wrong-network') {
    return (
      <Button type="button" onClick={() => void wallet.switchNetwork()}>
        Switch network
      </Button>
    )
  }

  if (wallet.account) {
    return (
      <Button
        type="button"
        variant="secondary"
        onClick={wallet.disconnect}
        aria-label={`Disconnect ${wallet.account}`}
      >
        {formatAddress(wallet.account)}
      </Button>
    )
  }

  return (
    <Button
      type="button"
      disabled={wallet.status === 'connecting'}
      onClick={() => void wallet.connect()}
    >
      {wallet.status === 'connecting' ? 'Connecting…' : 'Connect wallet'}
    </Button>
  )
}
