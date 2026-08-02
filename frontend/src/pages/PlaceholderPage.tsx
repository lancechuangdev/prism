import { useWallet } from '../wallet/WalletProvider'
import { Button } from '../components/Button'

type PlaceholderPageProps = {
  eyebrow: string
  title: string
  description: string
  requiresWallet?: boolean
}

export function PlaceholderPage({
  eyebrow,
  title,
  description,
  requiresWallet,
}: PlaceholderPageProps) {
  const wallet = useWallet()

  return (
    <section className="page placeholder">
      <p className="eyebrow">{eyebrow}</p>
      <h1>{title}</h1>
      <p>{description}</p>
      {requiresWallet && !wallet.account && (
        <Button onClick={() => void wallet.connect()}>Connect wallet</Button>
      )}
      {!requiresWallet && (
        <p className="status-pill">
          <span aria-hidden="true" /> Read-only mode available
        </p>
      )}
    </section>
  )
}
