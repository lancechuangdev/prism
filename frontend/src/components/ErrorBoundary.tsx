import {
  Component,
  type ErrorInfo,
  type PropsWithChildren,
  type ReactNode,
} from 'react'

type ErrorBoundaryState = { error?: Error }

export class ErrorBoundary extends Component<
  PropsWithChildren,
  ErrorBoundaryState
> {
  state: ErrorBoundaryState = {}

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('Prism frontend render failed', error, info.componentStack)
  }

  render(): ReactNode {
    if (this.state.error) {
      return (
        <main className="page placeholder" id="main-content">
          <p className="eyebrow">Application error</p>
          <h1>Prism could not load</h1>
          <p>
            Check the frontend configuration and refresh the page. No wallet
            transaction was submitted.
          </p>
          <button
            className="button button--primary"
            type="button"
            onClick={() => window.location.reload()}
          >
            Reload Prism
          </button>
        </main>
      )
    }

    return this.props.children
  }
}
