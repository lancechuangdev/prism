import { Link } from '../routing'

export function NotFoundPage() {
  return (
    <section className="page placeholder">
      <p className="eyebrow">404</p>
      <h1>Page not found</h1>
      <p>The page you requested does not exist.</p>
      <Link className="text-link" to="/">
        Return home →
      </Link>
    </section>
  )
}
