import { Link } from '../routing'

export function HomePage() {
  return (
    <section className="hero page">
      <p className="eyebrow">Fixed-rate, onchain credit</p>
      <h1>
        Clear terms.
        <br />
        Predictable capital.
      </h1>
      <p className="hero__lede">
        Explore transparent lending pools with fixed rates, visible collateral,
        and an auditable lifecycle from funding to repayment.
      </p>
      <div className="hero__actions">
        <Link className="button button--primary" to="/pools">
          Explore pools
        </Link>
        <Link className="text-link" to="/portfolio">
          View your positions <span aria-hidden="true">→</span>
        </Link>
      </div>
      <div className="principles" aria-label="Protocol principles">
        <article>
          <span>01</span>
          <h2>Fixed terms</h2>
          <p>
            Rates and maturity are established before capital enters a pool.
          </p>
        </article>
        <article>
          <span>02</span>
          <h2>Visible risk</h2>
          <p>Collateral and liquidation conditions stay inspectable onchain.</p>
        </article>
        <article>
          <span>03</span>
          <h2>User custody</h2>
          <p>Your wallet signs every approval, deposit, refund, and claim.</p>
        </article>
      </div>
    </section>
  )
}
